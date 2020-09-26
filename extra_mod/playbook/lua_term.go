package playbook

import (
	"context"
	lua "github.com/yuin/gopher-lua"
	"multi_ssh/m_terminal"
	"net"
)

func NewLuaTerm(state *lua.LState, term *m_terminal.Terminal, cancel context.CancelFunc) *lua.LTable {
	tab := state.NewTable()
	state.SetField(tab, "shell", state.NewFunction(newShell(term)))
	state.SetField(tab, "script", state.NewFunction(newScript(term)))
	state.SetField(tab, "copy", state.NewFunction(newCopy(term)))
	state.SetField(tab, "out", state.NewFunction(newOut(term)))
	state.SetField(tab, "outln", state.NewFunction(newOutLn(term)))
	state.SetField(tab, "extraInfo", state.NewFunction(newExtra(term)))
	state.SetField(tab, "hostInfo", state.NewFunction(newHostInfo(term)))
	state.SetField(tab, "setCode", state.NewFunction(newSetCode(term)))
	state.SetField(tab, "setErrInfo", state.NewFunction(newSetErrInfo(term)))
	state.SetField(tab, "sleep", state.NewFunction(luaSleep))
	state.SetField(tab, "hostInfo", initHostInfo(state, term))
	state.SetField(tab, "iota", lua.LNumber(term.GetBirthID()))
	state.SetField(tab, "exit", state.NewFunction(newExit(cancel)))
	return tab
}

func initHostInfo(state *lua.LState, term *m_terminal.Terminal) *lua.LTable {
	hostInfo := state.NewTable()
	s := term.GetUser()
	if s == nil {
		return state.NewTable()
	}
	hostInfo.RawSetString("line", lua.LNumber(s.Line()))
	ip, port, err := net.SplitHostPort(s.Host())
	if err == nil {
		hostInfo.RawSetString("ip", lua.LString(ip))
		hostInfo.RawSetString("port", lua.LString(port))
	}
	hostInfo.RawSetString("user", lua.LString(s.User()))
	e := s.Extra()
	if e != nil {
		hostInfo.RawSetString("extra", mapToLTable(state, e))
	}
	return hostInfo
}
