package auth

import (
	"github.com/console-dns/server/pkg/utils/route"
	"github.com/go-chi/chi/v5"
)

func LoginRoute(r chi.Router) {
	r.Get("/", route.MVCBind(LoginGet))
	r.Post("/", route.MVCBind(LoginPost))
}

func SessionRoute(r chi.Router) {
	r.Get("/", route.MVCBind(PageSessionsGet))
}
