package m_terminal

import (
	"errors"
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io"
	"multi_ssh/model"
	"multi_ssh/tools"
	"sync"
	"time"
)

const (
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
	id int
	// birthID 只能设置一个，birthID的有效值是大于0的
	birthID       int
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
	curSess       *TermSession
	mu            sync.RWMutex
	taskStat      bool
}

func DefaultWithPassphrase(user model.SHHUser) (*Terminal, error) {
	term, err := GetSSHClientByPassphrase(user)
	if err != nil {
		return nil, err
	}
	term.RreUse(ExpandCmd)
	term.RreUse(autoIncrSleep)
	term.PostUse(TrimSudo)
	return term, nil
}

const defaultOpenClientTimeout = time.Duration(time.Second * 10)

func getClientWithTimeout(user model.SHHUser) (client *ssh.Client, err error) {
	ch := make(chan struct{})
	config := ssh.ClientConfig{
		User:            user.User(),
		Auth:            user.Auth(),
		Timeout:         time.Second * 5,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	go func() {
		client, err = ssh.Dial("tcp", user.Host(), &config)
		ch <- struct{}{}
	}()
	select {
	case <-time.After(defaultOpenClientTimeout):
		return nil, errors.New("open client timeout")
	case <-ch:
		return
	}
}

func GetSSHClientByPassphrase(user model.SHHUser) (*Terminal, error) {
	client, err := getClientWithTimeout(user)
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

// 使用setbirth设置birthID 时，若修改的id小于0,则不修改，若birthID 已经设置则不修改
func (t *Terminal) SetBirthID(birthID int) {
	if t.birthID > 0 {
		return
	}
	if birthID <= 0 {
		return
	}
	t.birthID = birthID
}

func (t *Terminal) GetBirthID() int {
	return t.birthID
}

func (t *Terminal) SetShare(key string, val interface{}) {
	t.content.sharePool[key] = val
}

func (t *Terminal) GetOnceShare(key string) (interface{}, bool) {
	val, ok := t.content.sharePool[key]
	return val, ok
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

func (t *Terminal) Script(sudo bool, fil io.Reader, args string) *Result {
	filename := fmt.Sprintf(`__multi_ssh__.%s.sh`, tools.GenerateRandomStr(10))
	err := t.SftpUpdateByReaderWithFunc(filename, fil, `/tmp`)
	if err != nil {
		return buildRstByErr(err)
	}
	var prefix string
	if sudo {
		prefix = "sudo -i "
	}
	rst := t.Run(sudo, fmt.Sprintf(`%sbash /tmp/%s %s`, prefix, filename, args))
	_ = t.Remove(fmt.Sprintf(`/tmp/%s`, filename))
	return rst
}

func (t *Terminal) GetID() int {
	return t.id
}

func (t *Terminal) GetMsg() []byte {
	if t.curSess == nil {
		return nil
	}
	return t.curSess.GetMsg()
}

func (t *Terminal) Run(sudo bool, cmd string) *Result {
	defer func() {
		t.preIndex = 0
		t.postIndex = 0
	}()
	t.currentRawCmd = cmd
	for ; t.preIndex < uint8(len(t.preHook)); t.preIndex++ {
		t.preHook[t.preIndex](t)
	}
	bs, stdout, stderr, err := t.run(sudo, t.currentCmd)
	str := string(bs)
	rst := buildRst(str, err)
	rst.stdout, rst.stderr = string(stdout), string(stderr)
	t.content.result = rst
	for ; t.postIndex < uint8(len(t.postHook)); t.postIndex++ {
		t.postHook[t.postIndex](t)
	}
	return rst
}

func (t *Terminal) run(sudo bool, cmd string) ([]byte, []byte, []byte, error) {
	session := getSession()
	session.reset()
	session.withTerm(t)
	// 保证同一时间只有命令在执行
	//t.mu.Lock()
	defer func() {
		//t.mu.Unlock()
		putSession(session)
	}()
	t.curSess = session
	err := session.Run(sudo, cmd)
	return session.rst, session.stdout, session.stderr, err
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

func (t *Terminal) CfgStat() {
	t.mu.Lock()
	t.taskStat = !t.taskStat
	t.mu.Unlock()
}

func (t *Terminal) GetTaskStat() (stat bool) {
	t.mu.RLock()
	stat = t.taskStat
	t.mu.RUnlock()
	return
}
