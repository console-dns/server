//go:build !dev

package assets

import (
	"embed"
	"log/slog"
	"time"
)

//go:embed static template
var emb embed.FS

func init() {
	slog.Debug("使用内嵌模板")
	Cache = true
	StaticData = embedFs{
		emb,
	}
}

type embedFs struct {
	embed.FS
}

func (e embedFs) ModTime(_ string) (time.Time, error) {
	return BuildTime, nil
}
