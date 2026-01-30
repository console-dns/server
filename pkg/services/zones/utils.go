package zones

import (
	"encoding/json"
	"fmt"
	"text/scanner"

	zones_model "github.com/console-dns/server/pkg/models/zones"
	"github.com/console-dns/server/pkg/utils"

	self_errors "github.com/console-dns/server/pkg/errors"
	"github.com/console-dns/server/pkg/utils/route"
)

type ZoneStatus struct {
	Name string
	*zones_model.Zone
	_close func()
}

func (receiver *ZoneStatus) Close() {
	receiver._close()
}

func tryZone(ctx *route.WebRequest, rw bool) (*ZoneStatus, error) {
	zone := ctx.PathValue("zone")
	err := utils.RegexHost.Valid(zone)
	if err != nil {
		return nil, self_errors.BadRequestErrorf("区域不存在")
	}
	if rw {
		cfg, f := ctx.Content.SyncZones.WithReadWrite()
		if cfg.GetZone(zone) == nil {
			defer f()
			return nil, self_errors.BadRequestErrorf("区域不存在")
		} else {
			return &ZoneStatus{
				Name:   zone,
				Zone:   cfg.GetZone(zone),
				_close: f,
			}, nil
		}
	} else {
		cfg, f := ctx.Content.SyncZones.WithReadOnly()
		if cfg.GetZone(zone) == nil {
			defer f()
			return nil, self_errors.BadRequestErrorf("区域不存在")
		} else {
			return &ZoneStatus{
				Name:   zone,
				Zone:   cfg.GetZone(zone),
				_close: f,
			}, nil
		}
	}
}

func getOrNil[V any](data []*V, index int) *V {
	if index == scanner.EOF {
		return nil
	}
	if index >= len(data) {
		return nil
	}
	return data[index]
}

func modRecord(r *zones_model.Record, dnsType string, oldData func(string) string, newData func(string) string) error {
	switch dnsType {
	case "A":
		oRecord, err := zones_model.FromRecordA(oldData)
		if err != nil {
			return err
		}
		record, err := zones_model.FromRecordA(newData)
		if err != nil {
			return err
		}
		return r.ModA(oRecord, record)
	case "AAAA":
		oRecord, err := zones_model.FromRecordAAAA(oldData)
		if err != nil {
			return err
		}
		record, err := zones_model.FromRecordAAAA(newData)
		if err != nil {
			return err
		}
		return r.ModAAAA(oRecord, record)
	case "TXT":
		oldRecord, err := zones_model.FromRecordTXT(oldData)
		if err != nil {
			return err
		}
		record, err := zones_model.FromRecordTXT(newData)
		if err != nil {
			return err
		}
		return r.ModTXT(oldRecord, record)
	case "CNAME":
		oldRecord, err := zones_model.FromRecordCNAME(oldData)
		if err != nil {
			return err
		}
		record, err := zones_model.FromRecordCNAME(newData)
		if err != nil {
			return err
		}
		return r.ModCNAME(oldRecord, record)
	case "NS":
		oldRecord, err := zones_model.FromRecordNS(oldData)
		if err != nil {
			return err
		}
		record, err := zones_model.FromRecordNS(newData)
		if err != nil {
			return err
		}
		return r.ModNS(oldRecord, record)
	case "MX":
		oldRecord, err := zones_model.FromRecordMX(oldData)
		if err != nil {
			return err
		}
		record, err := zones_model.FromRecordMX(newData)
		if err != nil {
			return err
		}
		return r.ModMX(oldRecord, record)
	case "SRV":
		oldRecord, err := zones_model.FromRecordSRV(oldData)
		if err != nil {
			return err
		}
		record, err := zones_model.FromRecordSRV(newData)
		if err != nil {
			return err
		}
		return r.ModSRV(oldRecord, record)
	case "SOA":
		oldSOA, err := zones_model.FromRecordSOA(oldData)
		if err != nil {
			return err
		}
		data, err := zones_model.FromRecordSOA(newData)
		if err != nil {
			return err
		}
		return r.ModSOA(oldSOA, data)
	case "CAA":
		oldRecord, err := zones_model.FromRecordCAA(oldData)
		if err != nil {
			return err
		}
		data, err := zones_model.FromRecordCAA(newData)
		if err != nil {
			return err
		}
		return r.ModCAA(oldRecord, data)
	default:
		return fmt.Errorf("未知类型")
	}
}

