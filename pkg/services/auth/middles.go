package auth

//goland:noinspection GoSnakeCaseUsage
import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/console-dns/server/pkg/content"
	"github.com/console-dns/server/pkg/models/auth"
	"github.com/console-dns/server/pkg/utils"
)

type HackerMeta struct {
	isHacker bool
}

func (h *HackerMeta) TagHacker() {
	h.isHacker = true
}

func HackerIpsInject(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		meta := &HackerMeta{
			isHacker: false,
		}
		ctx := context.WithValue(r.Context(), "hacker", meta)
		next.ServeHTTP(w, r.WithContext(ctx))
		if meta.isHacker {
			ctn := content.GetContent(r)
			ctn.SyncSessions.ReadWrite(func(session *auth.Session) {
				session.AddHackerCount(utils.ParseAddress(r))
			})
		}
	})
}

func DenyIPsFilter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctn := content.GetContent(r)
		block := false
		ctn.SyncSessions.ReadOnly(func(session *auth.Session) {
			block = session.DenyIPs[utils.ParseAddress(r).String()] != nil
		})
		if block {
			http.Error(w, "你的 ip 位于位于禁用列表内，请联系管理员", http.StatusForbidden)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

func ConsoleLoginFilter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session")
		if err != nil || cookie.Value == "" {
			http.Redirect(w, r, "/login", 302)
			return
		}
		ip := utils.ParseAddress(r)
		ctn := content.GetContent(r)
		write, unlock := ctn.SyncSessions.WithReadWrite()
		currentSession := write.Sessions[cookie.Value]
		if currentSession != nil {
			if !currentSession.IpAddr.Equal(ip) {
				log.Printf("找到 session 但登录 IP 不正确")
				// 提前释放锁
				unlock()

				clearCookies(w, r)
				return
			}
			currentSession.UpdateTime = time.Now()
		} else {
			log.Printf("未找到 %s 的 session", cookie.Value)
			// 提前释放锁
			unlock()
			clearCookies(w, r)
			return
		}
		ctx := context.WithValue(r.Context(), "auth",
			auth.NewConsoleLoginMeta(currentSession.Account, cookie.Value, currentSession.CreateTime))
		// 提前释放锁
		unlock()
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func ApiLoginFilter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		token = strings.TrimPrefix(token, "Bearer ")
		ctn := content.GetContent(r)
		result, unlock := ctn.SyncTokens.WithReadWrite()
		user, client := result.GetBySecret(token)
		if client == nil {
			unlock()
			http.Error(w, "客户端不存在", http.StatusUnauthorized)
			return
		}
		ip := utils.ParseAddress(r)
		if !client.CheckAllowIP(ip) {
			unlock()
			http.Error(w, "客户端不在可信环境内", http.StatusForbidden)
			return
		}
		client.UpdateStatus(ip, r.UserAgent())
		ctx := context.WithValue(r.Context(), "auth", auth.NewConsoleApiMeta(user))
		unlock()
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
