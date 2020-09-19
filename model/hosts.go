package model

import (
	"regexp"
	"strings"
)

var (
	separate, _   = regexp.Compile(`\s*,\s*`)
	ignoreLine, _ = regexp.Compile(`^\s*#`)
	spaceLine, _  = regexp.Compile(`^\s+$`)
)

const (
	// 定义扩展信息开始标志
	keyWord = "`"
	// 定义包括扩展信息，主机信息结束标志
	keyWord1 = `;`
)

type RemoteHostInfo struct {
	UserName   string
	Passphrase string
	Host       string
	Extra      string
}

func ParseLine(line string) *RemoteHostInfo {
	if ignoreLine.MatchString(line) || ignoreLine.MatchString(line) {
		return nil
	}
	var (
		hostInfo  string
		extraInfo string
	)
	hostInfo = line
	if endIdx := strings.Index(line, keyWord1); endIdx != -1 {
		t := strings.Index(line, keyWord)
		if t != -1 {
			if endIdx < t {
				hostInfo = line[:endIdx]
				extraInfo = line[endIdx+1:]
			}
		}
		hostInfo = line[:endIdx]
	}
	r := new(RemoteHostInfo)
	if !parseBase(r, hostInfo) {
		return nil
	}
	r.Extra = extraInfo
	return r
}

func parseBase(r *RemoteHostInfo, str string) bool {
	arr := separate.Split(str, -1)
	if len(arr) != 3 {
		return false
	}
	r.UserName = strings.TrimSpace(arr[0])
	r.Passphrase = strings.TrimSpace(arr[1])
	r.Host = strings.TrimSpace(arr[2])
	return true
}
