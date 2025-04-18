package cmd

import (
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Console DNS 配置管理",
}

func init() {
	rootCmd.AddCommand(configCmd)
}
