package cmd

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	DefConfigPath = ""
	configPath    = ""
)

var rootCmd = &cobra.Command{
	Use:     filepath.Base(os.Args[0]),
	Short:   "一个简单的 dns 管理控制台",
	Version: "0.0.1",
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", DefConfigPath, "配置文件地址")
}

func Execute() error {
	return rootCmd.Execute()
}
