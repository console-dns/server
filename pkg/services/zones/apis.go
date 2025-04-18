package zones

import (
	self_errors "github.com/console-dns/server/pkg/errors"
	"github.com/console-dns/server/pkg/utils/route"
	"github.com/console-dns/spec/models"
	spec_utils "github.com/console-dns/spec/utils"
	"github.com/pkg/errors"
)

var recordTypes = []string{
	"A",
	"AAAA",
	"TXT",
	"CNAME",
	"NS",
	"MX",
	"SRV",
	"CAA",
	"SOA",
}

func ApiListZones(ctx *route.ApiRequest) (any, error) {
	zones, zUnlock := ctx.Content.SyncZones.WithReadOnly()
	clients, cUnlock := ctx.Content.SyncTokens.WithReadOnly()
	defer zUnlock()
	defer cUnlock()
	client := clients.Get(ctx.Client)
	result := models.NewZones()

	for zoneName, zone := range zones.ListRecords() {
		for recordName, record := range zone {
			for _, dnsType := range recordTypes {
				if client.VerifyRule(zoneName, recordName, dnsType).CanRead() {
					result.CopyFrom(&record, zoneName, recordName, dnsType)
				}
			}
		}
	}
	result.Clean(true)
	return result, nil
}

func ApiListZone(ctx *route.ApiRequest) (any, error) {
	zone := ctx.PathValue("zone")
	if err := spec_utils.RegexHost.Valid(zone); err != nil {
		return nil, self_errors.BadRequestErrorf("区域名称不合法")
	}
	zones, zUnlock := ctx.Content.SyncZones.WithReadOnly()
	clients, cUnlock := ctx.Content.SyncTokens.WithReadOnly()
	defer zUnlock()
	defer cUnlock()
	token := clients.Get(ctx.Client)
	result := models.NewZones()
	if zones.GetZone(zone) == nil {
		return nil, self_errors.NotfoundErrorf("区域不存在")
	}
	for recordName, record := range zones.GetRecords(zone) {
		for _, dnsType := range recordTypes {
			if token.VerifyRule(zone, recordName, dnsType).CanRead() {
				result.CopyFrom(&record, zone, recordName, dnsType)
			}
		}
	}
	result.Clean(false)
	return result.GetZone(zone), nil
}

func CreateRecord(ctx *route.ApiRequest) (any, error) {
	data, err := parseCommitData(ctx)
	if err != nil {
		return nil, err
	}
	zones, zUnlock := ctx.Content.SyncZones.WithReadWrite()
	clients, cUnlock := ctx.Content.SyncTokens.WithReadWrite()
	defer zUnlock()
	defer cUnlock()
	token := clients.Get(ctx.Client)
	if !token.VerifyRule(data.zone, data.record, data.dnsType).CanWrite() {
		return nil, self_errors.ForbiddenErrorf("用户没有操作此 dns 记录的权限")
	}
	z := zones.GetZone(data.zone)
	if z == nil {
		return nil, self_errors.BadRequestErrorf("区域不存在")
	}
	r := z.Records[data.record]
	if r == nil {
		r = models.NewRecord()
		z.Records[data.record] = r
	}
	err = modRecordWithIndex(r, data.dnsType, -1, func(s string) string {
		return data.items[s]
	})
	if err != nil {
		return nil, self_errors.BadRequestErrorf("记录插入错误: %s", err.Error())
	}
	ctx.PushLog("新增区域 $2 记录 $5 类型 $6", "zones", data.zone, "record", "add", data.record, data.dnsType)
	return nil, nil
}

func ModRecord(ctx *route.ApiRequest) (any, error) {
	data, err := parseMergeData(ctx)
	if err != nil {
		return nil, err
	}
	zones, zUnlock := ctx.Content.SyncZones.WithReadWrite()
	clients, cUnlock := ctx.Content.SyncTokens.WithReadWrite()
	defer zUnlock()
	defer cUnlock()
	token := clients.Get(ctx.Client)
	if !token.VerifyRule(data.zone, data.record, data.dnsType).CanWrite() {
		return nil, self_errors.ForbiddenErrorf("用户没有操作此 dns 记录的权限")
	}
	z := zones.GetZone(data.zone)
	if z == nil {
		return nil, self_errors.BadRequestErrorf("区域不存在")
	}
	r := z.Records[data.record]
	if r == nil {
		r = models.NewRecord()
		z.Records[data.record] = r
	}
	err = modRecord(r, data.dnsType, func(s string) string {
		return data.src[s]
	}, func(s string) string {
		return data.dst[s]
	})
	if err != nil {
		return nil, self_errors.BadRequestErrorf("记录修改错误: %s", err.Error())
	}
	ctx.PushLog("修改区域 $2 记录 $5 类型 $6", "zones", data.zone, "record", "mod", data.record, data.dnsType)
	return nil, nil
}

func DeleteRecord(ctx *route.ApiRequest) (any, error) {
	data, err := parseCommitData(ctx)
	if err != nil {
		return nil, err
	}
	zones, zUnlock := ctx.Content.SyncZones.WithReadWrite()
	clients, cUnlock := ctx.Content.SyncTokens.WithReadWrite()
	defer zUnlock()
	defer cUnlock()
	token := clients.Get(ctx.Client)
	if !token.VerifyRule(data.zone, data.record, data.dnsType).CanDelete() {
		return nil, self_errors.ForbiddenErrorf("用户没有操作此 dns 记录的权限")
	}
	z := zones.GetZone(data.zone)
	if z == nil {
		return nil, self_errors.BadRequestErrorf("区域不存在")
	}
	r := z.Records[data.record]
	if r == nil {
		return nil, self_errors.BadRequestErrorf("记录不存在")
	}
	getValue := func(s string) string {
		return data.items[s]
	}
	switch data.dnsType {
	case "A":
		if data, e := models.FromRecordA(getValue); e != nil {
			return nil, e
		} else {
			err = r.RemoveA(data)
		}
	case "AAAA":
		if data, e := models.FromRecordAAAA(getValue); e != nil {
			return nil, e
		} else {
			err = r.RemoveAAAA(data)
		}
	case "TXT":
		if data, e := models.FromRecordTXT(getValue); e != nil {
			return nil, e
		} else {
			err = r.RemoveTXT(data)
		}
	case "CNAME":
		if data, e := models.FromRecordCNAME(getValue); e != nil {
			return nil, e
		} else {
			err = r.RemoveCNAME(data)
		}
	case "NS":
		if data, e := models.FromRecordNS(getValue); e != nil {
			return nil, e
		} else {
			err = r.RemoveNS(data)
		}
	case "MX":
		if data, e := models.FromRecordMX(getValue); e != nil {
			return nil, e
		} else {
			err = r.RemoveMX(data)
		}
	case "SRV":
		if data, e := models.FromRecordSRV(getValue); e != nil {
			return nil, e
		} else {
			err = r.RemoveSRV(data)
		}
	case "CAA":
		if data, e := models.FromRecordCAA(getValue); e != nil {
			return nil, e
		} else {
			err = r.RemoveCAA(data)
		}
	case "SOA":
		if data, e := models.FromRecordSOA(getValue); e != nil {
			return nil, e
		} else {
			err = r.RemoveSOA(data)
		}
	default:
		err = errors.New("未知类型")
	}
	if err != nil {
		return nil, self_errors.BadRequestErrorf("内容删除错误: %s", err.Error())
	}
	ctx.PushLog("删除区域 $2 记录 $5 类型 $6", "zones", data.zone, "record", "remove", data.record, data.dnsType)
	return "删除完成", nil
}
