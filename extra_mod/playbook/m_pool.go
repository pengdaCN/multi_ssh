package playbook

import (
	"multi_ssh/m_terminal"
	"sync"
)

var (
	termPool *TermPool
)

func init() {
	termPool = new(TermPool)
}

type TermPool struct {
	mPool sync.Map
}

func (t *TermPool) Push(id int, term *m_terminal.Terminal) {
	t.mPool.Store(id, term)
}

func (t *TermPool) Get(id int) (m *m_terminal.Terminal, ok bool) {
	if v, _ok := t.mPool.Load(id); _ok {
		m = v.(*m_terminal.Terminal)
		ok = _ok
	}
	return
}

func SetTermPool(t *TermPool) {
	termPool = t
}

func Push(id int, t *m_terminal.Terminal) {
	termPool.Push(id, t)
}

func Get(id int) (*m_terminal.Terminal, bool) {
	return termPool.Get(id)
}
