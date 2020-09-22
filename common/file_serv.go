package common

import (
	"context"
	"fmt"
	"multi_ssh/tools"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const (
	fileServPrefix = "/__multi_ssh__/resource/static/"
)

type fileServe struct {
	shareFile map[string]string
	listen    net.Listener
	closed    bool
	mu        sync.RWMutex
	serve     http.Server
}

func (f *fileServe) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if ! strings.HasPrefix(r.URL.Path, fileServPrefix) {
		http.NotFound(w, r)
		return
	}
	downloadedFile := strings.Replace(r.URL.Path, fileServPrefix, "", 1)
	if !f.containFile(downloadedFile) {
		http.NotFound(w, r)
		return
	}
	hashStr, err := tools.HashFile(f.shareFile[downloadedFile])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Etag", fmt.Sprintf("\"%s\"", hashStr))
	file, err := os.Open(f.shareFile[downloadedFile])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.ServeContent(w, r, "", time.Time{}, file)
}

func (f *fileServe) containFile(fil string) bool {
	f.mu.RLock()
	if _, ok := f.shareFile[fil]; ok {
		return ok
	}
	f.mu.RUnlock()
	return false
}

func (f *fileServe) AddFile(str ...string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.shareFile == nil {
		f.shareFile = make(map[string]string)
	}
	for _, v := range str {
		filName := filepath.Base(v)
		f.shareFile[filName] = v
	}
}

func (f *fileServe) Start() {
	f.mu.Lock()
	defer f.mu.Unlock()
	if !f.closed {
		return
	}
	f.listen, _ = net.Listen("tcp", "0.0.0.0:0")
	f.serve = http.Server{
		Handler: f,
	}
	go func() {
		_ = f.serve.Serve(f.listen)
	}()
	f.closed = false
}

func (f *fileServe) Stop() {
	f.mu.Lock()
	defer f.mu.Unlock()
	_ = f.serve.Shutdown(context.Background())
	f.closed = true
}

func GetUrl(str string) string {

}
