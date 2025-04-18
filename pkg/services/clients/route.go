package clients

import (
	"github.com/console-dns/server/pkg/utils/route"
	"github.com/go-chi/chi/v5"
)

func Route(r chi.Router) {
	r.Get("/", route.MVCBind(ClientRootPage))
	r.Post("/", route.MVCBind(ClientAdd))
	r.Post("/{client}/delete", route.MVCBind(ClientDel))
	r.Post("/{client}/allow-ip", route.MVCBind(ClientAllowIpsAdd))
	r.Post("/{client}/allow-ip/{ip}", route.MVCBind(ClientAllowIpsDel))
	r.Post("/{client}/rules", route.MVCBind(ClientRulesMod))
	r.Post("/{client}/rules/{rule}", route.MVCBind(ClientRulesMod))
	r.Post("/{client}/rules/{rule}/delete", route.MVCBind(ClientRulesDelete))

	r.Get("/!/clients/{client}", route.MVCBind(UiClientDetail))
	r.Get("/!/clients/{client}/resetToken", route.MVCBind(UiClientResetTokenDialog))
	r.Post("/!/clients/{client}/resetToken", route.MVCBind(UiClientResetToken))
	r.Delete("/!/clients/{client}", route.MVCBind(UiClientDelete))
	r.Delete("/!/clients/{client}/allow-ip/{ip}", route.MVCBind(UiClientAllowIpDelete))
	r.Get("/!/clients/{client}/rules", route.MVCBind(UiClientModRule))
	r.Get("/!/clients/{client}/rules/{rule}", route.MVCBind(UiClientModRule))
	r.Get("/!/clients/{client}/rules/{rule}/verify", route.MVCBind(UiClientRuleVerify))
	r.Post("/!/clients/{client}/rules/{rule}/verify", route.MVCBind(UiClientRuleVerifyResult))
	r.Delete("/!/clients/{client}/rules/{rule}", route.MVCBind(UiClientRuleDelete))
}
