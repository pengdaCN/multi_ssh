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
		msg     string
		curInfo string
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
			if terminals == nil {
				continue
			}
			for _, term := range terminals {
				if o, ok := term.GetOnceShare(playbook.OutKey); ok {
					sb := o.(*strings.Builder)
					str := sb.String()
					curInfo = str
					//ch <- &execResult{
					//	code: 0,
					//	u:    term.GetUser(),
					//	msg:  msg,
					//}
					curInfo += fmt.Sprintf("\nshell cur: {\n%s\n}", tail(term.GetMsg(), 10))
					fmt.Printf("%s:%s {\n%s\n}", term.GetUser().User(), term.GetUser().Host(), curInfo)
					continue
				}
				curInfo = tail(term.GetMsg(), 10)
				fmt.Printf("%s:%s {\n%s\n}", term.GetUser().User(), term.GetUser().Host(), curInfo)
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
	if n <= 0 || src == nil {
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
