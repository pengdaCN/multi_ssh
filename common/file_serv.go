package common

import (
	"context"
	"fmt"
	"multi_ssh/tools"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	fileServPrefix = "/__multi_ssh__/resource/static/"
)

type fileServe struct {
	shareFile map[string]string
	fileMap   map[string]string
	listen    net.Listener
	started   bool
	mu        sync.RWMutex
	serve     http.Server
	ips       []*net.IPNet
}

var (
	o                sync.Once
	DefaultFileServe = NewFileServe()
)

func NewFileServe() *fileServe {
	ips, err := tools.ExternalIP()
	if err != nil {
		panic(err)
	}
	return &fileServe{
		ips:       ips,
		shareFile: make(map[string]string),
		fileMap:   make(map[string]string),
	}
}

func (f *fileServe) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, fileServPrefix) {
		http.NotFound(w, r)
		return
	}
	downloadedFile := strings.Replace(r.URL.Path, fileServPrefix, "", 1)
	var sealFile string
	if sealFile = f.getFile(downloadedFile); sealFile == "" {
		http.NotFound(w, r)
		return
	}
	hashStr, err := tools.HashFile(sealFile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Etag", fmt.Sprintf("\"%s\"", hashStr))
	file, err := os.Open(sealFile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.ServeContent(w, r, "", time.Time{}, file)
}

func (f *fileServe) containFile(fil string) bool {
	f.mu.RLock()
	if _, ok := f.fileMap[fil]; ok {
		return ok
	}
	f.mu.RUnlock()
	return false
}

func (f *fileServe) getFile(fil string) string {
	f.mu.RLock()
	if v, ok := f.fileMap[fil]; ok {
		return v
	}
	f.mu.RUnlock()
	return ""
}

func (f *fileServe) addFile(file string) {
	if _, ok := f.shareFile[file]; !ok {
		var n string
		for {
			n = tools.GenerateRandomStr(20)
			if _, ok := f.fileMap[n]; !ok {
				break
			}
		}
		f.shareFile[file] = n
		f.fileMap[n] = file
	} else {
		return
	}
}

func (f *fileServe) AddFile(str ...string) {
	o.Do(func() {
		f.Start()
	})
	f.mu.Lock()
	defer f.mu.Unlock()
	for _, v := range str {
		f.addFile(v)
	}
}

func (f *fileServe) Start() {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.started {
		return
	}
	f.listen, _ = net.Listen("tcp", "0.0.0.0:0")
	f.serve = http.Server{
		Handler: f,
	}
	go func() {
		_ = f.serve.Serve(f.listen)
	}()
	f.started = true
}

func (f *fileServe) Stop() {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.started {
		_ = f.serve.Shutdown(context.Background())
		f.started = false
	}
}

func (f *fileServe) buildUrl(ip net.IP, str string) string {
	for _, v := range f.ips {
		if v.Contains(ip) {
			return fmt.Sprintf(
				"http://%s:%d%s%s",
				v.IP.String(),
				f.listen.Addr().(*net.TCPAddr).Port,
				fileServPrefix,
				str,
			)
		}
	}
	return ""
}

func (f *fileServe) Exists(filename string) bool {
	f.mu.RLock()
	_, ok := f.shareFile[filename]
	f.mu.RUnlock()
	return ok
}

func (f *fileServe) GetUrl(ip net.IP, filename string) (str string) {
	if f.Exists(filename) {
		f.mu.RLock()
		str = f.buildUrl(ip, f.shareFile[filename])
		f.mu.RUnlock()
	}
	return
}
