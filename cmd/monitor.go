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
					println("playbook out")
					//ch <- &execResult{
					//	code: 0,
					//	u:    term.GetUser(),
					//	msg:  msg,
					//}
					fmt.Printf("%s:%s {\n%s\n}", term.GetUser().User(), term.GetUser().Host(), msg)
					println("playbook over")
					continue
				}
				println("get cmd data")
				msg = tail(term.GetMsg(), 10)
				println("cmd out")
				fmt.Printf("%s:%s {\n%s\n}", term.GetUser().User(), term.GetUser().Host(), msg)
				//ch <- &execResult{
				//	code: 0,
				//	u:    term.GetUser(),
				//	msg:  msg,
				//}
				println("cmd out")
				continue
			}
		}
		println(m)
	}
}

func tail(src []byte, n int) string {
	var cache []int
	for {
		i := bytes.IndexRune(src, '\n')
		if i == -1 {
			cache = append(cache, len(src))
		}
		cache = append(cache, i)
		if len(src) == 0 {
			break
		}
	}
	if n > len(cache) {
		return tools.ByteSlice2String(src)
	}
	return tools.ByteSlice2String(src[cache[len(cache)-n-1]:])
}
