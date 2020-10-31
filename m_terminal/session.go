package m_terminal

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"io"
	"log"
	"multi_ssh/model"
	"strings"
	"sync"
)

var (
	termSessionPool = sync.Pool{
		New: func() interface{} {
			return &TermSession{
				mu: sync.RWMutex{},
			}
		},
	}
)

func getSession() *TermSession {
	return termSessionPool.Get().(*TermSession)
}

func putSession(sess *TermSession) {
	termSessionPool.Put(sess)
}

type TermSession struct {
	*ssh.Session
	TermStdin io.WriteCloser
	rst       []byte
	stdout    []byte
	stderr    []byte
	uinfo     model.SHHUser
	mu        sync.RWMutex
}

func (t *Terminal) NewSession() (*TermSession, error) {
	s, err := t.GetSessionWithTerm()
	if err != nil {
		return nil, err
	}
	ts := new(TermSession)
	ts.Session = s
	ts.TermStdin, err = s.StdinPipe()
	if err != nil {
		panic(err)
	}
	ts.Stdout = make(out, 0)
	ts.Stderr = make(out, 0)
	ts.uinfo = t.user
	return ts, nil
}

func (s *TermSession) withTerm(t *Terminal) {
	ses, err := t.GetSessionWithTerm()
	if err != nil {
		panic(err)
	}
	s.Session = ses
	s.TermStdin, err = s.StdinPipe()
	if err != nil {
		panic(err)
	}
	if s.Stdout == nil {
		s.Stdout = make(out, 0)
	}
	if s.stderr == nil {
		s.Stderr = make(out, 0)
	}
	s.uinfo = t.user
}

func (s *TermSession) reset() {
	s.Session = nil
	s.rst = s.rst[:0]
	s.stdout = s.stdout[:0]
	s.stderr = s.stderr[:0]
	s.uinfo = nil
	s.TermStdin = nil
}

func (s *TermSession) GetMsg() (rst []byte) {
	s.mu.RLock()
	rst = make([]byte, len(s.rst))
	copy(rst, s.rst)
	s.mu.RUnlock()
	return
}

func (s *TermSession) Run(enableSudo bool, cmd string) error {
	defer func() {
		_ = s.Session.Close()
	}()
	go func() {
		for {
			stdout, ok := s.Stdout.(out)
			if !ok {
				panic("stdout 不是out类型")
			}
			stderr, ok := s.Stderr.(out)
			if !ok {
				panic("stderr 不是out类型")
			}
			if s.Stdout == nil && s.Stderr == nil {
				break
			}
			select {
			case o, ok := <-stdout:
				if !ok {
					stdout = nil
					continue
				}
				s.mu.Lock()
				s.rst = append(s.rst, o...)
				s.stdout = append(s.stdout, o...)
				s.mu.Unlock()
				if enableSudo {
					if err := sudo(s.uinfo, o, s.TermStdin); err != nil {
						log.Println("sudo执行错误", err)
						break
					}
				}
			case o2, ok := <-stderr:
				if !ok {
					stderr = nil
					continue
				}
				s.mu.Lock()
				s.rst = append(s.rst, o2...)
				s.stderr = append(s.stderr, o2...)
				s.mu.Unlock()
			}
		}
	}()
	return s.Session.Run(cmd)
}

const sudoPrefix = "[sudo] password for %s: "

func sudo(u model.SHHUser, in []byte, out io.Writer) error {
	line := string(in)
	beenMatched := fmt.Sprintf(sudoPrefix, u.User())
	if strings.Contains(beenMatched, line) {
		u, _ := u.(*model.SSHUserByPassphrase)
		_, err := out.Write([]byte(u.Password + "\n"))
		if err != nil {
			return err
		}
	}
	return nil
}
