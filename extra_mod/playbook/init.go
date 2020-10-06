package playbook

import (
	lua "github.com/yuin/gopher-lua"
)

var (
	VM *lua.LState
)

func init() {
	VM = lua.NewState()
	tools := VM.NewTable()
	VM.SetGlobal("tools", tools)
	tools.RawSetString("sleep", VM.NewFunction(luaSleep))
	{
		str := VM.NewTable()
		tools.RawSetString("str", str)
		initStr(VM, str)
	}
	{
		re := VM.NewTable()
		tools.RawSetString("re", re)
		initRe(VM, re)
	}
}
