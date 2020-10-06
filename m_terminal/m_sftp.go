package m_terminal

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/pkg/sftp"
	"io"
	"io/ioutil"
	"log"
	"multi_ssh/common"
	"multi_ssh/extra_mod/host_info"
	"net"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

type (
	CopyOption struct {
		Exists bool
		Sudo   bool
		Uid    int
		Gid    int
		Mode   string
	}
	CopyMode struct {
		Gid  int
		Uid  int
		Mode os.FileMode
	}
	HttpDownloader int
)

const (
	Wget HttpDownloader = iota
	Curl
	Disable
)

var (
	hdMethod = map[HttpDownloader]string{
		Wget: "wget -c %s -O %s",
		Curl: "curl -C %s -O %s",
	}
)

func (h HttpDownloader) String() string {
	switch h {
	case Wget:
		return "wget"
	case Curl:
		return "curl"
	case Disable:
		return "DISABLE"
	default:
		return "UNKNOWN"
	}
}

func (h HttpDownloader) buildUrl(url, filename string) string {
	if v, ok := hdMethod[h]; ok {
		return fmt.Sprintf(v, url, filename)
	}
	return ""
}

func (t *Terminal) dependEnvForHttpDownload() HttpDownloader {
	if r := t.Run(false, "which wget"); r.code == 0 {
		return Wget
	}
	if r := t.Run(false, "which curl"); r.code == 0 {
		return Curl
	}
	return Disable
}

// TODO 后续完善
func buildSetModeWithCopyMode(c CopyMode) string {
	var cmd string
	if c.Mode != 0 {
		cmd += fmt.Sprintf("chmod ")
	}
	return ""
}

func (t *Terminal) copyOnHttp(sudo bool, src []string, tar string, c *CopyMode) *Result {
	downloader := t.dependEnvForHttpDownload()
	u := t.GetUser()
	ip, _, err := net.SplitHostPort(u.Host())
	if err != nil {
		panic("解析ip错误")
	}
	common.DefaultFileServe.AddFile(src...)
	for _, v := range src {
		url := common.DefaultFileServe.GetUrl(net.ParseIP(ip), v)
		if url == "" {
			return buildErrWithText(fmt.Sprintf("file %s get url faild", v))
		}
		cmd := downloader.buildUrl(url, path.Join(tar, path.Base(v)))
		var _sudo string
		if sudo {
			_sudo = "sudo "
		}
		r := t.Run(sudo, fmt.Sprintf("%s%s", _sudo, cmd))
		if r.err != nil {
			return r
		}
	}
	// TODO 后续添加权限gid，uid设置的功能
	return buildRstWithOK()
}

// 将路径拆分
// @path 要拆分的路径
// 如 /tmp/a/b/c/d/e 拆分为/ tmp/ a/ b/ c/ d/ e
// 目前该方法只能在linux使用
func pathSplit(path string) []string {
	p := filepath.Clean(path)
	sep, _ := regexp.Compile(`(.*?/)`)
	t := sep.FindAllStringSubmatch(p, -1)
	rst := make([]string, 0, len(t))
	for _, v := range t {
		rst = append(rst, v[1])
	}
	if v := filepath.Base(p); v != "" {
		rst = append(rst, v)
	}
	return rst
}

// 根据终端中获取的信息，对目录进行完整性补全
// 如 ~/.bashrc 会补全为/home/.bashrc
func (t *Terminal) expandDir(path string) (string, bool) {
	dirs := pathSplit(path)
	if strings.HasPrefix(dirs[0], "~") {
		if info, ok := t.GetContent().GetHostInfo(); ok {
			if info.User.Home != "" {
				dirs[0] = info.User.Home
			}
		} else {
			return path, false
		}
	}
	return filepath.Join(dirs...), true
}

