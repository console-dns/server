package route

import (
	"net/http"

	"github.com/console-dns/server/pkg/content"
	"github.com/console-dns/server/pkg/models/logs"
)

type ApiHandler func(ctx *ApiRequest) (any, error)

type ApiRequest struct {
	*http.Request
	Client  string
	Content *content.Content
	logs.LogPusher
	http.Handler
	Body string
}
