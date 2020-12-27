package cmd

import (
	"github.com/spf13/cobra"
)

var (
	argsList string
)

func init() {
	rootCmd.AddCommand(&playbookCmd)
	playbookCmd.Flags().StringVarP(&argsList, "set-args", "S", "", "设置lua中全局的变量，key=val形式")
}

var playbookCmd = cobra.Command{
	Use:     "playbook <file>",
	Short:   "执行一系列的命令",
	Long:    "通过golang内置的lua虚拟机来执行一系列操作",
	Args:    cobra.MinimumNArgs(1),
	Example: "playbook example.play",
	Run: func(cmd *cobra.Command, args []string) {
		globalBuilder.NewPlaybookBuilder().SetVars(argsList).Path(args[0]).Builder().Run()
	},
}
