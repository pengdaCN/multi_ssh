package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(&pingCmd)
}

var pingCmd = cobra.Command{
	Use:   "ping",
	Short: "用于测试主机是否可以连通",
	Args:  cobra.MaximumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		globalBuilder.NewPingRunEnv().Run()
	},
}
