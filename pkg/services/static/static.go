package static

import (
	"strings"
	"time"

	"github.com/console-dns/server/assets"
	self_errors "github.com/console-dns/server/pkg/errors"
	"github.com/console-dns/server/pkg/utils/route"
)

func AssetsGet(ctx *route.WebRequest) error {
	targetPath := ""
	if ctx.RequestURI == "/favicon.ico" {
		targetPath = "static/favicon.png"
	} else if strings.HasPrefix(ctx.RequestURI, "/static/") {
		targetPath = strings.TrimPrefix(ctx.RequestURI, "/")
	} else {
		return self_errors.NotfoundErrorf("static assets not found")
	}

	file, err := assets.OpenResFile(targetPath)
	if err != nil {
		return self_errors.NotfoundErrorf("static assets not found")
	}
	ctx.W.Header().Add("Content-Type", file.Mine)
	if assets.Cache {
		ctx.W.Header().Add("Age", "86400")
		ctx.W.Header().Add("Cache-Control", "max-age=604800")
		ctx.W.Header().Add("Last-Modified", file.ModTime.Format(time.RFC822))
	} else {
		ctx.W.Header().Add("Cache-Control", "no-store, no-cache, max-age=0, must-revalidate, proxy-revalidate")
	}
	_, _ = ctx.W.Write(file.Data)
	return nil
}
