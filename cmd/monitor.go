package cmd

import (
	"bytes"
	"fmt"
	"multi_ssh/extra_mod/playbook"
	"multi_ssh/tools"
	"os"
	"strings"
)

func init() {
	go monitor()
}

const keywords = "show"

func monitor() {
	var (
		msg string
		ch  chan *execResult
	)
	go func() {
		c := output(ch, defaultOutputFormat, os.Stdout)
		<-c
	}()
	for {
		_, _ = fmt.Scanln(&msg)
		m := strings.TrimSpace(msg)
		if m == keywords {
			for _, term := range terminals {
				if o, ok := term.GetOnceShare(playbook.OutKey); ok {
					sb := o.(*strings.Builder)
					str := sb.String()
					msg = str
					ch <- &execResult{
						code: 0,
						u:    term.GetUser(),
						msg:  msg,
					}
					continue
				}
				msg = tail(term.GetMsg(), 10)
				ch <- &execResult{
					code: 0,
					u:    term.GetUser(),
					msg:  msg,
				}
				continue
			}
		}
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
