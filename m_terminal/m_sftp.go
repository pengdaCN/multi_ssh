package m_terminal

import (
	"bufio"
	"bytes"
	"github.com/pkg/sftp"
	"io"
	"io/ioutil"
	"log"
	"path"
)

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
		_ = f.Close()
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
