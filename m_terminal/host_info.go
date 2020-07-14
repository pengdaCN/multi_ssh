package m_terminal

import (
	"errors"
	"github.com/tidwall/gjson"
	"log"
	"multi_ssh/extra_mod/host_info"
	"regexp"
	"strings"
)

// 在content sharePool中存放主机信息的key
const HostInfoKey = `host_info`

// 获取主机信息，在sharePool中存放，键为host_info
func GetRemoteHostInfo(term *Terminal) {
	uinfo, err := getRemoteHostUserInfo(term)
	if err != nil {
		log.Printf("%s: 获取用户信息失败\n", term.user.User()+"@"+term.user.Host())
	}
	ninfo, err := getRemoteNetInterinfo(term)
	if err != nil {
		log.Printf("%s: 获取网络接口信息失败\n", term.user.User()+"@"+term.user.Host())
	}
	term.content.sharePool[HostInfoKey] = &host_info.HostGenericInfo{
		User: uinfo,
		Net: ninfo,
	}
}

var (
	enter, _ = regexp.Compile(`\n`)
	GetHostUserInfoErr = errors.New("获取hostUser信息错误")
	GetHostUserGroupErr = errors.New("获取user id与group信息错误")
)

// 获取远程主机的用户信息
func getRemoteHostUserInfo(term *Terminal) (*host_info.UserInfo, error) {
	s, err := term.GetSessionWithTerm()
	if err != nil {
		return nil, err
	}
	rst, err := s.Output(cmdPrefixGeneric + `echo "$USER"; echo "$HOME"; echo "$SHELL"; id -a;`)
	if err != nil {
		return nil, err
	}
	strRst := string(rst)
	strRst = strings.TrimSpace(strRst)
	sarr := enter.Split(strRst, -1)
	if len(sarr) != 4 {
		return nil, GetHostUserInfoErr
	}
	for i := 0; i < len(sarr); i++ {
		sarr[i] = strings.TrimSpace(sarr[i])
	}
	uinfo := host_info.UserInfo{
		UserName: sarr[0],
		Home: sarr[1],
		LoginShell: sarr[2],
	}
	uid, gid, groups, err := host_info.ParseIdAndGroup(sarr[3])
	if err != nil {
		return &uinfo, GetHostUserGroupErr
	}

	{
		uinfo.Uid = uid
		uinfo.Gid = gid
		uinfo.Groups = groups
		if uinfo.Uid == 0 {
			uinfo.IsRoot = true
		}
	}
	return &uinfo, nil
}

func getRemoteNetInterinfo(term *Terminal) (netinters []*host_info.NetInterInfo, err error) {
	s, err := term.GetSessionWithTerm()
	if err != nil {
		return nil, err
	}
	rst, err := s.Output(cmdPrefixGeneric + `ip -j a`)
	if err != nil {
		return nil, err
	}
	j := gjson.ParseBytes(rst)
	jarr := j.Array()
	for _, v := range jarr {
		t, err := host_info.ParseNetInterInfoWithGjson(v)
		if err != nil {
			continue
		}
		netinters = append(netinters, t)
	}
	return
}