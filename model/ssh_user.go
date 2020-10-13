package model

import (
	"bufio"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
	"multi_ssh/tools"
	"regexp"
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
	//read := bufio.NewReader(bytes.NewReader(context))
	//rst := make([]*SSHUserByPassphrase, 0)
	//// 用于选出重复的条目
	//m := make(map[string]struct{})
	//var lineNumber int
	//for {
	//	lineNumber++
	//	line, err := read.ReadString('\n')
	//	line = strings.TrimSpace(line)
	//	if err != nil {
	//		if s := ReadLine(line); s != nil {
	//			// 去除处重复的条目
	//			if !isRepeat(m, s) {
	//				s.line = lineNumber
	//				rst = append(rst, s)
	//			}
	//		}
	//		break
	//	}
	//	if s := ReadLine(line); s != nil {
	//		if !isRepeat(m, s) {
	//			s.line = lineNumber
	//			rst = append(rst, s)
	//		}
	//	}
	//}
	//return rst, nil
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
	return s
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

func (s *SSHUserByPassphrase) Line() int {
	return s.line
}

var (
	// 用于匹配 key是否是合法字符
	formatKey, _ = regexp.Compile(`[a-zA-Z][a-z0-9-A-Z]*`)
	// 用于匹配是否是合法分隔符
	_separate, _ = regexp.Compile(`\s`)
	_space, _    = regexp.Compile(`\s`)
)

const (
	// step 还没开始匹配key，准备key匹配阶段
	step uint8 = iota
	// step1 正在进行配置key，按照[a-zA-Z][a-z0-9A-Z]* 方式匹配
	step1
	// step2 结束key配置，当匹配到=号时结束key匹配，准备进行val的匹配
	step2
	// step3 当匹配到'"中的一个开始匹配val，去后面的字符都是val，当匹配到'"时结束匹配
	step3
	// step4 当匹配到转义字符时
	step4
	// step5 结果一次完整的key="val" 匹配，但是还没开始新的一次匹配
	step5
	// 中间符号，key="val" 两边没有空格
	middleSign = '='
	// 用于包括val的可允许的字符
	borderSign = `'"`
	// 用于转义的字符
	cfgWord = `\`
)

func parseExtraRetMap(str string) map[string]string {
	var (
		key    strings.Builder
		val    strings.Builder
		border string
		stat   uint8
	)
	rst := make(map[string]string)
	for _, w := range str {
		switch {
		case stat == step:
			if !_space.MatchString(string(w)) {
				stat = step1
				key.WriteRune(w)
				continue
			}
			continue
		case stat == step1:
			if w == middleSign {
				stat = step2
				continue
			}
			if !_space.MatchString(string(w)) {
				key.WriteRune(w)
				continue
			}
			goto END
		case stat == step2:
			if strings.Contains(borderSign, string(w)) {
				idx := strings.Index(borderSign, string(w))
				border = string(borderSign[idx])
				stat = step3
				continue
			}
			goto END
		case stat == step3:
			if cfgWord == string(w) {
				stat = step4
				continue
			}
			if border == string(w) {
				k := key.String()
				v := val.String()
				// key 不满足格式直接结束配置
				if !formatKey.MatchString(k) {
					goto END
				}
				rst[k] = v
				key = strings.Builder{}
				val = strings.Builder{}
				stat = step5
				continue
			}
			val.WriteRune(w)
			continue
		case stat == step4:
			val.WriteRune(w)
			continue
		case stat == step5:
			if _separate.MatchString(string(w)) {
				stat = step
				continue
			}
			// 当匹配结果后不是合法分割符，直接结束配置
			goto END
		}
	}
END:
	return rst
}

func (s *SSHUserByPassphrase) Identify() string {
	return fmt.Sprintf("SSHUserByPassphrase-%s@%s", s.UserName, s.RemoteHost)
}
