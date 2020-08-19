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
	VM.Register("extraInfo", extraInfo)
	VM.Register("hostInfo", hostInfo)
	VM.Register("out", out)
	VM.Register("setCode", setCode)
	VM.Register("setErrInfo", setErrInfo)
}
