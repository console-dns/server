package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/console-dns/server/cmd"
)

func main() {
	env, b := os.LookupEnv("DEBUG")
	if b && env == "true" {
		log.Printf("DEBUG 模式已启动")
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}
	_ = cmd.Execute()
}
