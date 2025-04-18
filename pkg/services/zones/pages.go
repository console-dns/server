package zones

import (
	"sort"

	"github.com/console-dns/server/pkg/utils/route"
)

func ZoneRootPage(ctx *route.WebRequest) error {
	cfg, f := ctx.Content.SyncZones.WithReadOnly()
	defer f()
	zones := cfg.ListZones()
	sort.Strings(zones)
	return ctx.WritePageTemplate("zones/pages/zones", map[string]any{
		"Title":  "区域管理",
		"Zones":  zones,
		"Select": ctx.FormValue("select"),
	})
}
