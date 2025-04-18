package zones

import (
	"github.com/console-dns/server/pkg/content/settings"
	"github.com/console-dns/server/pkg/utils"
	"github.com/console-dns/spec/models"
)

type Zones struct {
	localPath     string
	*models.Zones // 区域列表
}

func FromZones(cfg *settings.StaticConfig) (*Zones, error) {
	data := models.NewZones()
	err := utils.AutoUnmarshal(cfg.Storage.Zone, data, true)
	if err != nil {
		return nil, err
	}
	return &Zones{
		localPath: cfg.Storage.Zone,
		Zones:     data,
	}, nil
}

func (z *Zones) Flush() error {
	return utils.AutoMarshal(z.localPath, z.Zones)
}
