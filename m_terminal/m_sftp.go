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
	"multi_ssh/extra_mod/host_info"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

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

//@exists 参数为true，上传的目录不存在就创建
//@sudo 参数为true，上传放置在任何root可以方式目录
//@fn 对上传文件设置额外操作
func (t *Terminal) Copy(exist, sudo bool, srcPaths []string, remotePath string, fn HandleByFile) *Result {
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
			err := t.SftpUpdates(srcPaths, "/tmp", fn)
			if err != nil {
				return buildRstByErr(err)
			}
			filenames := make([]string, 0, len(srcPaths))
			for i := 0; i < len(srcPaths); i++ {
				filenames = append(filenames, filepath.Base(srcPaths[i]))
			}
			if len(filenames) < 1 {
				_, err := t.run(sudo, fmt.Sprintf(`sudo mv /tmp/%s /%s`, filenames[0], expandPath))
				if err != nil {
					return buildRstByErr(err)
				}
			} else {
				_, err := t.run(sudo, fmt.Sprintf(`sudo mv /tmp/{%s} %s`, strings.Join(filenames, ","), expandPath))
				if err != nil {
					return buildRstByErr(err)
				}
			}
			return buildRstWithOK()
		}
	}
	err := t.SftpUpdates(srcPaths, expandPath, fn)
	return buildRstByErr(err)
}

// 保证远程目录一定存在
func exists(term *Terminal, sudo bool, remotePath string) error {
	var prefix string
	if sudo {
		prefix = "sudo"
	}
	_, err := term.run(sudo, fmt.Sprintf(`%s test -d %s`, prefix, remotePath))
	if err != nil {
		_, err := term.run(sudo, fmt.Sprintf(`%s mkdir -p %s`, prefix, remotePath))
		if err != nil {
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

type HandleByFile func(*sftp.File) error

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

func (t *Terminal) SftpUpdates(srcPaths []string, remotePath string, fn HandleByFile) error {
	for _, s := range srcPaths {
		err := t.SftpUpdate(s, remotePath, fn)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *Terminal) SftpUpdate(_path, remotePath string, fn HandleByFile) error {
	b, err := ioutil.ReadFile(_path)
	if err != nil {
		panic(err)
	}
	rd := bytes.NewReader(b)
	filename := path.Base(_path)
	return t.SftpUpdateByReaderWithFunc(filename, rd, remotePath, fn)
}

func (t *Terminal) SftpUpdateByReaderWithFunc(filename string, reader io.Reader, remotePath string, fn HandleByFile) error {
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
	if err == nil && fn != nil {
		err := fn(f)
		if err != nil {
			return err
		}
	}
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
