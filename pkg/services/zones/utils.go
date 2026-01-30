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
	var cfg *zones_model.Zones
	var f func()
	if rw {
		cfg, f = ctx.Content.SyncZones.WithReadWrite()
	} else {
		cfg, f = ctx.Content.SyncZones.WithReadOnly()
	}
	z := cfg.GetZone(zone)
	if z == nil {
		f()
		return nil, self_errors.BadRequestErrorf("区域不存在")
	}
	return &ZoneStatus{
		Name:   zone,
		Zone:   z,
		_close: f,
	}, nil
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
		old, err1 := zones_model.FromRecordA(oldData)
		new, err2 := zones_model.FromRecordA(newData)
		if err1 != nil || err2 != nil {
			return utils.FirstError(err1, err2)
		}
		return r.ModA(old, new)
	case "AAAA":
		old, err1 := zones_model.FromRecordAAAA(oldData)
		new, err2 := zones_model.FromRecordAAAA(newData)
		if err1 != nil || err2 != nil {
			return utils.FirstError(err1, err2)
		}
		return r.ModAAAA(old, new)
	case "TXT":
		old, err1 := zones_model.FromRecordTXT(oldData)
		new, err2 := zones_model.FromRecordTXT(newData)
		if err1 != nil || err2 != nil {
			return utils.FirstError(err1, err2)
		}
		return r.ModTXT(old, new)
	case "CNAME":
		old, err1 := zones_model.FromRecordCNAME(oldData)
		new, err2 := zones_model.FromRecordCNAME(newData)
		if err1 != nil || err2 != nil {
			return utils.FirstError(err1, err2)
		}
		return r.ModCNAME(old, new)
	case "NS":
		old, err1 := zones_model.FromRecordNS(oldData)
		new, err2 := zones_model.FromRecordNS(newData)
		if err1 != nil || err2 != nil {
			return utils.FirstError(err1, err2)
		}
		return r.ModNS(old, new)
	case "MX":
		old, err1 := zones_model.FromRecordMX(oldData)
		new, err2 := zones_model.FromRecordMX(newData)
		if err1 != nil || err2 != nil {
			return utils.FirstError(err1, err2)
		}
		return r.ModMX(old, new)
	case "SRV":
		old, err1 := zones_model.FromRecordSRV(oldData)
		new, err2 := zones_model.FromRecordSRV(newData)
		if err1 != nil || err2 != nil {
			return utils.FirstError(err1, err2)
		}
		return r.ModSRV(old, new)
	case "SOA":
		old, err1 := zones_model.FromRecordSOA(oldData)
		new, err2 := zones_model.FromRecordSOA(newData)
		if err1 != nil || err2 != nil {
			return utils.FirstError(err1, err2)
		}
		return r.ModSOA(old, new)
	case "CAA":
		old, err1 := zones_model.FromRecordCAA(oldData)
		new, err2 := zones_model.FromRecordCAA(newData)
		if err1 != nil || err2 != nil {
			return utils.FirstError(err1, err2)
		}
		return r.ModCAA(old, new)
	default:
		return fmt.Errorf("未知类型")
	}
}

func modRecordWithIndex(r *zones_model.Record, dnsType string, oldIndex int, newData func(string) string) error {
	switch dnsType {
	case "A":
		record, err := zones_model.FromRecordA(newData)
		if err != nil {
			return err
		}
		return r.ModA(getOrNil(r.A, oldIndex), record)
	case "AAAA":
		record, err := zones_model.FromRecordAAAA(newData)
		if err != nil {
			return err
		}
		return r.ModAAAA(getOrNil(r.AAAA, oldIndex), record)
	case "TXT":
		record, err := zones_model.FromRecordTXT(newData)
		if err != nil {
			return err
		}
		return r.ModTXT(getOrNil(r.TXT, oldIndex), record)
	case "CNAME":
		record, err := zones_model.FromRecordCNAME(newData)
		if err != nil {
			return err
		}
		return r.ModCNAME(getOrNil(r.CNAME, oldIndex), record)
	case "NS":
		record, err := zones_model.FromRecordNS(newData)
		if err != nil {
			return err
		}
		return r.ModNS(getOrNil(r.NS, oldIndex), record)
	case "MX":
		record, err := zones_model.FromRecordMX(newData)
		if err != nil {
			return err
		}
		return r.ModMX(getOrNil(r.MX, oldIndex), record)
	case "SRV":
		record, err := zones_model.FromRecordSRV(newData)
		if err != nil {
			return err
		}
		return r.ModSRV(getOrNil(r.SRV, oldIndex), record)
	case "SOA":
		record, err := zones_model.FromRecordSOA(newData)
		if err != nil {
			return err
		}
		old := r.SOA
		if oldIndex == scanner.EOF {
			old = nil
		}
		return r.ModSOA(old, record)
	case "CAA":
		record, err := zones_model.FromRecordCAA(newData)
		if err != nil {
			return err
		}
		return r.ModCAA(getOrNil(r.CAA, oldIndex), record)
	default:
		return fmt.Errorf("未知类型")
	}
}

func removeRecord(r *zones_model.Record, dnsType string, data func(string) string) error {
	switch dnsType {
	case "A":
		it, err := zones_model.FromRecordA(data)
		if err != nil {
			return err
		}
		return r.RemoveA(it)
	case "AAAA":
		it, err := zones_model.FromRecordAAAA(data)
		if err != nil {
			return err
		}
		return r.RemoveAAAA(it)
	case "TXT":
		it, err := zones_model.FromRecordTXT(data)
		if err != nil {
			return err
		}
		return r.RemoveTXT(it)
	case "CNAME":
		it, err := zones_model.FromRecordCNAME(data)
		if err != nil {
			return err
		}
		return r.RemoveCNAME(it)
	case "NS":
		it, err := zones_model.FromRecordNS(data)
		if err != nil {
			return err
		}
		return r.RemoveNS(it)
	case "MX":
		it, err := zones_model.FromRecordMX(data)
		if err != nil {
			return err
		}
		return r.RemoveMX(it)
	case "SRV":
		it, err := zones_model.FromRecordSRV(data)
		if err != nil {
			return err
		}
		return r.RemoveSRV(it)
	case "SOA":
		it, err := zones_model.FromRecordSOA(data)
		if err != nil {
			return err
		}
		return r.RemoveSOA(it)
	case "CAA":
		it, err := zones_model.FromRecordCAA(data)
		if err != nil {
			return err
		}
		return r.RemoveCAA(it)
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
