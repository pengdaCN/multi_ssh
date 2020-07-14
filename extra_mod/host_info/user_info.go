package host_info

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

type UserInfo struct {
	Groups     []Group
	UserName   string
	Home       string
	LoginShell string
	Uid        int
	Gid        int
	IsRoot     bool
}

type Group struct {
	GroupName string
	GroupId   int
}

var (
	extractId, _   = regexp.Compile(`(\d+)`)
	extractName, _ = regexp.Compile(`\d+\((\w+)\)`)
	space, _ = regexp.Compile(`\s+`)
	ParseFormatErr = errors.New("解析用户id格式错误")
)

// uid=1000(pengda) gid=1000(pengda) groups=1000(pengda),4(adm)
func ParseIdAndGroup(str string) (uid, gid int, groups []Group, err error) {
	sarr := space.Split(str, -1)
	if len(sarr) < 2 {
		return 0, 0, nil, ParseFormatErr
	}
	t := parse1("uid=", sarr[0])
	uid, _, err = parse2(t)
	if err != nil {
		return 0, 0, nil, err
	}
	t = parse1("gid=", sarr[1])
	gid, _, err = parse2(t)
	if err != nil {
		return uid, 0, nil, err
	}
	if len(sarr) < 3 {
		return uid, gid, nil, err
	}
	for _, v := range strings.Split(sarr[2], ",") {
		id, name, err := parse2(v)
		if err != nil {
			continue
		}
		groups = append(groups, Group{GroupName: name, GroupId: id})
	}
	return
}

// 解析uid=1000(pengda) 返回1000(pengda)
func parse1(pre, str string) string {
	sarr := strings.Split(str, pre)
	if len(sarr) < 2 {
		return ""
	}
	return sarr[1]
}

// 解析1000(pengda)，返回1000, pengda
func parse2(str string) (int, string, error) {
	var (
		i int
		s string
	)
	id := extractId.FindString(str)
	if t, err := strconv.Atoi(id); err != nil {
		return 0, "", err
	} else {
		i = t
	}
	t := extractName.FindStringSubmatch(str)
	if len(t) > 2 {
		s = t[1]
	}
	return i, s, nil
}
