package auth

import (
	"log"
	"log/slog"
	"net/http"
	"strings"

	"github.com/alexedwards/argon2id"
	"github.com/console-dns/server/pkg/errors"
	"github.com/console-dns/server/pkg/models/auth"
	"github.com/console-dns/server/pkg/models/logs"
	"github.com/console-dns/server/pkg/utils/route"
	"github.com/google/uuid"
	"github.com/pquerna/otp/totp"
)

func LoginPost(ctx *route.WebRequest) error {
	totpCode := ctx.FormValue("2fa")
	authConfig := ctx.Content.Config.Auth
	username := ctx.FormValue("account")

	if username != authConfig.Username {
		ctx.TagHacker()
		ctx.PushLogByUser(logs.NewAccount("user", "guest"), "用户名不匹配 $4", "auth", "session", "login", username)
		return errors.ErrorRedirect("/login?error=auth")
	}

	if authConfig.TotpSecret != "" && totpCode == "" {
		return errors.ErrorRedirect("/login?error=2fa")
	}
	if authConfig.TotpSecret != "" && !totp.Validate(totpCode, authConfig.TotpSecret) {
		ctx.TagHacker()
		ctx.PushLogByUser(logs.NewAccount("user", "guest"), "totp 登录验证失败 $4", "auth", "session", "login", username)
		return errors.ErrorRedirect("/login?error=2fa-error")
	}
	password := ctx.FormValue("password")
	valid := false
	if strings.HasPrefix(authConfig.Password, "argon2:") {
		var err error
		valid, err = argon2id.ComparePasswordAndHash(password, strings.TrimPrefix(authConfig.Password, "argon2:"))
		if err != nil {
			valid = false
			slog.Debug("argon2 校验错误", "err", err)
		}
	} else {
		if authConfig.Password != "" && password == authConfig.Password {
			valid = true
		}
	}
	if !valid {
		ctx.TagHacker()
		ctx.PushLogByUser(logs.NewAccount("user", "guest"), "密码错误 $4", "auth", "session", "login", username)
		return errors.ErrorRedirect("/login?error=auth")
	}
	random, _ := uuid.NewRandom()
	token := random.String()
	var err error
	state, err := auth.NewSessionState(ctx.Request)
	state.Account = username
	if err != nil {
		return err
	}
	ctx.Content.SyncSessions.ReadWrite(func(session *auth.Session) {
		session.Sessions[token] = state
	})
	ctx.W.Header().Add("Set-Cookie", "session="+token)
	ctx.PushLogByUser(logs.NewAccount("user", username), "用户登录", "auth", "session", "login", username)
	log.Printf("新用户登录： %s, 来自: %s", token, state.IpAddr.String())
	return errors.ErrorRedirect("/")
}

func clearCookies(w http.ResponseWriter, r *http.Request) {
	if strings.Contains(r.Header.Get("Accept"), "json") {
		http.Error(w, "此请求已被禁止", 401)
	} else {
		w.Header().Set("Set-cookie", "session=deleted;expires=Thu, 01 Jan 1970 00:00:00 GMT")
		http.Redirect(w, r, "/login", 302)
	}
}

func LogoutGet(ctx *route.WebRequest) error {
	log.Printf("注销会话 %s", ctx.LoginMeta.Session)
	ctx.PushLog("用户 $4 注销会话", "auth", "session", "logout", ctx.LoginMeta.User().Name)
	ctx.Content.SyncSessions.ReadWrite(func(session *auth.Session) {
		delete(session.Sessions, ctx.LoginMeta.Session)
	})
	clearCookies(ctx.W, ctx.Request)
	return nil
}
