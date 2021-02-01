package playbook

import (
	lua "github.com/yuin/gopher-lua"
	"multi_ssh/common"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"
	"text/template"
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
	var (
		str  string
		re   string
		mode string
	)
	str = state.ToString(1)
	re = state.ToString(2)
	mode = state.ToString(3)
	_re := common.GetRe(re)
	if mode == "" || mode == lua.LTNil.String() {
		// 为兼容之前的版本
		mode = "sub"
	}
	switch mode {
	case "sub":
		arr := state.NewTable()
		_arr := _re.FindStringSubmatch(str)
		strSliceToTable(arr, _arr)
		state.Push(arr)
	case "sub_all":
		arr := state.NewTable()
		_arr := _re.FindAllStringSubmatch(str, -1)
		for i, v := range _arr {
			a := state.NewTable()
			strSliceToTable(a, v)
			arr.Insert(i+1, a)
		}
		state.Push(arr)
	case "str":
		r := _re.FindString(str)
		state.Push(lua.LString(r))
	case "str_all":
		arr := state.NewTable()
		_arr := _re.FindAllString(str, -1)
		strSliceToTable(arr, _arr)
		state.Push(arr)
	default:
		r := _re.FindString(str)
		state.Push(lua.LString(r))
	}
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
	_re := common.GetRe(re)
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
	_re := common.GetRe(re)
	_newStr := _re.ReplaceAllString(str, newStr)
	state.Push(lua.LString(_newStr))
	return 1
}

var (
	shareN int32
)

func setOnceShareNum(state *lua.LState) int {
	i := state.ToInt(1)
	if i < 0 {
		state.Push(lua.LFalse)
	}
	if atomic.LoadInt32(&shareN) > 0 {
		state.Push(lua.LFalse)
	}
	atomic.SwapInt32(&shareN, int32(i))
	state.Push(lua.LTrue)
	return 1
}

func getShareNum(state *lua.LState) int {
START:
	cur := atomic.LoadInt32(&shareN)
	if cur <= 0 {
		state.Push(lua.LNumber(0))
		return 1
	}
	if !atomic.CompareAndSwapInt32(&shareN, cur, cur-1) {
		goto START
	}
	state.Push(lua.LNumber(cur))
	return 1
}

func newMux(state *lua.LState) int {
	var rwLock sync.RWMutex
	tb := state.NewTable()
	lock := func(lState *lua.LState) int {
		rwLock.Lock()
		return 0
	}
	rLock := func(lState *lua.LState) int {
		rwLock.RLock()
		return 0
	}
	unLock := func(lState *lua.LState) int {
		rwLock.Unlock()
		return 0
	}
	rUnlock := func(lState *lua.LState) int {
		rwLock.RUnlock()
		return 0
	}
	tb.RawSetString("lock", state.NewFunction(lock))
	tb.RawSetString("rLock", state.NewFunction(rLock))
	tb.RawSetString("unLock", state.NewFunction(unLock))
	tb.RawSetString("rUnlock", state.NewFunction(rUnlock))
	state.Push(SetReadOnly(state, tb))
	return 1
}

func newWaitGroup(state *lua.LState) int {
	gw := state.NewTable()
	var w sync.WaitGroup
	// add func
	add := func(lState *lua.LState) int {
		i := lState.ToInt(1)
		if i < 0 {
			panic("add value cannot zero")
		}
		w.Add(i)
		return 0
	}
	// done func
	done := func(lState *lua.LState) int {
		w.Done()
		return 0
	}
	// wait func
	wait := func(lState *lua.LState) int {
		w.Wait()
		return 0
	}
	gw.RawSetString("add", state.NewFunction(add))
	gw.RawSetString("done", state.NewFunction(done))
	gw.RawSetString("wait", state.NewFunction(wait))
	state.Push(SetReadOnly(state, gw))
	return 1
}

func newTokenBucket(state *lua.LState) int {
	i := int32(state.ToInt(1))
	if i < 0 {
		panic("bucket size cannot zero")
	}
	t := state.NewTable()
	get := func(lState *lua.LState) int {
	START:
		cur := atomic.LoadInt32(&i)
		if cur <= 0 {
			state.Push(lua.LNumber(0))
			return 1
		}
		if !atomic.CompareAndSwapInt32(&i, cur, cur-1) {
			goto START
		}
		state.Push(lua.LNumber(cur))
		return 1
	}
	t.RawSetString("get", state.NewFunction(get))
	state.Push(SetReadOnly(state, t))
	return 1
}

