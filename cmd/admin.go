package cmd

import "github.com/spf13/cobra"

var adminCmd = &cobra.Command{
	Use:   "admin",
	Short: "管理员操作",
}

func init() {
	rootCmd.AddCommand(adminCmd)
}
