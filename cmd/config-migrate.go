package cmd

import (
	"fmt"
	"os"

	"github.com/console-dns/server/pkg/content/settings"
	"github.com/console-dns/server/pkg/utils"
	"github.com/spf13/cobra"
)

var configMigrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "迁移配置文件",
	Run: func(cmd *cobra.Command, args []string) {
		if _, err := os.Stat(configPath); err != nil {
			println("配置文件不存在")
			os.Exit(1)
		}
		staticConfig := settings.NewStaticConfig()
		err := utils.AutoUnmarshal(configPath, staticConfig, false)
		if err != nil {
			fmt.Printf("配置文件加载失败： %v\n", err.Error())
			os.Exit(1)
		}
		err = utils.AutoMarshal(configPath, staticConfig)
		if err != nil {
			fmt.Printf("配置文件迁移失败： %v\n", err.Error())
			os.Exit(1)
		}
		fmt.Printf("配置文件迁移完成！\n")
	},
}

func init() {
	configCmd.AddCommand(configMigrateCmd)
}
