package zones

import (
	"context"
	"regexp"

	"gopkg.d7z.net/middleware/kv"
)

type Zones struct {
	storage   kv.KV
	*ZoneData // 区域列表
}

func FromZones(ctx context.Context, storage kv.KV) (*Zones, error) {
	zStorage := storage.Child("zones")
	data := NewZones()
	list, err := zStorage.List(ctx, "")
	if err != nil {
		return nil, err
	}
	for key := range list {
		parts := regexp.MustCompile(`/`).Split(key, 2)
		if len(parts) > 0 {
			zoneName := parts[0]
			if _, ok := data.Data[zoneName]; !ok {
				zone := NewZone()
				err := zone.Load(ctx, zStorage.Child(zoneName))
				if err != nil {
					return nil, err
				}
				data.Data[zoneName] = zone
			}
		}
	}

	return &Zones{
		storage:  storage,
		ZoneData: data,
	}, nil
}

func (z *Zones) Flush(ctx context.Context) error {
	zStorage := z.storage.Child("zones")
	for name, zone := range z.Data {
		err := zone.Save(ctx, zStorage.Child(name))
		if err != nil {
			return err
		}
	}
	return nil
}

func (z *Zones) ListZones() []string {
	res := make([]string, 0, len(z.Data))
	for k := range z.Data {
		res = append(res, k)
	}
	return res
}
