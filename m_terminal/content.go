package m_terminal

import (
	"multi_ssh/extra_mod/host_info"
	"sync"
	"time"
)

type par struct {
	data       [line]string
	updateTime time.Time
	floor      int
	heap       int
	length     int
	cap        int
	sync.RWMutex
	chg chan struct{}
}

func NewPar() *par {
	return &par{
		length: line,
	}
}

func (p *par) pop(src []byte) {
	str := string(src)
	if p.length != p.cap {
		p.length++
		p.heap++
		p.data[p.heap] = str
		return
	}
	if p.floor+1 != p.cap {
		p.floor++
	} else {
		p.floor = 0
	}
	if p.heap+1 != p.cap {
		p.heap++
	} else {
		p.heap = 0
	}
	p.data[p.heap] = str
}

func (p *par) Write(src []byte) (n int, err error) {
	p.Lock()
	p.updateTime = time.Now()
	p.pop(src)
	p.Unlock()
	return len(src), nil
}

func (p *par) GetLastPar() string {
	p.RLock()
	tmp := p.data[p.heap]
	p.RUnlock()
	return tmp
}

type content struct {
	out       *par
	stdout    *par
	stderr    *par
	stdin     *par
	sharePool map[string]interface{}
}

func NewContent() *content {
	return &content{
		out:       NewPar(),
		stderr:    NewPar(),
		stdout:    NewPar(),
		stdin:     NewPar(),
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
