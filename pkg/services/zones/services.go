package zones

import (
	"strconv"

	zones_model "github.com/console-dns/spec/models"
	spec_utils "github.com/console-dns/spec/utils"
	"github.com/pkg/errors"

	self_errors "github.com/console-dns/server/pkg/errors"
	"github.com/console-dns/server/pkg/utils/route"
)

func ZoneAdd(ctx *route.WebRequest) error {
	cfg, f := ctx.Content.SyncZones.WithReadWrite()
	defer f()
	zone := ctx.FormValue("zone")
	if err := cfg.AddZone(zone); err != nil {
		return self_errors.BadRequestErrorf("添加区域失败. %s", err.Error())
	}
	ctx.PushLog("新增区域 $2", "zones", zone, "create")
	return self_errors.ErrorRedirect("/zones")
}

func ZoneDelete(ctx *route.WebRequest) error {
	cfg, f := ctx.Content.SyncZones.WithReadWrite()
	defer f()
	zone := ctx.FormValue("zone")
	if err := cfg.RemoveZone(zone); err != nil {
		return self_errors.BadRequestErrorf("删除区域失败. %s", err.Error())
	}
	ctx.PushLog("删除区域 $2", "zones", zone, "remove")
	return self_errors.ErrorRedirect("/zones")
}

func ZoneRecordMod(ctx *route.WebRequest) error {
	cfg, f := ctx.Content.SyncZones.WithReadWrite()
	defer f()
	zone := ctx.PathValue("zone")
	dnsName := ctx.FormValue("name")
	action := ctx.FormValue("action")
	isUpdate := action == "update"
	index := -1
	if isUpdate {
		indexStr := ctx.FormValue("index")
		if indexStr != "" {
			var err error
			index, err = strconv.Atoi(indexStr)
			if err != nil {
				return self_errors.BadRequestErrorf("操作失败. %s", err.Error())
			}
		}
	}
	zoneCfg := cfg.GetZone(zone)
	if zoneCfg == nil {
		return self_errors.BadRequestErrorf("操作失败，区域不存在")
	}

	if err := spec_utils.RegexDnsName.Valid(dnsName); err != nil {
		return self_errors.BadRequestErrorf("操作失败,%s", err.Error())
	}
	recordType := ctx.FormValue("type")
	if err := zoneCfg.ModRecord(dnsName, func(r *zones_model.Record) error {
		return modRecordWithIndex(r, recordType, index, func(s string) string {
			return ctx.FormValue(s)
		})
	}); err != nil {
		return self_errors.BadRequestErrorf("操作失败,%s", err.Error())
	}
	if isUpdate {
		ctx.PushLog("修改区域 $2 记录 $5 类型 $6", "zones", zone, "record", "mod", dnsName, recordType)
	} else {
		ctx.PushLog("新增区域 $2 记录 $5 类型 $6", "zones", zone, "record", "add", dnsName, recordType)
	}
	return self_errors.ErrorRedirect("/zones?select=" + zone)
}

func ZoneRecordDel(ctx *route.WebRequest) error {
	zone := ctx.PathValue("zone")
	recordType := ctx.FormValue("type")
	dnsName := ctx.FormValue("name")
	indexStr := ctx.FormValue("index")
	index, err := strconv.Atoi(indexStr)
	cfg, f := ctx.Content.SyncZones.WithReadWrite()
	defer f()
	if err != nil {
		return self_errors.BadRequestErrorf("删除失败， %s", err.Error())
	}
	zoneCfg := cfg.GetZone(zone)
	if zoneCfg == nil {
		return self_errors.BadRequestErrorf("删除失败，区域不存在")
	}
	if err = zoneCfg.ModRecord(dnsName, func(r *zones_model.Record) error {
		switch recordType {
		case "A":
			return r.RemoveA(r.A[index])
		case "AAAA":
			return r.RemoveAAAA(r.AAAA[index])
		case "TXT":
			return r.RemoveTXT(r.TXT[index])
		case "CNAME":
			return r.RemoveCNAME(r.CNAME[index])
		case "NS":
			return r.RemoveNS(r.NS[index])
		case "MX":
			return r.RemoveMX(r.MX[index])
		case "SRV":
			return r.RemoveSRV(r.SRV[index])
		case "CAA":
			return r.RemoveCAA(r.CAA[index])
		case "SOA":
			return r.RemoveSOA(r.SOA)
		default:
			return errors.New("未知类型")
		}
	}); err != nil {
		return self_errors.BadRequestErrorf("删除失败， %s", err.Error())
	}
	ctx.PushLog("删除区域 $2 记录 $5 类型 $6", "zones", zone, "record", "remove", dnsName, recordType)
	return self_errors.ErrorRedirect("/zones?select=" + zone)
}
