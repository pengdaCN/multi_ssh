package m_terminal

import (
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
	chg        chan struct{}
	wantAtTime bool
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
	if p.wantAtTime {
		<-p.chg
	}
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

func (p *par) GetLastParByTime(t time.Time) string {
	p.RLock()
	if p.updateTime.Before(t) {
		p.wantAtTime = true
		p.chg <- struct{}{}
	}
	p.RUnlock()
	tmp := p.GetLastPar()
	return tmp
}

type content2 struct {
	out       par
	stdout    par
	stderr    par
	sharePool map[string]interface{}
}

//func GetOutLast() string {
//
//}
