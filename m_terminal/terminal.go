package m_terminal

import (
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"multi_ssh/model"
	"time"
)

const (
	line = 1000
)

type content struct {
	data       [line]string
	updateTime time.Time
	floor      int
	heap       int
	length     int
	cap        int
}

func newContent() *content {
	return &content{
		cap: line,
	}
}

func (c *content) pop(src []byte) {
	str := string(src)
	if c.length != c.cap {
		c.length++
		c.heap++
		c.data[c.heap] = str
		return
	}
	if c.floor+1 != c.cap {
		c.floor++
	} else {
		c.floor = 0
	}
	if c.heap+1 != c.cap {
		c.heap++
	} else {
		c.heap = 0
	}
	c.data[c.heap] = str
}

func (c *content) getLast() string {
	return c.data[c.heap]
}

func (c *content) Write(src []byte) (n int, err error) {
	c.pop(src)
	c.updateTime = time.Now()
	return len(src), nil
}

type Terminal struct {
	user            model.SHHUser
	client          *ssh.Client
	sftpClient      *sftp.Client
	termCache       *content
	termStdoutCache *content
	termStderrCache *content
	iBefore         uint8
	iAfter          uint8
}

func GetSSHClientByPassphrase(user model.SHHUser) (*Terminal, error) {
	config := ssh.ClientConfig{
		User:            user.User(),
		Auth:            user.Auth(),
		Timeout:         time.Second * 5,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	client, err := ssh.Dial("tcp", user.Host(), &config)
	if err != nil {
		return nil, err
	}
	return &Terminal{
		user:            user,
		client:          client,
		termCache:       newContent(),
		termStderrCache: newContent(),
		termStdoutCache: newContent(),
	}, nil
}

func (t *Terminal) Run(sudo bool, cmd string) ([]byte, error) {
	session, err := t.NewSession()
	if err != nil {
		return nil, err
	}
	// 为了sudo的字符串可以匹配
	cmd = "LANG=en_US.utf8;LANGUAGE=en_US.utf8;" + cmd
	err = session.Run(t, sudo, cmd)
	return session.rst, err
}

func (t *Terminal) GetUser() model.SHHUser {
	return t.user
}

func (t *Terminal) Close() error {
	return t.client.Close()
}