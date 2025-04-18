package logs

import (
	"time"

	"github.com/console-dns/server/pkg/utils/route"
)

func PageLogsGet(ctx *route.WebRequest) error {
	return ctx.WritePageTemplate("logs/pages/logs", map[string]any{
		"Title": "审计日志",
		"Begin": time.Now().AddDate(0, -1, 0).Format("2006-01-02T15:04"),
		"End":   time.Now().Add(1 * time.Minute).Format("2006-01-02T15:04"),
	})
}