type safeTable struct {
	mu sync.RWMutex
	tb *lua.LTable
}

func (s *safeTable) append(value lua.LValue) {
	s.mu.Lock()
	s.tb.Append(value)
	s.mu.Unlock()
}

func (s *safeTable) set(key, val lua.LValue) {
	s.mu.Lock()
	s.tb.RawSet(key, val)
	s.mu.Unlock()
	return
}

func (s *safeTable) len() int {
	var length int
	s.mu.RLock()
	length = s.tb.Len()
	s.mu.RUnlock()
	return length
}

func (s *safeTable) get(key lua.LValue) lua.LValue {
	var val lua.LValue
	s.mu.RLock()
	val = s.tb.RawGet(key)
	s.mu.RUnlock()
	return val
}

func (s *safeTable) rLock() {
	s.mu.RLock()
}

func (s *safeTable) rUnlock() {
	s.mu.RUnlock()
}

func (s *safeTable) into() *lua.LTable {
	return s.tb
}

func newSafeTable(state *lua.LState) int {
	st := new(safeTable)
	st.tb = state.NewTable()
	lAppend := func(lState *lua.LState) int {
		val := lState.Get(1)
		st.append(val)
		return 0
	}
	lSet := func(lState *lua.LState) int {
		key, val := lState.Get(1), lState.Get(2)
		st.set(key, val)
		return 0
	}
	lGet := func(lState *lua.LState) int {
		key := lState.Get(1)
		val := st.get(key)
		lState.Push(val)
		return 1
	}
	lLen := func(lState *lua.LState) int {
		val := st.len()
		lState.Push(lua.LNumber(val))
		return 1
	}
	rLock := func(lState *lua.LState) int {
		st.rLock()
		return 0
	}
	rUnlock := func(lState *lua.LState) int {
		st.rUnlock()
		return 0
	}
	into := func(lState *lua.LState) int {
		t := st.into()
		lState.Push(t)
		return 1
	}
	tb := state.NewTable()
	tb.RawSetString("append", state.NewFunction(lAppend))
	tb.RawSetString("set", state.NewFunction(lSet))
	tb.RawSetString("get", state.NewFunction(lGet))
	tb.RawSetString("len", state.NewFunction(lLen))
	tb.RawSetString("rLock", state.NewFunction(rLock))
	tb.RawSetString("rUnlock", state.NewFunction(rUnlock))
	tb.RawSetString("into", state.NewFunction(into))
	state.Push(SetReadOnly(state, tb))
	return 1
}

func newOnce(state *lua.LState) int {
	var once sync.Once
	do := func(lState *lua.LState) int {
		val := lState.Get(1)
		fn, ok := val.(*lua.LFunction)
		if !ok {
			lState.RaiseError("need a function")
			return 0
		}
		once.Do(func() {
			_ = lState.CallByParam(lua.P{
				Fn:      fn,
				NRet:    0,
				Protect: true,
			})
		})
		return 0
	}
	tb := state.NewTable()
	tb.RawSetString("Do", state.NewFunction(do))
	state.Push(SetReadOnly(state, tb))
	return 1
}

func newTmpl(state *lua.LState) int {
	tmpl := template.New("MULTI_SSH_TMPL")
	tb := state.NewTable()
	parse := func(lState *lua.LState) int {
		_name := lState.ToString(1)
		if _name == "" {
			state.RaiseError("new template must give a name")
			return 0
		}
		text := lState.ToString(2)
		if text == "" {
			lState.RaiseError("context must not be empty")
			return 0
		}
		var err error
		tmpl, err = tmpl.New(_name).Parse(text)
		if err != nil {
			lState.RaiseError(err.Error())
			return 0
		}
		return 0
	}
	execute := func(lState *lua.LState) int {
		var (
			name string
			val  interface{}
		)
		name = lState.ToString(1)
		if name == "" {
			state.RaiseError("name must not be empty")
			return 0
		}
		lval := lState.Get(2)
		val = LuaValueToGoVal(lval)
		if val == nil {
			lState.RaiseError("nonsupport type")
			return 0
		}
		var sb strings.Builder

		err := tmpl.ExecuteTemplate(&sb, name, val)
		if err != nil {
			lState.RaiseError(err.Error())
			return 0
		}
		lState.Push(lua.LString(sb.String()))
		return 1
	}
	tb.RawSetString("parse", state.NewFunction(parse))
	tb.RawSetString("execute", state.NewFunction(execute))
	state.Push(SetReadOnly(state, tb))
	return 1
}
