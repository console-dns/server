package clients

import (
	"net"

	self_errors "github.com/console-dns/server/pkg/errors"
	clients_model "github.com/console-dns/server/pkg/models/clients"
	"github.com/console-dns/server/pkg/utils/route"
	spec_utils "github.com/console-dns/spec/utils"
	"github.com/pkg/errors"
)

func ClientAdd(ctx *route.WebRequest) error {
	cfg, f := ctx.Content.SyncTokens.WithReadWrite()
	defer f()
	if err := cfg.NewClient(ctx.FormValue("name"), clients_model.TypeClient); err != nil {
		return self_errors.BadRequestErrorf(err.Error())
	}
	ctx.PushLog("创建新的客户端: $2", "clients", ctx.FormValue("name"), "create")
	return self_errors.ErrorRedirect("/clients")
}

func ClientDel(ctx *route.WebRequest) error {
	cfg, f := ctx.Content.SyncTokens.WithReadWrite()
	defer f()
	if err := cfg.Delete(ctx.PathValue("client")); err != nil {
		return self_errors.BadRequestErrorf(err.Error())
	}
	ctx.PushLog("删除已有客户端: $2", "clients", ctx.PathValue("client"), "remove")
	return self_errors.ErrorRedirect("/clients")
}

func ClientAllowIpsAdd(ctx *route.WebRequest) error {
	ip := ctx.FormValue("ip")
	err := spec_utils.RegexIPv4.Valid(ip)
	if err != nil {
		err = spec_utils.RegexIPv6.Valid(ip)
	}
	if err != nil {
		return errors.New("未知 IP 类型")
	}
	parseIP := net.ParseIP(ip)
	if parseIP == nil {
		return errors.New("解析 IP 失败")
	}
	c, err := tryClient(ctx, true)
	if err != nil {
		return err
	}
	defer c.Close()
	if err := c.AddAllowedIp(parseIP); err != nil {
		return self_errors.BadRequestErrorf(err.Error())
	}
	ctx.PushLog("客户端 $2 新增 IP 放行白名单 $4", "clients", c.Name, "allows-ip", ip, "create")
	return self_errors.ErrorRedirect("/clients?tab=allow-ip&select=" + c.Name)
}

func ClientAllowIpsDel(ctx *route.WebRequest) error {
	ip := ctx.PathValue("ip")
	parseIP := net.ParseIP(ip)
	if parseIP == nil {
		return errors.New("解析 IP 失败")
	}
	c, err := tryClient(ctx, true)
	if err != nil {
		return err
	}
	defer c.Close()

	if err = c.RemoveAllowedIp(parseIP); err == nil {
		ctx.PushLog("客户端 $2 移除 IP 放行白名单 $4", "clients", c.Name, "allows-ip", ip, "remove")
		return self_errors.ErrorRedirect("/clients?tab=allow-ip&select=" + c.Name)
	}
	return err
}

func ClientRulesMod(ctx *route.WebRequest) error {
	name := ctx.FormValue("name")
	clientRule, err := clients_model.NewClientRule(
		ctx.FormValue("host"),
		ctx.FormValue("type"),
		ctx.FormValue("permission"),
	)
	if err != nil {
		return self_errors.BadRequestErrorf(err.Error())
	}
	c, err := tryClient(ctx, true)
	if err != nil {
		return err
	}
	defer c.Close()
	if ctx.PathValue("rule") != "" {
		name = ctx.PathValue("rule")
	}
	if err = c.ModRule(name, *clientRule); err != nil {
		return self_errors.BadRequestErrorf(err.Error())
	}
	if ctx.PathValue("rule") != "" {
		ctx.PushLog("修改客户端 $2 策略 $4", "clients", c.Name, "rules", name, "mod")
	} else {
		ctx.PushLog("新增客户端 $2 策略 $4", "clients", c.Name, "rules", name, "create")
	}
	return self_errors.ErrorRedirect("/clients?tab=rules&select=" + c.Name)
}

func ClientRulesDelete(ctx *route.WebRequest) error {
	rule := ctx.PathValue("rule")
	c, err := tryClient(ctx, true)
	if err != nil {
		return err
	}
	defer c.Close()
	if err = c.DeleteRule(rule); err != nil {
		return self_errors.NotfoundErrorf(err.Error())
	}
	ctx.PushLog("删除客户端 $2 策略 $4", "clients", c.Name, "rules", rule, "remove")
	return self_errors.ErrorRedirect("/clients?tab=rules&select=" + c.Name)
}
