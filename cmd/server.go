package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/console-dns/server/pkg"
	"github.com/console-dns/server/pkg/content"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "启动 Console DNS HTTP 控制台",
	Run: func(cmd *cobra.Command, args []string) {
		content, err := content.NewContent(configPath)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		ctx, exitFunc := context.WithCancel(context.Background())
		childCtx, signExit := context.WithCancel(ctx)
		route := pkg.NewConsoleRoute(content, childCtx)
		route.StartAsync(exitFunc)
		log.Printf("Web 服务启动，请访问 http://%s", content.Config.Server.AddressPort())
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		<-sig
		signExit()   // 主动停止
		<-ctx.Done() // 等待退出
		_ = content.Close()
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
	cobra.OnInitialize()
}
