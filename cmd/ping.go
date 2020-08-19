package cmd

import (
	"github.com/spf13/cobra"
	"multi_ssh/m_terminal"
	"os"
)

func init() {
	rootCmd.AddCommand(&pingCmd)
}

var pingCmd = cobra.Command{
	Use:   "ping",
	Short: "用于测试主机是否可以连通",
	Args:  cobra.MaximumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		const pingShowFormat = "#{user}@#{host}:{\n\tmsg: #{msg},\n\tcode: #{code}\n}\n"
		ch := make(chan *execResult, 0)
		outFinish := output(ch, pingShowFormat, os.Stdout)
		execFinish := eachTerm(terminals, func(term *m_terminal.Terminal) {
			rst := term.Run(false, "whoami")
			r := buildExecResultFromResult(rst)
			r.u = term.GetUser()
			ch <- r
		})
		<-execFinish
		close(ch)
		<-outFinish
	},
}
