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
	"path"
	"path/filepath"
	"strings"
)

// 判断路径是否已用户家目录表示符号开始的
// 目前只支持~方式
func switchStartHomeSymbol(path string) bool {
	paths := filepath.SplitList(path)
	return strings.HasPrefix(paths[0], "~")
}

func (t *Terminal) expandDir(path string) (string, bool) {
	dirs := filepath.SplitList(path)
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
//fn 对上传文件设置额外操作
func (t *Terminal) Copy(exists, sudo bool, srcPaths []string, remotePath string, fn handleByFile) error {
	paths := filepath.SplitList(remotePath)
	expandPath, ok := t.expandDir(remotePath)
	if strings.HasPrefix(paths[0], "~") || strings.HasPrefix(remotePath, "/tmp") {
		if exists {
			if ok {
				_, err := t.run(false, fmt.Sprintf(`test -d %s`, expandPath))
				if err != nil {
					err = t.sftpClient.MkdirAll(expandPath)
					return err
				}
			} else {
				return errors.New("不能对目录进行扩展")
			}
		}
		err := t.SftpUpdates(srcPaths, expandPath, fn)
		if err != nil {
			return err
		}
	} else if sudo {
		if exists {
			if ok {
				_, err := t.run(sudo , fmt.Sprintf(`sudo test -d %s`, expandPath))
				if err != nil {
					_, err := t.run(sudo, fmt.Sprintf("sudo mkdir -p %s", expandPath))
					if err != nil {
						return errors.New("copy 时创建目录失败")
					}
				}
			} else {
				return errors.New("不能对目录进行扩展")
			}
		}
		if t.SftpUpdates(srcPaths, "/tmp", fn) != nil {
			return errors.New("上传文件失败")
		}
		filenames := make([]string, 0, len(srcPaths))
		for i:=0 ; i<len(srcPaths); i++ {
			filenames = append(filenames, path.Base(srcPaths[i]))
		}
		_, err := t.run(sudo, fmt.Sprintf(`sudo mv /tmp/{%s} %s`, strings.Join(filenames, ","), expandPath))
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *Terminal) Remove(path string) error {
	return t.sftpClient.Remove(path)
}

type handleByFile func(*sftp.File) error

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

func (t *Terminal) SftpUpdates(srcPaths []string, remotePath string, fn handleByFile) error {
	// 对~ 表示的目录进行扩展，注意~cat, ~panda，能表示方式不会进行扩展
	dirs := filepath.SplitList(remotePath)
	if strings.HasPrefix(dirs[0], "~") {
		if info, ok := t.GetContent().GetHostInfo(); ok {
			if info.User.Home != "" {
				dirs[0] = info.User.Home
			}
		} else {
			log.Println("对~字符的扩展失败")
			return errors.New("对~字符的扩展失败")
		}
	}
	remotePath = filepath.Join(dirs...)
	for _, s := range srcPaths {
		err := t.SftpUpdate(s, remotePath, fn)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *Terminal) SftpUpdate(_path, remotePath string, fn handleByFile) error {
	b, err := ioutil.ReadFile(_path)
	if err != nil {
		panic(err)
	}
	rd := bytes.NewReader(b)
	filename := path.Base(_path)
	return t.SftpUpdateByReaderWithFunc(filename, rd, remotePath, fn)
}

func (t *Terminal) SftpUpdateByReaderWithFunc(filename string, reader io.Reader, remotePath string, fn handleByFile) error {
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
