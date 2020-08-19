package m_terminal

import (
	"multi_ssh/extra_mod/host_info"
	"strings"
)

type content struct {
	out       strings.Builder
	sharePool map[string]interface{}
}

func NewContent() *content {
	return &content{
		sharePool: make(map[string]interface{}),
	}
}

func (c *content) GetHostInfo() (*host_info.HostGenericInfo, bool) {
	if v, ok := c.sharePool[HostInfoKey]; ok {
		if info, sok := v.(*host_info.HostGenericInfo); sok {
			return info, true
		} else {
			return nil, false
		}
	}
	return nil, false
}
