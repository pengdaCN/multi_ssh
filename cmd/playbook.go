package cmd

import "github.com/spf13/cobra"

func init() {
	rootCmd.AddCommand(&playbookCmd)
}

var playbookCmd = cobra.Command{
	Use:   "playbook [flag] [arg] <file>",
	Short: "执行一系列的命令",
	Example: "playbook example.play",
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {

	},
}
