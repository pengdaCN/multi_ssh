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

//type chip func(*Terminal, *TermSession, []byte)

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
	//hookBeforeExec  []chip
	//hookAfterExec   []chip
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

//func (t *Terminal) RegistryHookBeforeExec(fn ...chip) {
//	t.hookBeforeExec = append(t.hookBeforeExec, fn...)
//}
//
//func (t *Terminal) RegistryHookAfterExec(fn ...chip) {
//	t.hookAfterExec = append(t.hookAfterExec, fn...)
//}

//func (t *Terminal) Run2(cmd string, sudo bool) ([]byte, error) {
//	rst := make([]byte, 0)
//	s, err := t.client.NewSession()
//	if err != nil {
//		return nil, err
//	}
//	defer func() {
//		_ = s.Close()
//	}()
//	{
//		modes := ssh.TerminalModes{
//			ssh.ECHO:          0,
//			ssh.TTY_OP_ISPEED: 14400,
//			ssh.TTY_OP_OSPEED: 14400,
//		}
//		if err := s.RequestPty("xterm", 40, 80, modes); err != nil {
//			return nil, err
//		}
//	}
//	stdout := make(out, 0)
//	stderr := make(out, 0)
//	s.Stdout = stdout
//	s.Stderr = stderr
//	stdin, err := s.StdinPipe()
//	if err != nil {
//		panic(err)
//	}
//	go func() {
//		for {
//			select {
//			case o, ok := <-stdout:
//				if !ok {
//					stdout = nil
//					continue
//				}
//				_, _ = t.termStdoutCache.Write(o)
//				rst = append(rst, o...)
//				str := string(o)
//				parten := fmt.Sprintf(sudoPrefix, t.user.User())
//				if sudo && strings.Contains(str, parten) {
//					u, _ := t.user.(*model.SSHUserByPassphrase)
//					_, err := stdin.Write([]byte(u.Password + "\n"))
//					if err != nil {
//						panic(err)
//					}
//				}
//			case o2, ok := <-stderr:
//				if !ok {
//					stderr = nil
//					continue
//				}
//				_, _ = t.termStderrCache.Write(o2)
//				rst = append(rst, o2...)
//			}
//		}
//	}()
//	err = s.Run(cmd)
//	return rst, err
//}

func (t *Terminal) Run(sudo bool, cmd string) ([]byte, error) {
	session, err := t.NewSession()
	if err != nil {
		return nil, err
	}
	err = session.Run(t, sudo, cmd)
	return session.rst, err
}

func (t *Terminal) GetUser() model.SHHUser {
	return t.user
}

func (t *Terminal) Close() error {
	return t.client.Close()
}