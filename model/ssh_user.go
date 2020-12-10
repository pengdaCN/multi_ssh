package model

import (
	"bufio"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
	"multi_ssh/common"
	"multi_ssh/tools"
	"strings"
)

type SHHUser interface {
	Host() string
	User() string
	Auth() []ssh.AuthMethod
	Extra() map[string]string
	Identify() string
	Line() int
}

type SSHUserByPassphrase struct {
	line       int
	RemoteHost string
	UserName   string
	Password   string
	PriFile    []byte
	ExtraField map[string]string
}

func isRepeat(m map[string]struct{}, user SHHUser) bool {
	id := user.Identify()
	if _, ok := m[id]; ok {
		return true
	}
	m[id] = struct{}{}
	return false
}

func ReadHosts(fil string) ([]*SSHUserByPassphrase, error) {
	context, err := ioutil.ReadFile(fil)
	if err != nil {
		return nil, err
	}
	return ReadLines(tools.ByteSlice2String(context))
}

func ReadLines(context string) ([]*SSHUserByPassphrase, error) {
	read := bufio.NewReader(strings.NewReader(context))
	rst := make([]*SSHUserByPassphrase, 0)
	// 用于选出重复的条目
	m := make(map[string]struct{})
	var lineNumber int
	for {
		lineNumber++
		line, err := read.ReadString('\n')
		line = strings.TrimSpace(line)
		if err != nil {
			if s := ReadLine(line); s != nil {
				// 去除处重复的条目
				if !isRepeat(m, s) {
					s.line = lineNumber
					rst = append(rst, s)
				}
			}
			break
		}
		if s := ReadLine(line); s != nil {
			if !isRepeat(m, s) {
				s.line = lineNumber
				rst = append(rst, s)
			}
		}
	}
	return rst, nil
}

func ReadLine(line string) *SSHUserByPassphrase {
	info := ParseLine(line)
	if info == nil {
		return nil
	}
	return ParseFromRHostInfo(info)
}

const priKey = "PRIKEY"

func ParseFromRHostInfo(info *RemoteHostInfo) *SSHUserByPassphrase {
	s := new(SSHUserByPassphrase)
	s.UserName = info.UserName
	s.RemoteHost = info.Host
	s.Password = info.Passphrase
	if info.Extra != "" {
		if !s.ParseExtra(info.Extra) {
			log.Printf("host %s 解析扩展字符串失败", info.Host)
		}
	}
	if path, ok := s.ExtraField[priKey]; ok {
		b, err := ioutil.ReadFile(path)
		if err != nil {
			log.Println(err)
			return nil
		}
		s.PriFile = b
	}
	return s
}

func (s *SSHUserByPassphrase) Host() string {
	return s.RemoteHost
}

func (s *SSHUserByPassphrase) User() string {
	return s.UserName
}

func (s *SSHUserByPassphrase) Auth() []ssh.AuthMethod {
	if s.PriFile != nil {
		return []ssh.AuthMethod{publicKeyAuthFunc(s.PriFile)}
	}
	return []ssh.AuthMethod{ssh.Password(s.Password)}
}

func (s *SSHUserByPassphrase) Extra() map[string]string {
	return s.ExtraField
}

func (s *SSHUserByPassphrase) ParseExtra(str string) bool {
	_s := strings.ReplaceAll(str, "`", "")
	m := parseExtra(_s)
	if m == nil {
		return false
	}
	s.ExtraField = m
	return true
}

func (s *SSHUserByPassphrase) Line() int {
	return s.line
}

var (
	startSpaceChar = common.GetRe(`^\s`)
)

func parseExtra(extraStr string) (rst map[string]string) {
	extraStr, _ = common.ReadBetween(extraStr, [2]rune{'`', '`'}, false)
	rst = make(map[string]string)
	var once bool
	if extraStr == "" {
		return
	}
	for {
		var (
			key string
			sym string
			val string
		)
		if extraStr == "" || spaceLine.MatchString(extraStr) {
			return
		}
		if (!startSpaceChar.MatchString(extraStr)) && once {
			return
		}
		key, extraStr = common.ReadWord(extraStr)
		sym, extraStr = common.ReadNotSpaceChar(extraStr)
		if sym != "=" {
			panic("parse extra error")
		}
		val, extraStr = common.ReadStr(extraStr)
		rst[key] = val
		once = true
	}
}

func (s *SSHUserByPassphrase) Identify() string {
	return fmt.Sprintf("SSHUserByPassphrase-%s@%s", s.UserName, s.RemoteHost)
}

func publicKeyAuthFunc(b []byte) ssh.AuthMethod {
	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(b)
	if err != nil {
		log.Fatal("ssh key signer failed", err)
	}
	return ssh.PublicKeys(signer)
}
