package cmd

import (
	"bytes"
	"fmt"
	"multi_ssh/extra_mod/playbook"
	"multi_ssh/tools"
	"strings"
)

func init() {
	go monitor()
}

const keywords = "show"

func monitor() {
	var (
		msg string
		//ch  chan *execResult
	)
	//go func() {
	//	c := output(ch, defaultOutputFormat, os.Stdout)
	//	<-c
	//}()
	for {
		_, _ = fmt.Scanln(&msg)
		m := strings.TrimSpace(msg)
		if m == keywords {
			println("是关键字")
			for i, term := range terminals {
				fmt.Printf("开始读取: %d", i)
				if o, ok := term.GetOnceShare(playbook.OutKey); ok {
					sb := o.(*strings.Builder)
					str := sb.String()
					msg = str
					//ch <- &execResult{
					//	code: 0,
					//	u:    term.GetUser(),
					//	msg:  msg,
					//}
					fmt.Printf("%s:%s {\n%s\n}", term.GetUser().User(), term.GetUser().Host(), msg)
					continue
				}
				msg = tail(term.GetMsg(), 10)
				fmt.Printf("%s:%s {\n%s\n}", term.GetUser().User(), term.GetUser().Host(), msg)
				//ch <- &execResult{
				//	code: 0,
				//	u:    term.GetUser(),
				//	msg:  msg,
				//}
				continue
			}
		}
		println(m)
	}
}

func tail(src []byte, n int) string {
	if n <= 0 {
		return ""
	}
	cache := []int{0}
	_src := src
	for {
		i := bytes.IndexRune(_src, '\n')
		if i == -1 {
			break
		}
		cache = append(cache, cache[len(cache)-1]+i+1)
		_src = _src[i+1:]
	}
	if n > len(cache) {
		return tools.ByteSlice2String(src)
	}
	return tools.ByteSlice2String(src[cache[len(cache)-(n-1)-1]:])
}
