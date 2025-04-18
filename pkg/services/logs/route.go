package logs

import (
	"github.com/console-dns/server/pkg/utils/route"
	"github.com/go-chi/chi/v5"
)

func Route(r chi.Router) {
	r.Get("/", route.MVCBind(PageLogsGet))
	r.Get("/!/query", route.MVCBind(QueryByGroup))
}
