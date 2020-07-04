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

func (t *Terminal) SftpUpdate(_path, remotePath string) error {
	b, err := ioutil.ReadFile(_path)
	if err != nil {
		panic(err)
	}
	rd := bytes.NewReader(b)
	filename := path.Base(_path)
	return t.SftpUpdateByReader(filename, rd, remotePath)
}

func (t *Terminal) SftpUpdateByReader(filename string, reader io.Reader, remotePath string) error {
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
