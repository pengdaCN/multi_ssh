package playbook

import (
	lua "github.com/yuin/gopher-lua"
)

var (
	VM *lua.LState
)

func init() {
	VM = lua.NewState()
	VM.Register("shell", shell)
	VM.Register("script", script)
	VM.Register("copy", cp)
}
