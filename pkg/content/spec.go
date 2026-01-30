package content

import (
	"context"
	"log/slog"

	"github.com/console-dns/server/pkg/content/settings"
	"github.com/console-dns/server/pkg/models/auth"
	clients_model "github.com/console-dns/server/pkg/models/clients"
	logs_model "github.com/console-dns/server/pkg/models/logs"
	zones_model "github.com/console-dns/server/pkg/models/zones"
	"github.com/console-dns/server/pkg/utils"
	"github.com/robfig/cron/v3"
	"gopkg.d7z.net/middleware/kv"
)

// Content 程序信息
type Content struct {
	Config       *settings.StaticConfig                      // 静态配置
	Storage      kv.CloserKV                                 // 存储
	SyncSessions *utils.DataRwLocker[*auth.Session]          // 用户会话管理
	SyncZones    *utils.DataRwLocker[*zones_model.Zones]     // DNS 区域
	SyncTokens   *utils.DataRwLocker[*clients_model.Clients] // 服务账户的 token
	Logs         logs_model.ConsoleLog                       // 审计日志
	cron         *cron.Cron
}

func NewContent(configPath string) (*Content, error) {
	staticCfg, err := settings.FromStaticConfig(configPath)
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	storage, err := kv.NewKVFromURL(staticCfg.Storage.DSN)
	if err != nil {
		return nil, err
	}

	sessions, err := auth.FromSession(ctx, staticCfg, storage)
	if err != nil {
		return nil, err
	}
	zones, err := zones_model.FromZones(ctx, storage)
	if err != nil {
		return nil, err
	}
	tokens, err := clients_model.FromClients(ctx, storage)
	if err != nil {
		return nil, err
	}
	kvLog := logs_model.NewKVLog(storage)

	wrapper := NewSystemLogs(kvLog)
	sessions.Logs = wrapper
	content := &Content{
		cron:         cron.New(),
		Config:       staticCfg,
		Storage:      storage,
		SyncSessions: utils.NewDataRwLocker(sessions),
		SyncZones:    utils.NewDataRwLocker(zones),
		SyncTokens:   utils.NewDataRwLocker(tokens),
		Logs:         kvLog,
	}
	_, _ = content.cron.AddFunc("@every 1m", func() {
		err := content.Refresh()
		if err != nil {
			slog.Debug("refresh content error", "err", err.Error())
		}
	})
	content.cron.Start()
	return content, nil
}

func (c *Content) Refresh() error {
	var err error
	c.SyncSessions.ReadWrite(func(session *auth.Session) {
		err = session.Refresh()
	})
	return err
}

func (c *Content) Close() error {
	ctx := context.Background()
	c.SyncZones.ReadWrite(func(zones *zones_model.Zones) {
		_ = zones.Flush(ctx)
	})
	c.SyncTokens.ReadWrite(func(tokens *clients_model.Clients) {
		_ = tokens.Flush(ctx)
	})
	c.SyncSessions.ReadWrite(func(zs *auth.Session) {
		_ = zs.Flush(ctx)
	})
	c.cron.Stop()
	return c.Storage.Close()
}

type SystemLogs struct {
	a logs_model.ConsoleLog
}

var systemAccount = logs_model.NewAccount("system", "")

func (s *SystemLogs) PushLog(msg string, group ...string) {
	err := s.a.Append(logs_model.NewMeta(systemAccount, "127.0.0.1", msg, group...))
	if err != nil {
		slog.Debug("审计写入失败", "err", err)
	}
}

func NewSystemLogs(a logs_model.ConsoleLog) *SystemLogs {
	return &SystemLogs{
		a: a,
	}
}
