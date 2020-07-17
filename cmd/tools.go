package cmd

import (
	"fmt"
	"io"
	"log"
	"multi_ssh/tools"
	"sync"
)

const defaultOutputFormat = "%s@%s:{%s}\n"

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

func formatParse(format string, data *commandResult) string {
	return fmt.Sprintf(defaultOutputFormat, data.u.Host(), data.u.User(), data.msg)
}

func output(in <-chan *commandResult, format string, writer ...io.Writer) chan struct{} {
	finish := make(chan struct{}, 0)
	go func() {
		defer func() {
			finish <- struct{}{}
		}()
		var w sync.WaitGroup
		var look sync.Mutex
		for v := range in {
			w.Add(1)
			str := formatParse(format, v)
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
