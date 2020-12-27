package cmd

import (
	"github.com/spf13/cobra"
)

var (
	scriptSudo     bool
	scriptArgs     string
	scriptSaveFile string
)

func init() {
	rootCmd.AddCommand(&scriptCmd)
	scriptCmd.Flags().BoolVarP(&scriptSudo, "sudo", "S", false, "是否以sudo方式执行脚本")
	scriptCmd.Flags().StringVar(&scriptArgs, "args", "", "添加脚本执行的参数")
	scriptCmd.Flags().StringVar(&scriptSaveFile, "save", "", "将脚本输出保存到文件中")
}

var scriptCmd = cobra.Command{
	Use:   "script file",
	Short: "将本地脚本上传到远端并执行",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		globalBuilder.NewScriptBuilder().Sudo(scriptSudo).Path(args[0]).Args(scriptArgs).Builder().Run()
	},
}
