package playbook

import (
	lua "github.com/yuin/gopher-lua"
)

var (
	VM *lua.LState
)

func init() {
	VM = lua.NewState()
}
