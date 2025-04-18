package zones

import (
	"sort"
	"strconv"

	"github.com/console-dns/server/pkg/errors"
	zones_model "github.com/console-dns/server/pkg/models/zones"
	"github.com/console-dns/server/pkg/utils"
	"github.com/console-dns/server/pkg/utils/route"
	"github.com/fatih/structs"
)

func UiZone(ctx *route.WebRequest) error {
	zone, err := tryZone(ctx, false)
	if err != nil {
		return err
	}
	defer zone.Close()
	list := make([]*zones_model.DNSView, 0)
	for s, record := range zone.Records {
		list = append(list, zones_model.ToList(record, s)...)
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].Name < list[j].Name
	})
	return ctx.WriteComponentTemplate("zones/components/zone", map[string]any{
		"Zone":  zone.Name,
		"Data":  list,
		"Regex": utils.DataRegex,
	})
}

func UiZonesAddDialog(ctx *route.WebRequest) error {
	return ctx.WriteComponentTemplate("zones/components/zones-add", map[string]any{
		"Regex": utils.DataRegex,
	})
}

func UiZonesDelete(ctx *route.WebRequest) error {
	zone, err := tryZone(ctx, false)
	if err != nil {
		return err
	}
	defer zone.Close()
	return ctx.WriteComponentTemplate("zones/components/zones-delete", map[string]any{
		"Zone": zone.Name,
	})
}

func UiZoneCreate(ctx *route.WebRequest) error {
	zone, err := tryZone(ctx, false)
	if err != nil {
		return err
	}
	defer zone.Close()
	return ctx.WriteComponentTemplate("zones/components/zone-mod", map[string]any{
		"Zone":   zone.Name,
		"Mode":   "create",
		"Type":   "A",
		"Record": make(map[string]any),
		"Regex":  utils.DataRegex,
	})
}

func UiZoneMod(ctx *route.WebRequest) error {
	recordName := ctx.PathValue("name")
	indexStr := ctx.PathValue("index")
	dnsType := ctx.PathValue("type")
	zone, err := tryZone(ctx, false)
	if err != nil {
		return err
	}
	defer zone.Close()
	index, err := strconv.Atoi(indexStr)
	if err != nil {
		return errors.BadRequestErrorf("未知 index 类型")
	}

	return ctx.WriteComponentTemplate("zones/components/zone-mod", map[string]any{
		"Zone":   zone.Name,
		"Mode":   "update",
		"Name":   recordName,
		"Type":   dnsType,
		"Index":  index,
		"Record": structs.Map(zone.Records[recordName]),
		"Regex":  utils.DataRegex,
	})
}

func UiZoneDelete(ctx *route.WebRequest) error {
	return ctx.WriteComponentTemplate("zones/components/zone-delete", map[string]any{
		"Zone":  ctx.PathValue("zone"),
		"Type":  ctx.PathValue("type"),
		"Name":  ctx.PathValue("name"),
		"Index": ctx.PathValue("index"),
		"Value": ctx.PathValue("value"),
	})
}
