package cro

import (
	"bytes"
	"fmt"
	"multi_ssh/extra_mod/playbook"
	"multi_ssh/tools"
	"strings"
)

const keywords = "show"

func (s *baseRunEnv) monitor() {
	var (
		msg     string
		curInfo string
	)
	for {
		_, _ = fmt.Scanln(&msg)
		m := strings.TrimSpace(msg)
		if m == keywords {
			if s.terms == nil {
				continue
			}
			for _, term := range s.terms {
				if term.GetTaskStat() {
					continue
				}
				if o, ok := term.GetOnceShare(playbook.OutKey); ok {
					sb := o.(*strings.Builder)
					str := sb.String()
					curInfo = str
					curInfo += fmt.Sprintf("\nshell cur: {\n%s\n}", tail(term.GetMsg(), 10))
					fmt.Printf("%s:%s {\n%s\n}", term.GetUser().User(), term.GetUser().Host(), curInfo)
					continue
				}
				curInfo = tail(term.GetMsg(), 10)
				fmt.Printf("%s:%s {\n%s\n}", term.GetUser().User(), term.GetUser().Host(), curInfo)
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
