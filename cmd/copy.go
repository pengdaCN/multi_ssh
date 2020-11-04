package cmd

import (
	"github.com/spf13/cobra"
	"multi_ssh/m_terminal"
	"os"
)

var (
	//mode       string
	//uid        int
	//gid        int
	copySudo   bool
	copyExists bool
)

func init() {
	rootCmd.AddCommand(&copyCmd)
	//copyCmd.Flags().StringVar(&mode, "mode", "", "设置文件上传后的权限")
	//copyCmd.Flags().IntVar(&uid, "uid", -1, "设置上传后文件uid")
	//copyCmd.Flags().IntVar(&gid, "gid", -1, "设置上传后文件的gid")
	copyCmd.Flags().BoolVarP(&copySudo, "sudo", "S", false, "可将本地文件无限制的拷贝到远端")
	copyCmd.Flags().BoolVarP(&copyExists, "exists", "e", false, "当远程目录不存在则创建")
}

var copyCmd = cobra.Command{
	Use:   "copy src dst",
	Short: "copy 命令将本地的文件拷贝到远端",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		srcPaths := args[:len(args)-1]
		dstPath := args[len(args)-1]
		ch := make(chan *execResult, 0)
		outFinish := output(ch, outFormat, os.Stdout)
		execFinish := eachTerm(terminals, func(term *m_terminal.Terminal) {
			rst := term.Copy(copyExists, copySudo, srcPaths, dstPath)
			term.CfgStat()
			r := buildExecResultFromResult(rst)
			r.u = term.GetUser()
			ch <- r
		})
		<-execFinish
		close(ch)
		<-outFinish
	},
}