func modRecordWithIndex(r *zones_model.Record, dnsType string, oldIndex int, newData func(string) string) error {
	switch dnsType {
	case "A":
		recordA, err := zones_model.FromRecordA(newData)
		if err != nil {
			return err
		}
		return r.ModA(getOrNil(r.A, oldIndex), recordA)
	case "AAAA":
		aaaa, err := zones_model.FromRecordAAAA(newData)
		if err != nil {
			return err
		}
		return r.ModAAAA(getOrNil(r.AAAA, oldIndex), aaaa)
	case "TXT":
		txt, err := zones_model.FromRecordTXT(newData)
		if err != nil {
			return err
		}
		return r.ModTXT(getOrNil(r.TXT, oldIndex), txt)
	case "CNAME":
		data, err := zones_model.FromRecordCNAME(newData)
		if err != nil {
			return err
		}
		return r.ModCNAME(getOrNil(r.CNAME, oldIndex), data)
	case "NS":
		data, err := zones_model.FromRecordNS(newData)
		if err != nil {
			return err
		}
		return r.ModNS(getOrNil(r.NS, oldIndex), data)
	case "MX":
		data, err := zones_model.FromRecordMX(newData)
		if err != nil {
			return err
		}
		return r.ModMX(getOrNil(r.MX, oldIndex), data)
	case "SRV":
		data, err := zones_model.FromRecordSRV(newData)
		if err != nil {
			return err
		}
		return r.ModSRV(getOrNil(r.SRV, oldIndex), data)
	case "SOA":
		data, err := zones_model.FromRecordSOA(newData)
		if err != nil {
			return err
		}
		if oldIndex == scanner.EOF {
			return r.ModSOA(nil, data)
		} else {
			return r.ModSOA(r.SOA, data)
		}
	case "CAA":
		data, err := zones_model.FromRecordCAA(newData)
		if err != nil {
			return err
		}
		return r.ModCAA(getOrNil(r.CAA, oldIndex), data)
	default:
		return fmt.Errorf("未知类型")
	}
}

type mergeDataMod struct {
	Src map[string]any `json:"src"`
	Dst map[string]any `json:"dst"`
}

type mergeData struct {
	zone    string
	record  string
	dnsType string
	src     map[string]string
	dst     map[string]string
}

func parseMergeData(ctx *route.ApiRequest) (*mergeData, error) {
	zone, record, dnsType, err := checkParams(ctx)
	if err != nil {
		return nil, err
	}
	mod := mergeDataMod{
		Src: make(map[string]any),
		Dst: make(map[string]any),
	}
	err = json.Unmarshal([]byte(ctx.Body), &mod)
	if err != nil {
		return nil, self_errors.BadRequestErrorf("请求解析失败: %s", err.Error())
	}

	src := make(map[string]string)
	dst := make(map[string]string)
	for s, a := range mod.Src {
		src[s] = fmt.Sprintf("%v", a)
	}
	for s, a := range mod.Dst {
		dst[s] = fmt.Sprintf("%v", a)
	}
	return &mergeData{
			zone:    zone,
			record:  record,
			dnsType: dnsType,
			src:     src,
			dst:     dst,
		},
		nil
}

type commitData struct {
	zone    string
	record  string
	dnsType string
	items   map[string]string
}

func parseCommitData(ctx *route.ApiRequest) (*commitData, error) {
	zone, record, dnsType, err := checkParams(ctx)
	if err != nil {
		return nil, err
	}
	items := map[string]any{}
	err = json.Unmarshal([]byte(ctx.Body), &items)
	if err != nil {
		return nil, self_errors.BadRequestErrorf("请求解析失败: %s", err.Error())
	}
	values := make(map[string]string)
	for s, a := range items {
		values[s] = fmt.Sprintf("%v", a)
	}
	return &commitData{
		zone:    zone,
		record:  record,
		dnsType: dnsType,
		items:   values,
	}, nil
}

func checkParams(ctx *route.ApiRequest) (string, string, string, error) {
	zone := ctx.PathValue("zone")
	if err := utils.RegexHost.Valid(zone); err != nil {
		return "", "", "", self_errors.BadRequestErrorf("区域名称不合法")
	}
	record := ctx.PathValue("record")
	if err := utils.RegexDnsName.Valid(record); err != nil {
		return "", "", "", self_errors.BadRequestErrorf("记录名称不合法")
	}
	dnsType := ctx.PathValue("type")
	if err := utils.RegexDnsType.Valid(dnsType); err != nil {
		return "", "", "", self_errors.BadRequestErrorf("DNS 类型不合法")
	}
	return zone, record, dnsType, nil
}
