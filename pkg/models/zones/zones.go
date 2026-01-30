package zones

import (
	"context"

	"github.com/console-dns/server/pkg/utils"
	"gopkg.d7z.net/middleware/kv"
)

type Zones struct {
	storage   kv.KV
	*ZoneData // 区域列表
}

func FromZones(ctx context.Context, storage kv.KV) (*Zones, error) {
	data := NewZones()
	err := utils.AutoKVUnmarshal(ctx, storage, "zones.json", data)
	if err != nil {
		return nil, err
	}
	return &Zones{
		storage:  storage,
		ZoneData: data,
	}, nil
}

func (z *Zones) Flush(ctx context.Context) error {
	return utils.AutoKVMarshal(ctx, z.storage, "zones.json", z.ZoneData)
}

func (z *Zones) ListZones() []string {
	res := make([]string, 0, len(z.Data))
	for k := range z.Data {
		res = append(res, k)
	}
	return res
}