// TODO 后续添加权限设置功能
//@exists 参数为true，上传的目录不存在就创建
//@sudo 参数为true，上传放置在任何root可以方式目录
//@fn 对上传文件设置额外操作
func (t *Terminal) Copy(exist, sudo bool, srcPaths []string, remotePath string) *Result {
	info, _ := t.GetContent().GetHostInfo()
	expandPath, _ := t.expandDir(remotePath)
	if exist {
		err := exists(t, sudo, expandPath)
		if err != nil {
			return buildRstByErr(err)
		}
	}
	if sudo {
		if !inRange(info, remotePath) {
			//err := t.SftpUpdates(srcPaths, "/tmp", fn)
			//if err != nil {
			//	return buildRstByErr(err)
			//}
			//filenames := make([]string, 0, len(srcPaths))
			//for i := 0; i < len(srcPaths); i++ {
			//	filenames = append(filenames, filepath.Base(srcPaths[i]))
			//}
			//if len(filenames) < 1 {
			//	_, err := t.run(sudo, fmt.Sprintf(`sudo mv /tmp/%s /%s`, filenames[0], expandPath))
			//	if err != nil {
			//		return buildRstByErr(err)
			//	}
			//} else {
			//	_, err := t.run(sudo, fmt.Sprintf(`sudo mv /tmp/{%s} %s`, strings.Join(filenames, ","), expandPath))
			//	if err != nil {
			//		return buildRstByErr(err)
			//	}
			//}
			//return buildRstWithOK()
			r := t.copyOnHttp(sudo, srcPaths, remotePath, nil)
			return r
		}
	}
	err := t.SftpUpdates(srcPaths, expandPath)
	return buildRstByErr(err)
}

// 保证远程目录一定存在
func exists(term *Terminal, sudo bool, remotePath string) error {
	var prefix string
	if sudo {
		prefix = "sudo"
	}
	r := term.Run(sudo, fmt.Sprintf(`%s test -d %s`, prefix, remotePath))
	if r.Err() != nil {
		r := term.Run(sudo, fmt.Sprintf(`%s mkdir -p %s`, prefix, remotePath))
		if r.Err() != nil {
			return errors.New("copy 时创建目录失败")
		}
	}
	return nil
}

// 用于判断目标路径是否在远程上有写权限
func inRange(info *host_info.HostGenericInfo, remotePath string) bool {
	paths := pathSplit(remotePath)
	if info != nil {
		if info.User.IsRoot {
			return true
		}
		if info.User.Home == paths[0] {
			return true
		}
	}
	if strings.HasPrefix(remotePath, "~") {
		return true
	}
	if strings.HasPrefix(remotePath, "/tmp") {
		return true
	}
	return false
}

func (t *Terminal) Remove(path string) error {
	return t.sftpClient.Remove(path)
}

//type HandleByFile func(*sftp.File) error

func (t *Terminal) GetSftpClient() (*sftp.Client, error) {
	return sftp.NewClient(t.client)
}

func (t *Terminal) SftpOpen(path string) ([]byte, error) {
	t.sftpReady()
	b, err := t.sftpClient.Open(path)
	defer func() {
		_ = b.Close()
	}()
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(b)
}

func (t *Terminal) SftpUpdates(srcPaths []string, remotePath string) error {
	for _, s := range srcPaths {
		err := t.SftpUpdate(s, remotePath)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *Terminal) SftpUpdate(_path, remotePath string) error {
	b, err := ioutil.ReadFile(_path)
	if err != nil {
		panic(err)
	}
	rd := bytes.NewReader(b)
	filename := path.Base(_path)
	return t.SftpUpdateByReaderWithFunc(filename, rd, remotePath)
}

func (t *Terminal) SftpUpdateByReaderWithFunc(filename string, reader io.Reader, remotePath string) error {
	t.sftpReady()
	f, err := t.sftpClient.Create(path.Join(remotePath, filename))
	defer func() {
		if f != nil {
			_ = f.Close()
		}
	}()
	if err != nil {
		return err
	}
	w := bufio.NewWriter(f)
	_, err = w.ReadFrom(reader)
	return err
}

func (t *Terminal) sftpReady() {
	if t.sftpClient == nil {
		var err error
		t.sftpClient, err = t.GetSftpClient()
		if err != nil {
			log.Println("sftp client打开失败", err.Error())
			panic(err)
		}
	}
}
