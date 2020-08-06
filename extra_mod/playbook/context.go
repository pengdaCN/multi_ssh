package playbook

import (
	lua "github.com/yuin/gopher-lua"
	"multi_ssh/m_terminal"
	"os"
)

var (
	_context *Context
)

type Context struct {
	term []m_terminal.Terminal
	out  _output
	exec *lua.LFunction
}

func SetCurrentContext(ctx *Context) {
	_context = ctx
}

func AppendTerm(term *m_terminal.Terminal) {
	_context.AppendTerm(term)
}

func GetTerm(id int) *m_terminal.Terminal {
	return _context.GetTerm(id)
}

func (c *Context) AppendTerm(term *m_terminal.Terminal) {
	if c.term == nil {
		c.term = []m_terminal.Terminal{*term}
		return
	}
	c.term = append(c.term, *term)
}

func newContext(term []m_terminal.Terminal, exec *lua.LFunction) *Context {
	return &Context{
		term: term,
		out:  newOutput(os.Stdin),
		exec: exec,
	}
}

func (c *Context) GetTerm(id int) *m_terminal.Terminal {
	return &c.term[id]
}
