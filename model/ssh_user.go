package model

import (
	"bufio"
	"errors"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
)

type SHHUser interface {
	Host() string
	User() string
	Auth() []ssh.AuthMethod
}

type SSHUserByPassphrase struct {
	RemoteHost string
	UserName   string
	Password   string
	ExtraField map[string]string
}

func ReadLine(line string) *SSHUserByPassphrase {
	info := ParseLine(line)
	if info == nil {
		return nil
	}
	return ParseFromRHostInfo(info)
}

func ParseFromRHostInfo(info *RemoteHostInfo) *SSHUserByPassphrase {
	s := new(SSHUserByPassphrase)
	s.UserName = info.UserName
	s.RemoteHost = info.Host
	s.Password = info.Passphrase
	if info.Extra != "" {
		if ! s.ParseExtra(info.Extra) {
			log.Printf("host %s 解析扩展字符串失败", info.Host)
		}
	}
	return s
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
	piece := separate.Split(line, -1)
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

func (s *SSHUserByPassphrase) Host() string {
	return s.RemoteHost
}

func (s *SSHUserByPassphrase) User() string {
	return s.UserName
}

func (s *SSHUserByPassphrase) Auth() []ssh.AuthMethod {
	return []ssh.AuthMethod{ssh.Password(s.Password)}
}

func (s *SSHUserByPassphrase) Extra() map[string]string {
	return s.ExtraField
}

func (s *SSHUserByPassphrase) ParseExtra(str string) bool {
	_s := strings.ReplaceAll(str, "`", "")
	m := parseExtraRetMap(_s)
	if m == nil {
		return false
	}
	s.ExtraField = m
	return true
}

var (
	formatKey, _ = regexp.Compile(`[a-zA-Z][a-z0-9-A-Z]*`)
	middleSign = '='
	borderSign = `'"`
	cfgWord = `\`
)

func parseExtraRetMap(str string) map[string]string {
	var (
		key strings.Builder
		val strings.Builder
		border string
		stat uint8
	)
	rst := make(map[string]string)
	for _, w := range str {
		switch  {
		case stat == 0b0000:
			if formatKey.MatchString(string(w)) {
				stat = 0b0001
				key.WriteRune(w)
				continue
			}
			continue
		case stat == 0b0001:
			if w == middleSign {
				stat = 0b0011
				continue
			}
			if formatKey.MatchString(string(w)) {
				key.WriteRune(w)
				continue
			}
			break
		case stat == 0b0011:
			if strings.Contains(borderSign, string(w)) {
				idx := strings.Index(borderSign, string(w))
				border = string(borderSign[idx])
				stat = 0b0111
			}
			break
		case stat == 0b0111:
			if cfgWord == string(w) {
				stat = 0b1111
				continue
			}
			if border == string(w) {
				k := key.String()
				v := val.String()
				rst[k] = v
				key = strings.Builder{}
				val = strings.Builder{}
				stat = 0b0000
				continue
			}
			val.WriteRune(w)
			continue
		case stat == 0b1111:
			val.WriteRune(w)
			continue
		}
	}
	return rst
}