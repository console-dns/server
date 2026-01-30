package clients

import (
	"github.com/console-dns/server/pkg/models/clients"
	"github.com/console-dns/server/pkg/utils/route"
	"github.com/pkg/errors"
)

type ClientStatus struct {
	Name string
	*clients.Client
}

func (receiver *ClientStatus) Close() {
}

func tryClient(ctx *route.WebRequest, rw bool) (*ClientStatus, error) {
	client := ctx.PathValue("client")
	if client == "" {
		return nil, errors.New("客户端名称格式错误")
	}
	cfg := ctx.Content.SyncTokens
	c := cfg.Get(client)
	if c == nil {
		return nil, errors.New("客户端不存在")
	}
	return &ClientStatus{
		Name:   client,
		Client: c,
	}, nil
}
