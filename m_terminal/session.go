package m_terminal

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"io"
	"log"
	"multi_ssh/model"
	"strings"
)

type TermSession struct {
	*ssh.Session
	TermStdin io.WriteCloser
	rst       []byte
	uinfo     model.SHHUser
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
				s.rst = append(s.rst, o...)
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
				s.rst = append(s.rst, o2...)
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
