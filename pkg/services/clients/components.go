package clients

import (
	"github.com/console-dns/server/pkg/utils"
	"github.com/console-dns/server/pkg/utils/route"
	"github.com/pkg/errors"
)

func UiClientDelete(ctx *route.WebRequest) error {
	client, err := tryClient(ctx, false)
	if err != nil {
		return err
	}
	defer client.Close()
	return ctx.WriteComponentTemplate("clients/components/clients-delete", map[string]any{
		"Client": client.Name,
	})
}

func UiClientDetail(ctx *route.WebRequest) error {
	client, err := tryClient(ctx, false)
	if err != nil {
		return err
	}
	defer client.Close()
	return ctx.WriteComponentTemplate("clients/components/client-detail", map[string]any{
		"Name":   client.Name,
		"Client": client.Client,
		"Regex":  utils.DataRegex,
	})
}

func UiClientResetTokenDialog(ctx *route.WebRequest) error {
	client, err := tryClient(ctx, false)
	if err != nil {
		return err
	}
	defer client.Close()
	return ctx.WriteComponentTemplate("clients/components/client-reset-token-dialog", map[string]any{
		"Name": client.Name,
	})
}

func UiClientResetToken(ctx *route.WebRequest) error {
	client, err := tryClient(ctx, true)
	if err != nil {
		return err
	}
	defer client.Close()

	ctx.PushLog("重置客户端 $2 的 token", "clients", client.Name, "reset-token")
	token := client.ResetToken()
	return ctx.WriteComponentTemplate("clients/components/client-reset-token", map[string]any{
		"Name":  client.Name,
		"Token": token,
	})
}

func UiClientAllowIpDelete(ctx *route.WebRequest) error {
	client, err := tryClient(ctx, false)
	if err != nil {
		return err
	}
	defer client.Close()
	return ctx.WriteComponentTemplate("clients/components/client-delete-allow-ip", map[string]any{
		"Client": client.Name,
		"IP":     ctx.PathValue("ip"),
	})
}

func UiClientModRule(ctx *route.WebRequest) error {
	client, err := tryClient(ctx, false)
	if err != nil {
		return err
	}
	defer client.Close()
	ruleName := ctx.PathValue("rule")
	if ruleName != "" {
		rule := client.Rules[ruleName]
		if rule == nil {
			return errors.New("策略不存在")
		}
		return ctx.WriteComponentTemplate("clients/components/client-mod-rule", map[string]any{
			"Client":     client.Name,
			"IsUpdate":   true,
			"RuleName":   ruleName,
			"Host":       rule.Host,
			"Record":     rule.Record,
			"Permission": rule.Permission,
		})
	} else {
		return ctx.WriteComponentTemplate("clients/components/client-mod-rule", map[string]any{
			"Client":   client.Name,
			"IsUpdate": false,
			"Regex":    utils.DataRegex,
		})
	}
}

func UiClientRuleDelete(ctx *route.WebRequest) error {
	client, err := tryClient(ctx, false)
	if err != nil {
		return err
	}
	defer client.Close()
	ruleName := ctx.PathValue("rule")
	if client.Rules[ruleName] == nil {
		return errors.New("策略不存在")
	}
	return ctx.WriteComponentTemplate("clients/components/client-delete-rule", map[string]any{
		"Client":   client.Name,
		"RuleName": ruleName,
		"Rule":     client.Rules[ruleName],
	})
}

func UiClientRuleVerify(ctx *route.WebRequest) error {
	rule := ctx.PathValue("rule")
	client, err := tryClient(ctx, false)
	if err != nil {
		return err
	}
	defer client.Close()
	if client.Rules[rule] == nil {
		return errors.New("策略不存在")
	}
	return ctx.WriteComponentTemplate("clients/components/client-test-rule", map[string]any{
		"Client":   client.Name,
		"RuleName": rule,
		"Rule":     client.Rules[rule],
		"Regex":    utils.DataRegex,
	})
}

func UiClientRuleVerifyResult(ctx *route.WebRequest) error {
	rule := ctx.PathValue("rule")
	client, err := tryClient(ctx, false)
	if err != nil {
		return err
	}
	defer client.Close()
	if client.Rules[rule] == nil {
		return errors.New("策略不存在")
	}
	return ctx.WriteComponentTemplate("clients/components/client-test-rule-result", map[string]any{
		"Status": client.Rules[rule].VerifyRule(ctx.FormValue("zone"), ctx.FormValue("name"), ctx.FormValue("type")),
	})
}
