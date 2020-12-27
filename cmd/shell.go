package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	shellCmd.Flags().BoolVar(&shellSudo, "sudo", false, "设置是否以sudo方式执行命令")
	shellCmd.Flags().StringVar(&shellSaveFile, "save", "", "将shell命令执行的输出保存到本地的文件中")
	rootCmd.AddCommand(&shellCmd)
}

var (
	shellSudo     bool
	shellSaveFile string
)

var shellCmd = cobra.Command{
	Use:     "shell command",
	Short:   "执行一行shell命令",
	Example: "shell --sudo 'command'",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		globalBuilder.NewShellBuilder().Sudo(shellSudo).Cmds(args[0]).Builder().Run()
	},
}
