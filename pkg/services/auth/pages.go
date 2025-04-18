package auth

import (
	"net"
	"sort"
	"strings"
	"time"

	"github.com/console-dns/server/pkg/models/auth"
	"github.com/console-dns/server/pkg/models/clients"
	"github.com/console-dns/server/pkg/utils/route"
	"github.com/mileusna/useragent"
)

func LoginGet(ctx *route.WebRequest) error {
	return ctx.WritePageTemplate("auth/pages/login", map[string]any{
		"Title": "管理员登录",
		"Error": strings.ToLower(ctx.FormValue("error")),
	})
}

type PageSessionsGetResult struct {
	Self bool
	Type string
	Name string
	Last time.Time
	IP   net.IP
	UA   *useragent.UserAgent
}

func PageSessionsGet(ctx *route.WebRequest) error {
	result := make([]*PageSessionsGetResult, 0)
	ctx.Content.SyncSessions.ReadOnly(func(session *auth.Session) {
		for s, state := range session.Sessions {
			agent := useragent.Parse(state.UserAgent)
			result = append(result, &PageSessionsGetResult{
				Self: ctx.LoginMeta.Session == s,
				Name: s[0:4],
				Type: "web",
				Last: state.UpdateTime,
				IP:   state.IpAddr,
				UA:   &agent,
			})
		}
	})
	ctx.Content.SyncTokens.ReadOnly(func(clients *clients.Clients) {
		for _, s := range clients.GetClientStatus() {
			userAgent := useragent.Parse(s.UA)
			result = append(result, &PageSessionsGetResult{
				Type: "api",
				Name: s.Name,
				Last: s.At,
				IP:   s.IP,
				UA:   &userAgent,
			})
		}
	})
	sort.Slice(result, func(i, j int) bool {
		return result[i].Last.After(result[j].Last)
	})
	return ctx.WritePageTemplate("auth/pages/sessions", map[string]any{
		"Title": "会话概览",
		"Data":  result,
	})
}
