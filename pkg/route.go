package pkg

import (
	"net/http"

	auth_service "github.com/console-dns/server/pkg/services/auth"
	client_service "github.com/console-dns/server/pkg/services/clients"
	logs_service "github.com/console-dns/server/pkg/services/logs"
	"github.com/console-dns/server/pkg/services/static"
	"github.com/console-dns/server/pkg/services/zones"
	"github.com/console-dns/server/pkg/utils/route"
	"github.com/go-chi/chi/v5"
)

func Index() http.Handler {
	r := chi.NewRouter()
	r.Use(auth_service.DenyIPsFilter)
	r.Use(auth_service.HackerIpsInject)
	r.Use(logs_service.DefaultLogsInject)
	r.Route("/login", auth_service.LoginRoute)
	r.Get("/static/*", route.MVCBind(static.AssetsGet))
	r.Get("/favicon.ico", route.MVCBind(static.AssetsGet))
	r.Route("/api/v1", func(r chi.Router) {
		r.Use(auth_service.ApiLoginFilter)
		r.Use(logs_service.ApiLogsInject)
		r.Route("/zones", zones.ApiRoute)
	})
	r.Route("/", func(r chi.Router) {
		r.Use(auth_service.ConsoleLoginFilter)
		r.Use(logs_service.ConsoleRequestLogsInject)
		r.Get("/", func(writer http.ResponseWriter, request *http.Request) {
			http.Redirect(writer, request, "/zones", 302)
		})
		r.Route("/zones", zones.Route)
		r.Route("/clients", client_service.Route)
		r.Route("/logs", logs_service.Route)
		r.Route("/sessions", auth_service.SessionRoute)
		r.Get("/logout", route.MVCBind(auth_service.LogoutGet))
	})
	return r
}
