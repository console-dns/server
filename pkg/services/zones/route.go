package zones

import (
	"github.com/console-dns/server/pkg/utils/route"
	"github.com/go-chi/chi/v5"
)

func Route(r chi.Router) {
	r.Get("/", route.MVCBind(ZoneRootPage))
	r.Post("/", route.MVCBind(ZoneAdd))
	r.Post("/{zone}", route.MVCBind(ZoneRecordMod))
	r.Post("/{zone}/delete", route.MVCBind(ZoneRecordDel))
	r.Get("/delete", route.MVCBind(ZoneDelete))

	r.Post("/!/zones", route.MVCBind(UiZonesAddDialog))
	r.Get("/!/zones/{zone}", route.MVCBind(UiZone))
	r.Delete("/!/zones/{zone}", route.MVCBind(UiZonesDelete))
	r.Post("/!/zone/{zone}", route.MVCBind(UiZoneCreate))
	r.Put("/!/zone/{zone}/{type}/{name}/{index}", route.MVCBind(UiZoneMod))
	r.Delete("/!/zone/{zone}/{type}/{name}/{index}/{value}", route.MVCBind(UiZoneDelete))
}

func ApiRoute(r chi.Router) {
	r.Get("/", route.ApiBind(ApiListZones))
	r.Get("/{zone}", route.ApiBind(ApiListZone))
	r.Post("/{zone}/{record}/{type}", route.ApiBind(CreateRecord))
	r.Post("/{zone}/{record}/{type}/mod", route.ApiBind(ModRecord))
	r.Post("/{zone}/{record}/{type}/delete", route.ApiBind(DeleteRecord))
}
