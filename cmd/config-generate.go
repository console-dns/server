package cmd

import (
	"github.com/console-dns/server/pkg/content/settings"
	"github.com/console-dns/server/pkg/utils"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var configGenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "生成默认配置文件",
	Run: func(cmd *cobra.Command, args []string) {
		staticConfig := settings.NewStaticConfig()
		if len(configPath) != 0 {
			err := utils.AutoMarshal(configPath, staticConfig)
			if err != nil {
				panic(err)
			}
		} else {
			out, _ := yaml.Marshal(staticConfig)
			println(string(out))
		}
	},
}

func init() {
	configCmd.AddCommand(configGenerateCmd)
}
