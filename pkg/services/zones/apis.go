package zones

import (
	self_errors "github.com/console-dns/server/pkg/errors"
	zones_model "github.com/console-dns/server/pkg/models/zones"
	"github.com/console-dns/server/pkg/utils"
	"github.com/console-dns/server/pkg/utils/route"
)

func ApiListZones(ctx *route.ApiRequest) (any, error) {
	zones, zUnlock := ctx.Content.SyncZones.WithReadOnly()
	clients, cUnlock := ctx.Content.SyncTokens.WithReadOnly()
	defer zUnlock()
	defer cUnlock()
	client := clients.Get(ctx.Client)
	result := zones_model.NewZones()

	for zoneName, zone := range zones.ListRecords() {
		for recordName, record := range zone {
			for _, dnsType := range zones_model.RecordTypes {
				if client.VerifyRule(zoneName, recordName, dnsType).CanRead() {
					r := record
					result.CopyFrom(&r, zoneName, recordName, dnsType)
				}
			}
		}
	}
	result.Clean(true)
	return result, nil
}

func ApiListZone(ctx *route.ApiRequest) (any, error) {
	zone := ctx.PathValue("zone")
	if err := utils.RegexHost.Valid(zone); err != nil {
		return nil, self_errors.BadRequestErrorf("区域名称不合法")
	}
	zones, zUnlock := ctx.Content.SyncZones.WithReadOnly()
	clients, cUnlock := ctx.Content.SyncTokens.WithReadOnly()
	defer zUnlock()
	defer cUnlock()
	token := clients.Get(ctx.Client)
	result := zones_model.NewZones()
	if zones.GetZone(zone) == nil {
		return nil, self_errors.NotfoundErrorf("区域不存在")
	}
	for recordName, record := range zones.GetRecords(zone) {
		for _, dnsType := range zones_model.RecordTypes {
			if token.VerifyRule(zone, recordName, dnsType).CanRead() {
				result.CopyFrom(record, zone, recordName, dnsType)
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
		r = zones_model.NewRecord()
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
		r = zones_model.NewRecord()
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
		r = zones_model.NewRecord()
		z.Records[data.record] = r
	}
	err = removeRecord(r, data.dnsType, func(s string) string {
		return data.items[s]
	})
	if err != nil {
		return nil, self_errors.BadRequestErrorf("内容删除错误: %s", err.Error())
	}
	ctx.PushLog("删除区域 $2 记录 $5 类型 $6", "zones", data.zone, "record", "remove", data.record, data.dnsType)
	return "删除完成", nil
}
