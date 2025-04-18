package auth

import (
	"net/http"
	"time"

	"github.com/console-dns/server/pkg/models/logs"
)

type AccountMeta interface {
	User() *logs.Account
}

type ConsoleLoginMeta struct {
	user      string
	Session   string
	LoginTime time.Time
}

type ConsoleApiMeta struct {
	client string
}

func NewConsoleApiMeta(client string) *ConsoleApiMeta {
	return &ConsoleApiMeta{
		client,
	}
}

func GetLoginMeta(r *http.Request) *ConsoleLoginMeta {
	if r, ok := r.Context().Value("auth").(*ConsoleLoginMeta); ok {
		return r
	}
	return nil
}

func NewConsoleLoginMeta(user string, session string, loginTime time.Time) *ConsoleLoginMeta {
	return &ConsoleLoginMeta{
		user:      user,
		Session:   session,
		LoginTime: loginTime,
	}
}

func GetApiMeta(r *http.Request) *ConsoleApiMeta {
	if r, ok := r.Context().Value("auth").(*ConsoleApiMeta); ok {
		return r
	}
	return nil
}

func (c *ConsoleApiMeta) User() *logs.Account {
	return logs.NewAccount("client", c.client)
}

func (c *ConsoleLoginMeta) User() *logs.Account {
	return logs.NewAccount("user", c.user)
}
