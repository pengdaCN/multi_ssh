package cmd

import (
	"github.com/spf13/cobra"
)

var (
	copySudo   bool
	copyExists bool
)

func init() {
	rootCmd.AddCommand(&copyCmd)
	copyCmd.Flags().BoolVarP(&copySudo, "sudo", "S", false, "可将本地文件无限制的拷贝到远端")
	copyCmd.Flags().BoolVarP(&copyExists, "exists", "e", false, "当远程目录不存在则创建")
}

var copyCmd = cobra.Command{
	Use:   "copy src dst",
	Short: "copy 命令将本地的文件拷贝到远端",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		globalBuilder.NewCopyBuilder().Sudo(copySudo).Exists(copyExists).Src(args[:len(args)-1]).Dst(args[len(args)-1]).Builder().Run()
	},
}
