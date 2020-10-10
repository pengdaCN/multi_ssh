package playbook

import (
	lua "github.com/yuin/gopher-lua"
	"regexp"
	"strings"
)

func initStr(state *lua.LState, table *lua.LTable) {
	table.RawSetString("split", state.NewFunction(strSplit))
	table.RawSetString("hasPrefix", state.NewFunction(strHasPrefix))
	table.RawSetString("hasSuffix", state.NewFunction(strHasSuffix))
	table.RawSetString("trim", state.NewFunction(strTrimSpace))
	table.RawSetString("replace", state.NewFunction(strReplace))
	table.RawSetString("contain", state.NewFunction(strContain))
}

func strContain(state *lua.LState) int {
	b := lua.LFalse
	defer func() {
		state.Push(b)
	}()
	var (
		str string
		sub string
	)
	str = state.ToString(1)
	{
		val := state.Get(2)
		switch val.Type() {
		case lua.LTNil:
			sub = ""
		case lua.LTString:
			sub = val.String()
		default:
			panic("ERROR require str")
		}
	}
	b = lua.LBool(strings.Contains(str, sub))
	return 1
}

func strSplit(state *lua.LState) int {
	arr := state.NewTable()
	defer func() {
		state.Push(arr)
	}()
	var (
		str string
		sep string
	)
	str = state.ToString(1)
	{
		val := state.Get(2)
		switch val.Type() {
		case lua.LTNil:
			str = " "
		case lua.LTString:
			str = val.String()
		}
	}
	_arr := strings.Split(str, sep)
	strSliceToTable(arr, _arr)
	return 1
}

func strHasPrefix(state *lua.LState) int {
	var (
		str    string
		prefix string
	)
	str = state.ToString(1)
	prefix = state.ToString(2)
	b := strings.HasPrefix(str, prefix)
	state.Push(lua.LBool(b))
	return 1
}

func strHasSuffix(state *lua.LState) int {
	var (
		str    string
		suffix string
	)
	str = state.ToString(1)
	suffix = state.ToString(2)
	b := strings.HasSuffix(str, suffix)
	state.Push(lua.LBool(b))
	return 1
}

func strTrimSpace(state *lua.LState) int {
	var (
		str string
	)
	str = state.ToString(1)
	newStr := strings.TrimSpace(str)
	state.Push(lua.LString(newStr))
	return 1
}

func strReplace(state *lua.LState) int {
	var (
		str   string
		old   string
		n     string
		count int
	)
	str = state.ToString(1)
	old = state.ToString(2)
	n = state.ToString(3)
	{
		val := state.Get(4)
		switch val.Type() {
		case lua.LTNil:
			count = -1
		case lua.LTNumber:
			m := val.(lua.LNumber)
			count = int(m)
		}
	}
	_newStr := strings.Replace(str, old, n, count)
	state.Push(lua.LString(_newStr))
	return 1
}

func initRe(state *lua.LState, table *lua.LTable) {
	table.RawSetString("match", state.NewFunction(reMatch))
	table.RawSetString("find", state.NewFunction(reFind))
	table.RawSetString("replace", state.NewFunction(reReplace))
	table.RawSetString("split", state.NewFunction(reSplit))
	table.RawSetString("splitSpace", state.NewFunction(reSplitSpace))
}

func reMatch(state *lua.LState) int {
	var (
		str string
		re  string
	)
	str = state.ToString(1)
	{
		val := state.Get(2)
		switch val.Type() {
		case lua.LTNil:
			state.Push(lua.LFalse)
			return 1
		case lua.LTString:
			re = val.String()
		}
	}
	b, err := regexp.MatchString(re, str)
	if err != nil {
		state.Error(lua.LString(err.Error()), 1)
		return 0
	}
	state.Push(lua.LBool(b))
	return 1
}

func reFind(state *lua.LState) int {
	arr := state.NewTable()
	defer func() {
		state.Push(arr)
	}()
	var (
		str string
		re  string
	)
	str = state.ToString(1)
	re = state.ToString(2)
	_re := regexp.MustCompile(re)
	_arr := _re.FindStringSubmatch(str)
	strSliceToTable(arr, _arr)
	return 1
}

func reSplit(state *lua.LState) int {
	arr := state.NewTable()
	defer func() {
		state.Push(arr)
	}()
	var (
		str string
		re  string
	)
	str = state.ToString(1)
	re = state.ToString(2)
	_re := regexp.MustCompile(re)
	_arr := _re.Split(str, -1)
	strSliceToTable(arr, _arr)
	return 1
}

var (
	space = regexp.MustCompile(`\s+`)
)

func reSplitSpace(state *lua.LState) int {
	arr := state.NewTable()
	defer func() {
		state.Push(arr)
	}()
	str := state.ToString(1)
	_arr := space.Split(str, -1)
	strSliceToTable(arr, _arr)
	return 1
}

func reReplace(state *lua.LState) int {
	var (
		str    string
		re     string
		newStr string
	)
	str = state.ToString(1)
	re = state.ToString(2)
	newStr = state.ToString(3)
	_re := regexp.MustCompile(re)
	_newStr := _re.ReplaceAllString(str, newStr)
	state.Push(lua.LString(_newStr))
	return 1
}
