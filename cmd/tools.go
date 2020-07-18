package cmd

import (
	"fmt"
	"io"
	"log"
	"multi_ssh/model"
	"multi_ssh/tools"
	"reflect"
	"strings"
	"sync"
)

const defaultOutputFormat = "#{user}@#{host}: {#{msg}}\n"

var printf reflect.Value

type outAttribute func(*commandResult) string

var outFunc map[string]outAttribute

type commandResult struct {
	u   model.SHHUser
	msg []byte
}

func init() {
	outFunc = make(map[string]outAttribute)
	printf = reflect.ValueOf(fmt.Sprintf)
	outRegistry("user", func(result *commandResult) string {
		return result.u.User()
	})
	outRegistry("host", func(result *commandResult) string {
		return result.u.Host()
	})
	outRegistry("msg", func(result *commandResult) string {
		return tools.ByteSlice2String(result.msg)
	})
}

func outputByFormat(format string, result *commandResult, out ...io.Writer) {
	rst := fmt.Sprintf(format, result.u.User(), result.u.Host(), string(result.msg))
	for _, o := range out {
		if o == nil {
			continue
		}
		_, err := o.Write([]byte(rst))
		if err != nil {
			log.Println(err)
		}
	}
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

func output(in <-chan *commandResult, format string, writer ...io.Writer) chan struct{} {
	finish := make(chan struct{}, 0)
	f, funcs := formatParse(format)
	args := make([]reflect.Value, 0, len(funcs)+1)
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
				for i, fn := range funcs {
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
					_, err := val.Write(tools.String2ByteSlice(o))
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
