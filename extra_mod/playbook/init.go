package playbook

import (
	lua "github.com/yuin/gopher-lua"
	"multi_ssh/m_terminal"
)

var (
	VM *lua.LState
)

func init() {
	VM = lua.NewState()
	VM.Register("shell", shell)
	VM.Register("script", script)
	VM.Register("copy", cp)
	_context = newContext(make([]m_terminal.Terminal, 0), nil)
}
