package model

import (
	"bufio"
	"errors"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io"
	"os"
	"regexp"
	"strings"
)

type SHHUser interface {
	Host() string
	User() string
	Auth() []ssh.AuthMethod
}

var (
	seperate      *regexp.Regexp
	ignoreLine, _ = regexp.Compile(`^\s*#`)
	spaceLine, _  = regexp.Compile(`^\s+$`)
)

func init() {
	seperate, _ = regexp.Compile(`\s*,\s*`)
	ignoreLine, _ = regexp.Compile(`^\s*#`)
	spaceLine, _ = regexp.Compile(`^\s+$`)
}

type SSHUserByPassphrase struct {
	RemoteHost string
	UserName   string
	Password   string
}

func ReadHosts(path string) ([]*SSHUserByPassphrase, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("%s 打开失败 %s", path, err.Error())
	}
	buf := bufio.NewReader(f)
	rst := make([]*SSHUserByPassphrase, 0)
	for {
		line, err := buf.ReadString('\n')
		if line == "" && err == io.EOF {
			return rst, nil
		}
		if spaceLine.MatchString(line) || line == "" {
			continue
		}
		if ignoreLine.MatchString(line) {
			continue
		}
		u, err := NewSSHUserByPassphraseWithStringLine(line)
		if err != nil {
			return rst, err
		}
		rst = append(rst, u)
	}
	return rst, nil
}

func NewSSHUserByPassphraseWithStringLine(line string) (*SSHUserByPassphrase, error) {
	line = strings.TrimSpace(line)
	piece := seperate.Split(line, -1)
	if len(piece) < 3 {
		return nil, errors.New("解析一行数据错误")
	}
	u := SSHUserByPassphrase{
		UserName:   piece[0],
		Password:   piece[1],
		RemoteHost: piece[2],
	}
	return &u, nil
}

//func OpenSHHUsers(users ...SHHUser) ([]*m_terminal.Terminal, error) {
//	rst := make([]*m_terminal.Terminal, 0)
//	for _, u := range users {
//		c, err := m_terminal.GetSSHClientByPassphrase(u)
//		if err != nil {
//			return rst, err
//		}
//		rst = append(rst, c)
//	}
//	return rst, nil
//}

func (s *SSHUserByPassphrase) Host() string {
	return s.RemoteHost
}

func (s *SSHUserByPassphrase) User() string {
	return s.UserName
}

func (s *SSHUserByPassphrase) Auth() []ssh.AuthMethod {
	return []ssh.AuthMethod{ssh.Password(s.Password)}
}
