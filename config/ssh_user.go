package config

import (
	"bufio"
	"io"
	"multi_ssh/model"
	"os"
	"regexp"
	"strings"
)

/*
hosts文件格式
username, password, host:port
例如:
	cat, 123456, 192.168.0.10:22
可以通过 # 开头代表注释，忽略该行
例如：
	# cat, 123456, 192.168.0.10:22
*/

var (
	ignoreLine *regexp.Regexp
	seperate   *regexp.Regexp
	spaceLine  *regexp.Regexp
)

func init() {
	ignoreLine, _ = regexp.Compile(`^\s*#`)
	seperate, _ = regexp.Compile(`\s*,\s*`)
	spaceLine, _ = regexp.Compile(`^\s+$`)
}

func ReadHosts(path string) ([]*model.SSHUserByPassphrase, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	buf := bufio.NewReader(f)
	rst := make([]*model.SSHUserByPassphrase, 0)
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
		line = strings.TrimSpace(line)
		piece := seperate.Split(line, -1)
		u := model.SSHUserByPassphrase{
			UserName:   piece[0],
			Password:   piece[1],
			RemoteHost: piece[2],
		}
		rst = append(rst, &u)
	}
	return rst, nil
}
