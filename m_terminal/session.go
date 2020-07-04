package m_terminal

import (
	"golang.org/x/crypto/ssh"
	"io"
)

type TermSession struct {
	*ssh.Session
	TermStdin io.WriteCloser
	rst []byte
}

func (t *Terminal) NewSession() (*TermSession, error) {
	s, err := t.client.NewSession()
	if err != nil {
		return nil, err
	}
	{
		modes := ssh.TerminalModes{
			ssh.ECHO:          0,
			ssh.TTY_OP_ISPEED: 14400,
			ssh.TTY_OP_OSPEED: 14400,
		}
		if err := s.RequestPty("xterm", 40, 80, modes); err != nil {
			return nil, err
		}
	}
	ts := new(TermSession)
	ts.Session = s
	ts.TermStdin, err = s.StdinPipe()
	if err != nil {
		panic(err)
	}
	ts.Stdout = make(out, 0)
	ts.Stderr = make(out, 0)
	return ts, nil
}

//
//import (
//	"errors"
//	"golang.org/x/crypto/ssh"
//	"io"
//	"multi_ssh/model"
//	"time"
//)
//
//type OUT chan []byte
//
//func (o OUT) Write(src []byte) (n int, err error) {
//	o <- src
//	return len(src), nil
//}
//
//func (o OUT) Close() error {
//	close(o)
//	return nil
//}
//
//type TermSession struct {
//	stdin    io.WriteCloser
//	session  *ssh.Session
//	stdout   OUT
//	stderr   OUT
//	execable chan struct{}
//}
//
//func (t *TermSession) Run(cmd ...string) error {
//	if err := t.session.Shell(); err != nil {
//		return err
//	}
//	for _, c := range cmd {
//		_, err := t.stdin.Write([]byte(c + "\n"))
//		if err != nil {
//			return err
//		}
//	}
//	return t.session.Wait()
//}
//
//func (t *TermSession) RunBySudo(user model.SHHUser, cmd ...string) error {
//	if err := t.session.Shell(); err != nil {
//		return err
//	}
//	time.Sleep(time.Second * 2)
//	if _, err := t.stdin.Write([]byte("sudo -s su\n")); err != nil {
//		return err
//	}
//	u, ok := user.(*model.SSHUserByPassphrase)
//	if !ok {
//		return errors.New("ERROR sudo 用户密码有误")
//	}
//	time.Sleep(time.Second * 2)
//	if _, err := t.stdin.Write([]byte(u.Password + "\n")); err != nil {
//		return err
//	}
//	for _, c := range cmd {
//		_, err := t.stdin.Write([]byte(c + "\n"))
//		if err != nil {
//			return err
//		}
//	}
//	_, err := t.stdin.Write([]byte("exit" + "\n"))
//	if err != nil {
//		return err
//	}
//	return t.session.Wait()
//}
