package m_terminal

import (
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io"
	"multi_ssh/model"
	"multi_ssh/tools"
	"time"
)

const (
	line             = 1000
	cmdPrefixGeneric = `LANG=en_US.utf8;LANGUAGE=en_US.utf8;`
)

var (
	modes = ssh.TerminalModes{
		ssh.ECHO:          0,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}
)

type HookFunc func(*Terminal)

type Terminal struct {
	id            int
	user          model.SHHUser
	client        *ssh.Client
	sftpClient    *sftp.Client
	content       *content
	currentCmd    string
	currentRawCmd string
	preHook       []HookFunc
	preIndex      uint8
	postHook      []HookFunc
	postIndex     uint8
}

func DefaultWithPassphrase(user model.SHHUser) (*Terminal, error) {
	term, err := GetSSHClientByPassphrase(user)
	if err != nil {
		return nil, err
	}
	term.RreUse(ExpandCmd)
	return term, nil
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
		id:      tools.GetID(),
		user:    user,
		client:  client,
		content: NewContent(),
	}, nil
}

func (t *Terminal) RreUse(fn ...HookFunc) {
	t.preHook = append(t.preHook, fn...)
}

func (t *Terminal) PostUse(fn ...HookFunc) {
	t.postHook = append(t.postHook, fn...)
}

func (t *Terminal) pressCmd(cmd string) {
	t.currentCmd = cmdPrefixGeneric + cmd
}

func (t *Terminal) Script(sudo bool, fil io.Reader, args string) ([]byte, error) {
	filename := fmt.Sprintf(`__multi_ssh__.%s.sh`, tools.GenerateRandomStr(10))
	err := t.SftpUpdateByReaderWithFunc(filename, fil, `/tmp`, nil)
	if err != nil {
		return nil, err
	}
	var prefix string
	if sudo {
		prefix = "sudo -s "
	}
	rst, err := t.Run(sudo, fmt.Sprintf(`%sbash /tmp/%s %s`, prefix, filename, args))
	_ = t.Remove(fmt.Sprintf(`/tmp/%s`, filename))
	return rst, err
}

func (t *Terminal) GetID() int {
	return t.id
}

func (t *Terminal) Run(sudo bool, cmd string) ([]byte, error) {
	defer func() {
		t.preIndex = 0
		t.postIndex = 0
	}()
	t.currentRawCmd = cmd
	for ; t.preIndex < uint8(len(t.preHook)); t.preIndex++ {
		t.preHook[t.preIndex](t)
	}
	rst, err := t.run(sudo, t.currentCmd)
	if err != nil {
		//	TODO 后续添加，对于执行命令报错的处理
	}
	for ; t.postIndex < uint8(len(t.postHook)); t.postIndex++ {
		t.postHook[t.postIndex](t)
	}
	return rst, err
}

func (t *Terminal) run(sudo bool, cmd string) ([]byte, error) {
	session, err := t.NewSession()
	if err != nil {
		return nil, err
	}
	err = session.Run(sudo, cmd)
	return session.rst, err
}

func (t *Terminal) GetSessionWithTerm() (*ssh.Session, error) {
	s, err := t.client.NewSession()
	if err != nil {
		return nil, err
	}
	{
		if err := s.RequestPty("xterm", 40, 80, modes); err != nil {
			return nil, err
		}
	}
	return s, nil
}

func (t *Terminal) GetContent() *content {
	return t.content
}

func (t *Terminal) GetUser() model.SHHUser {
	return t.user
}

func (t *Terminal) Close() error {
	return t.client.Close()
}
