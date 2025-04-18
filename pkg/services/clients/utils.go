package clients

import (
	"github.com/console-dns/server/pkg/models/clients"
	"github.com/console-dns/server/pkg/utils/route"
	"github.com/pkg/errors"
)

type ClientStatus struct {
	Name string
	*clients.Client
	close func()
}

func (receiver *ClientStatus) Close() {
	receiver.close()
}

func tryClient(ctx *route.WebRequest, rw bool) (*ClientStatus, error) {
	client := ctx.PathValue("client")
	if client == "" {
		return nil, errors.New("客户端名称格式错误")
	}
	if rw {
		cfg, f := ctx.Content.SyncTokens.WithReadWrite()
		if cfg.Get(client) == nil {
			defer f()
			return nil, errors.New("客户端不存在")
		} else {
			return &ClientStatus{
				Name:   client,
				Client: cfg.Get(client),
				close:  f,
			}, nil
		}
	} else {
		cfg, f := ctx.Content.SyncTokens.WithReadOnly()
		if cfg.Get(client) == nil {
			defer f()
			return nil, errors.New("客户端不存在")
		} else {
			return &ClientStatus{
				Name:   client,
				Client: cfg.Get(client),
				close:  f,
			}, nil
		}
	}
}
