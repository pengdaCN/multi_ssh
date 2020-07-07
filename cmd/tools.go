package cmd

import (
	"fmt"
	"io"
	"log"
)

const defaultOutputFormat = "%s@%s:{%s}\n"

func outputByFormat(format string, result *commandResult, out ...io.Writer)  {
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
