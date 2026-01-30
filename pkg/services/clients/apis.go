package clients

import (
	"encoding/json"
	"net"

	self_errors "github.com/console-dns/server/pkg/errors"
	clients_model "github.com/console-dns/server/pkg/models/clients"
	"github.com/console-dns/server/pkg/utils"
	"github.com/console-dns/server/pkg/utils/route"
)

func ApiAddClient(ctx *route.ApiRequest) (any, error) {
	type req struct {
		Name string `json:"name"`
		Type string `json:"type"`
	}
	var r req
	if err := json.Unmarshal([]byte(ctx.Body), &r); err != nil {
		return nil, self_errors.BadRequestErrorf("请求解析失败: %s", err.Error())
	}
	if r.Type == "" {
		r.Type = string(clients_model.TypeClient)
	}

	cfg := ctx.Content.SyncTokens
	if err := cfg.NewClient(r.Name, clients_model.ClientType(r.Type)); err != nil {
		return nil, self_errors.BadRequestErrorf(err.Error())
	}
	ctx.PushLog("创建新的客户端: $2", "clients", r.Name, "create")
	return "创建成功", nil
}

func ApiDeleteClient(ctx *route.ApiRequest) (any, error) {
	client := ctx.PathValue("client")
	cfg := ctx.Content.SyncTokens
	if err := cfg.Delete(client); err != nil {
		return nil, self_errors.BadRequestErrorf(err.Error())
	}
	ctx.PushLog("删除已有客户端: $2", "clients", client, "remove")
	return "删除成功", nil
}

func ApiAddAllowedIp(ctx *route.ApiRequest) (any, error) {
	client := ctx.PathValue("client")
	type req struct {
		IP string `json:"ip"`
	}
	var r req
	if err := json.Unmarshal([]byte(ctx.Body), &r); err != nil {
		return nil, self_errors.BadRequestErrorf("请求解析失败: %s", err.Error())
	}

	err := utils.RegexIPv4.Valid(r.IP)
	if err != nil {
		err = utils.RegexIPv6.Valid(r.IP)
	}
	if err != nil {
		return nil, self_errors.BadRequestErrorf("未知 IP 类型")
	}
	parseIP := net.ParseIP(r.IP)
	if parseIP == nil {
		return nil, self_errors.BadRequestErrorf("解析 IP 失败")
	}

	cfg := ctx.Content.SyncTokens
	c := cfg.Get(client)
	if c == nil {
		return nil, self_errors.NotfoundErrorf("客户端不存在")
	}

	if err := c.AddAllowedIp(parseIP); err != nil {
		return nil, self_errors.BadRequestErrorf(err.Error())
	}
	ctx.PushLog("客户端 $2 新增 IP 放行白名单 $4", "clients", client, "allows-ip", r.IP, "create")
	return "添加成功", nil
}

func ApiDeleteAllowedIp(ctx *route.ApiRequest) (any, error) {
	client := ctx.PathValue("client")
	ip := ctx.PathValue("ip")
	parseIP := net.ParseIP(ip)
	if parseIP == nil {
		return nil, self_errors.BadRequestErrorf("解析 IP 失败")
	}

	cfg := ctx.Content.SyncTokens
	c := cfg.Get(client)
	if c == nil {
		return nil, self_errors.NotfoundErrorf("客户端不存在")
	}

	if err := c.RemoveAllowedIp(parseIP); err != nil {
		return nil, self_errors.BadRequestErrorf(err.Error())
	}
	ctx.PushLog("客户端 $2 移除 IP 放行白名单 $4", "clients", client, "allows-ip", ip, "remove")
	return "删除成功", nil
}

func ApiModRule(ctx *route.ApiRequest) (any, error) {
	client := ctx.PathValue("client")
	ruleName := ctx.PathValue("rule")

	type req struct {
		Name       string `json:"name"`
		Host       string `json:"host"`
		Type       string `json:"type"`
		Permission string `json:"permission"`
	}
	var r req
	if err := json.Unmarshal([]byte(ctx.Body), &r); err != nil {
		return nil, self_errors.BadRequestErrorf("请求解析失败: %s", err.Error())
	}

	if ruleName != "" {
		r.Name = ruleName
	}

	clientRule, err := clients_model.NewClientRule(r.Host, r.Type, r.Permission)
	if err != nil {
		return nil, self_errors.BadRequestErrorf(err.Error())
	}

	cfg := ctx.Content.SyncTokens
	c := cfg.Get(client)
	if c == nil {
		return nil, self_errors.NotfoundErrorf("客户端不存在")
	}

	if err = c.ModRule(r.Name, *clientRule); err != nil {
		return nil, self_errors.BadRequestErrorf(err.Error())
	}

	if ruleName != "" {
		ctx.PushLog("修改客户端 $2 策略 $4", "clients", client, "rules", r.Name, "mod")
	} else {
		ctx.PushLog("新增客户端 $2 策略 $4", "clients", client, "rules", r.Name, "create")
	}
	return "操作成功", nil
}

func ApiDeleteRule(ctx *route.ApiRequest) (any, error) {
	client := ctx.PathValue("client")
	ruleName := ctx.PathValue("rule")

	cfg := ctx.Content.SyncTokens
	c := cfg.Get(client)
	if c == nil {
		return nil, self_errors.NotfoundErrorf("客户端不存在")
	}

	if err := c.DeleteRule(ruleName); err != nil {
		return nil, self_errors.NotfoundErrorf(err.Error())
	}
	ctx.PushLog("删除客户端 $2 策略 $4", "clients", client, "rules", ruleName, "remove")
	return "删除成功", nil
}

func ApiVerifyRule(ctx *route.ApiRequest) (any, error) {
	client := ctx.PathValue("client")
	ruleName := ctx.PathValue("rule")

	type req struct {
		Zone string `json:"zone"`
		Name string `json:"name"`
		Type string `json:"type"`
	}
	var r req
	if err := json.Unmarshal([]byte(ctx.Body), &r); err != nil {
		return nil, self_errors.BadRequestErrorf("请求解析失败: %s", err.Error())
	}

	cfg := ctx.Content.SyncTokens
	c := cfg.Get(client)
	if c == nil {
		return nil, self_errors.NotfoundErrorf("客户端不存在")
	}

	rule := c.Rules[ruleName]
	if rule == nil {
		return nil, self_errors.NotfoundErrorf("策略不存在")
	}

	return rule.VerifyRule(r.Zone, r.Name, r.Type), nil
}

func ApiResetToken(ctx *route.ApiRequest) (any, error) {
	client := ctx.PathValue("client")
	cfg := ctx.Content.SyncTokens
	c := cfg.Get(client)
	if c == nil {
		return nil, self_errors.NotfoundErrorf("客户端不存在")
	}
	newToken := c.ResetToken()
	ctx.PushLog("重置客户端 $2 的 token", "clients", client, "reset-token")
	return newToken, nil
}
