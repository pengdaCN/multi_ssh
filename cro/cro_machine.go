package cro

import (
	"bytes"
	"context"
	"github.com/pkg/errors"
	lua "github.com/yuin/gopher-lua"
	"io"
	"io/ioutil"
	"log"
	"multi_ssh/extra_mod/playbook"
	"multi_ssh/m_terminal"
	"multi_ssh/model"
	"multi_ssh/tools"
	"os"
	"sort"
	"strings"
	"sync"
	"time"
)

type userSlice []model.SHHUser

func (u userSlice) Less(v1, v2 int) bool {
	return u[v1].Line() < u[v2].Line()
}

func (u userSlice) Swap(v1, v2 int) {
	u[v1], u[v2] = u[v2], u[v1]
}

func (u userSlice) Len() int {
	return len(u)
}

type baseRunEnv struct {
	conf           *baseTaskBuilder
	users          userSlice
	outFormat      string
	filter string
	outSite        io.Writer
	execTimeout    time.Duration
	maxExecSeveral int
	terms          []*m_terminal.Terminal
}

func getBaseRunEnvFromBaseBuilder(b *baseTaskBuilder) (br *baseRunEnv, err error) {
	br = new(baseRunEnv)
	if b.rawHostsInfo != "" {
		us, err := model.ReadLines(b.rawHostsInfo)
		if err != nil {
			return nil, err
		}
		for _, v := range us {
			br.users = append(br.users, v)
		}
	} else {
		us, err := model.ReadHosts(b.hostsF)
		if err != nil {
			return nil, err
		}
		for _, v := range us {
			br.users = append(br.users, v)
		}
	}
	sort.Sort(br.users)
	br.execTimeout = b.timeout
	br.outFormat = b.format
	br.filter = b.filerStr
	if br.outFormat == "" {
		br.outFormat = DefaultOutputFormat
	}
	if b.out == nil {
		br.outSite = os.Stdout
	}
	br.conf = b
	return
}

func (b *baseRunEnv) ready() {
	if b.filter != "" {
		b.users = filters(b.users, b.filter)
	}
	ch := make(chan *m_terminal.Terminal, 0)
	if b.terms != nil {
		b.terms = b.terms[:0]
	}
	var w sync.WaitGroup
	for i, u := range b.users {
		w.Add(1)
		go func(user model.SHHUser, bi int) {
			defer w.Done()
			c, err := m_terminal.DefaultWithPassphrase(user)
			if err != nil {
				log.Printf("打开%s失败 %s", user.Host(), err)
				return
			} else {
				log.Printf("打开%s成功", user.Host())
			}
			if b.conf.preInfo {
				m_terminal.GetRemoteHostInfo(c)
			}
			c.SetBirthID(bi + 1)
			ch <- c
		}(u, i)
	}
	var w2 sync.WaitGroup
	w2.Add(1)
	go func() {
		defer w2.Done()
		for i := range ch {
			b.terms = append(b.terms, i)
		}
	}()
	w.Wait()
	close(ch)
	w2.Wait()
	go b.monitor()
}

func (b *baseRunEnv) runEach(fn func(m *m_terminal.Terminal)) chan struct{} {
	b.ready()
	finish := make(chan struct{}, 0)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		// 当timeout 设置为-1时，没有任务超时
		if b.execTimeout == -1 {
			return
		}
		<-time.NewTimer(b.execTimeout).C
		cancel()
	}()
	go func() {
		defer func() {
			finish <- struct{}{}
		}()
		var w sync.WaitGroup
		for i := 0; i < len(b.terms); i++ {
			if b.maxExecSeveral > 0 && b.maxExecSeveral < i {
				continue
			}
			w.Add(1)
			go func(term *m_terminal.Terminal) {
				defer w.Done()
				ch := make(chan struct{}, 0)
				go func() {
					err := tools.PanicToErr(func() {
						fn(term)
					})
					if err != nil {
						log.Println(err)
					}
					ch <- struct{}{}
				}()
				// 设置任务超时
				select {
				case <-ch:
				case <-ctx.Done():
					log.Printf("Host: %s timeout", term.GetUser().Host())
				}
			}(b.terms[i])
		}
		w.Wait()
	}()
	return finish
}

type shellRunEnv struct {
	conf *shellTBuilder
	b    *baseRunEnv
}

func (s *shellRunEnv) Run() {
	ch := make(chan *execResult, 0)
	outFinish := output(ch, s.b.outFormat, s.b.outSite)
	execFinish := s.b.runEach(func(term *m_terminal.Terminal) {
		rst := term.Run(s.conf.sudo, s.conf.cmds)
		term.CfgStat()
		r := buildExecResultFromResult(rst)
		r.u = term.GetUser()
		ch <- r
	})
	<-execFinish
	close(ch)
	<-outFinish
}

type ScriptEunEnv struct {
	conf *scriptTBuilder
	b    *baseRunEnv
}

