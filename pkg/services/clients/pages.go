package clients

import (
	"sort"

	clients_model "github.com/console-dns/server/pkg/models/clients"
	"github.com/console-dns/server/pkg/utils"
	"github.com/console-dns/server/pkg/utils/route"
)

func ClientRootPage(ctx *route.WebRequest) error {
	cfg, f := ctx.Content.SyncTokens.WithReadOnly()
	defer f()
	keys := cfg.ListByGroup(clients_model.TypeClient)
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})
	return ctx.WritePageTemplate("clients/pages/clients", map[string]any{
		"Title":   "客户端管理",
		"Clients": keys,
		"Regex":   utils.DataRegex,
		"Select":  ctx.FormValue("select"),
	})
}
