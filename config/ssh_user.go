package config

import (
	"bufio"
	"io"
	"multi_ssh/model"
	"os"
	"regexp"
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
	spaceLine  *regexp.Regexp
)

func init() {
	ignoreLine, _ = regexp.Compile(`^\s*#`)
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
		u, err := model.NewSSHUserByPassphraseWithStringLine(line)
		if err != nil {
			return rst, err
		}
		rst = append(rst, u)
	}
	return rst, nil
}