func (s *ScriptEunEnv) Run() {
	ch := make(chan *execResult, 0)
	outFinish := output(ch, s.b.outFormat, s.b.outSite)
	var scriptReader io.Reader
	if s.conf.text != "" {
		scriptReader = strings.NewReader(s.conf.text)
	} else {
		f, err := ioutil.ReadFile(s.conf.path)
		if err != nil {
			panic(err)
		}
		scriptReader = bytes.NewBuffer(f)
	}
	execFinish := s.b.runEach(func(term *m_terminal.Terminal) {
		rst := term.Script(s.conf.sudo, scriptReader, s.conf.args)
		term.CfgStat()
		r := buildExecResultFromResult(rst)
		r.u = term.GetUser()
		ch <- r
	})
	<-execFinish
	close(ch)
	<-outFinish
}

type playbookTRunEnv struct {
	conf *playbookTBuilder
	b    *baseRunEnv
}

func (s *playbookTRunEnv) Run() {
	ch := make(chan *execResult, 0)
	outFinish := output(ch, s.b.outFormat, s.b.outSite)
	if s.conf.text != "" {
		if err := playbook.VM.DoString(s.conf.text); err != nil {
			log.Println(errors.WithStack(err))
			return
		}
	} else {
		if err := playbook.VM.DoFile(s.conf.path); err != nil {
			log.Println(errors.WithStack(err))
			return
		}
	}
	if s.conf.vars != "" {
		setGlobalVal(s.conf.vars)
	}
	var (
		fn *lua.LFunction
		ok bool
	)
	_ = playbook.VM.CallByParam(lua.P{
		Fn:      playbook.VM.GetGlobal("BEGIN"),
		NRet:    0,
		Protect: true,
	})
	if fn, ok = playbook.VM.GetGlobal("exec").(*lua.LFunction); !ok {
		log.Println(errors.New("未读取到exec函数，请检查代码"))
		return
	}
	m := playbook.VM.NewTable()
	playbook.VM.SetGlobal("m", m)
	m.RawSetString("hosts_num", lua.LNumber(len(s.b.terms)))
	finished := s.b.runEach(func(term *m_terminal.Terminal) {
		// 执行begin
		beginCo, beginCancel := playbook.VM.NewThread()
		beginT := playbook.NewLuaTerm(beginCo, term, beginCancel)
		_ = beginCo.CallByParam(lua.P{
			Fn:      playbook.VM.GetGlobal("EXEC_BEGIN"),
			NRet:    0,
			Protect: true,
		}, beginT)

		// 执行 exec
		co, cancel := playbook.VM.NewThread()
		t := playbook.NewLuaTerm(co, term, cancel)
		_, err, _ := playbook.VM.Resume(co, fn, t)
		if err != nil {
			log.Println("exec : ", err.Error())
		}

		// 执行over
		overCo, overCancel := playbook.VM.NewThread()
		overT := playbook.NewLuaTerm(beginCo, term, overCancel)
		_ = overCo.CallByParam(lua.P{
			Fn:      playbook.VM.GetGlobal("EXEC_OVER"),
			NRet:    0,
			Protect: true,
		}, overT)
		term.CfgStat()
		var (
			msg     string
			code    int
			errInfo string
		)
		out, ok := term.GetOnceShare(playbook.OutKey)
		if ok {
			sb := out.(*strings.Builder)
			str := sb.String()
			msg = str
		}
		c, ok := term.GetOnceShare(playbook.Code)
		if ok {
			code, _ = c.(int)
		}
		_errInfo, ok := term.GetOnceShare(playbook.ErrInfo)
		if ok {
			errInfo, _ = _errInfo.(string)
		}
		rst := new(execResult)
		{
			rst.errInfo = errInfo
			rst.msg = msg
			rst.code = code
			rst.u = term.GetUser()
		}
		ch <- rst
	})
	_ = playbook.VM.CallByParam(lua.P{
		Fn:      playbook.VM.GetGlobal("OVER"),
		NRet:    0,
		Protect: true,
	})
	<-finished
	close(ch)
	<-outFinish
}

type copyTRunEnv struct {
	conf *copyTBuilder
	b    *baseRunEnv
}

func (s *copyTRunEnv) Run() {
	srcPaths := s.conf.src
	dstPath := s.conf.dst
	ch := make(chan *execResult, 0)
	outFinish := output(ch, s.b.outFormat, s.b.outSite)
	execFinish := s.b.runEach(func(term *m_terminal.Terminal) {
		rst := term.Copy(s.conf.exists, s.conf.sudo, srcPaths, dstPath)
		term.CfgStat()
		r := buildExecResultFromResult(rst)
		r.u = term.GetUser()
		ch <- r
	})
	<-execFinish
	close(ch)
	<-outFinish
}

type pingTRunEnv struct {
	b *baseRunEnv
}

func (s *pingTRunEnv) Run() {
	const pingShowFormat = "#{user}@#{host}:{\n\tmsg: #{msg},\n\tcode: #{code}\n}\n"
	ch := make(chan *execResult, 0)
	outFinish := output(ch, pingShowFormat, os.Stdout)
	execFinish := s.b.runEach(func(term *m_terminal.Terminal) {
		rst := term.Run(false, "whoami")
		r := buildExecResultFromResult(rst)
		r.u = term.GetUser()
		ch <- r
	})
	<-execFinish
	close(ch)
	<-outFinish
}
