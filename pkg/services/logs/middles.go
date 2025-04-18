package logs

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/console-dns/server/pkg/content"
	auth_model "github.com/console-dns/server/pkg/models/auth"
	logs_model "github.com/console-dns/server/pkg/models/logs"
	"github.com/console-dns/server/pkg/utils"
)

type RequestWrapper struct {
	logs_model.ConsoleLog
	ip   net.IP
	user *logs_model.Account
}

func (a *RequestWrapper) PushLog(msg string, groups ...string) {
	if a.user == nil {
		panic("no user :" + fmt.Sprintf("%s : %s", strings.Join(groups, "/"), msg))
	}
	_ = a.Append(logs_model.NewMeta(a.user, a.ip.String(), msg, groups...))
}

func (a *RequestWrapper) PushLogByUser(user *logs_model.Account, msg string, groups ...string) {
	_ = a.Append(logs_model.NewMeta(user, a.ip.String(), msg, groups...))
}

func DefaultLogsInject(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctn := content.GetContent(r)
		groups := &RequestWrapper{
			ConsoleLog: ctn.Logs,
			ip:         utils.ParseAddress(r),
			user:       nil,
		}
		ctx := context.WithValue(r.Context(), logs_model.CONTENT, groups)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func ConsoleRequestLogsInject(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wrapper := logs_model.GetContent(r).(*RequestWrapper)
		auth := r.Context().Value("auth").(*auth_model.ConsoleLoginMeta)
		wrapper.user = auth.User()
		ctx := context.WithValue(r.Context(), logs_model.CONTENT, wrapper)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func ApiLogsInject(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Context().Value("auth").(*auth_model.ConsoleApiMeta)
		wrapper := logs_model.GetContent(r).(*RequestWrapper)
		wrapper.user = auth.User()
		ctx := context.WithValue(r.Context(), logs_model.CONTENT, wrapper)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
