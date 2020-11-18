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
	tools.RawSetString("sleep", VM.NewFunction(luaSleep))
	tools.RawSetString("setShareIotaMax", VM.NewFunction(setOnceShareNum))
	tools.RawSetString("getShareIota", VM.NewFunction(getShareNum))
	tools.RawSetString("newWaitGroup", VM.NewFunction(newWaitGroup))
	tools.RawSetString("newTokenBucket", VM.NewFunction(newTokenBucket))
	tools.RawSetString("newMux", VM.NewFunction(newMux))
	tools.RawSetString("newSafeTable", VM.NewFunction(newSafeTable))
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
	VM.SetGlobal("tools", SetReadOnly(VM, tools))
}
