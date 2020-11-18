package cmd

import (
	"context"
	"fmt"
	"io"
	"log"
	"multi_ssh/m_terminal"
	"multi_ssh/model"
	"multi_ssh/tools"
	"reflect"
	"strings"
	"sync"
	"time"
)

const defaultOutputFormat = "#{user}@#{host}: {\n#{msg}\n}\n"

var (
	printf  reflect.Value
	outFunc map[string]outAttribute
)

type (
	outAttribute func(*execResult) string
	execResult   struct {
		u       model.SHHUser
		msg     string
		errInfo string
		code    int
	}
	eachFunc func(term *m_terminal.Terminal)
)

func init() {
	outFunc = make(map[string]outAttribute)
	printf = reflect.ValueOf(fmt.Sprintf)
	outRegistry("user", func(result *execResult) string {
		return result.u.User()
	})
	outRegistry("host", func(result *execResult) string {
		return result.u.Host()
	})
	outRegistry("msg", func(result *execResult) string {
		return result.msg
	})
	outRegistry("err", func(result *execResult) string {
		return result.errInfo
	})
	outRegistry("code", func(result *execResult) string {
		return fmt.Sprintf("%d", result.code)
	})
}

func eachTerm(terms []*m_terminal.Terminal, fn eachFunc) chan struct{} {
	finish := make(chan struct{}, 0)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		// 当timeout 设置为-1时，没有任务超时
		if timeout == -1 {
			return
		}
		<-time.NewTimer(timeout).C
		cancel()
	}()
	go func() {
		defer func() {
			finish <- struct{}{}
		}()
		var w sync.WaitGroup
		for i := 0; i < len(terms); i++ {
			if execableNums > 0 && execableNums < i {
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
			}(terms[i])
		}
		w.Wait()
	}()
	return finish
}

func buildExecResultFromResult(r *m_terminal.Result) *execResult {
	rst := new(execResult)
	{
		rst.msg = r.Msg()
		rst.errInfo = r.ErrInfo()
		rst.code = r.Code()
	}
	return rst
}

func outRegistry(key string, val outAttribute) {
	outFunc[key] = val
}

func formatParse(format string) (f string, in []outAttribute) {
	_format := strings.Builder{}
	var stat uint8
	attributeName := strings.Builder{}
	for _, v := range format {
		if v == '#' {
			stat |= 1
			continue
		}
		if v == '{' {
			if stat == 1 {
				stat |= 1 << 1
				continue
			}
			if stat == 3 {
				panic("格式不对")
			}
			_format.WriteRune(v)
			continue
		}
		// 当在获取属性遇到}，直接结束属性名字的获取
		if v == '}' && stat == 3 {
			stat &= 0
			aName := attributeName.String()
			if v, ok := outFunc[aName]; ok {
				in = append(in, v)
			} else {
				panic("错误的属性名")
			}
			_format.WriteString(`%s`)
			attributeName = strings.Builder{}
			continue
		}
		if stat == 3 {
			attributeName.WriteRune(v)
			continue
		}
		_format.WriteRune(v)
	}
	return _format.String(), in
}

func output(in <-chan *execResult, format string, writer ...io.Writer) chan struct{} {
	finish := make(chan struct{}, 0)
	f, fns := formatParse(format)
	{
		f = tools.SpecialStrTransform(f)
	}
	args := make([]reflect.Value, len(fns)+1)
	args[0] = reflect.ValueOf(f)
	go func() {
		defer func() {
			finish <- struct{}{}
		}()
		var w sync.WaitGroup
		var look sync.Mutex
		for v := range in {
			w.Add(1)
			var str string
			{
				for i, fn := range fns {
					args[i+1] = reflect.ValueOf(fn(v))
				}
				t := printf.Call(args)
				str = t[0].String()
			}
			go func(o string) {
				defer w.Done()
				for _, val := range writer {
					if val == nil {
						continue
					}
					look.Lock()
					_, err := fmt.Fprint(val, o)
					//_, err := val.Write(tools.String2ByteSlice(o))
					look.Unlock()
					if err != nil {
						log.Println(err)
					}
				}
			}(str)
		}
		w.Wait()
	}()
	return finish
}
