package logs

import (
	"strconv"
	"strings"
	"time"

	logs_model "github.com/console-dns/server/pkg/models/logs"
	"github.com/console-dns/server/pkg/utils/route"
	"github.com/thoas/go-funk"
)

func QueryByGroup(ctx *route.WebRequest) error {
	logs := ctx.Content.Logs
	index, err := ctx.UIntForm("index", 0)
	if err != nil {
		return err
	}
	length, err := ctx.UIntForm("length", 15)
	if err != nil {
		return err
	}
	// 构建过滤器
	params := make([]string, 0)
	filters := make([]logs_model.Filter, 0)
	ctx.FormIfExists("group", func(value string) {
		params = append(params, "group="+value)
		filters = append(filters, logs_model.QueryGroup(value))
	})
	if err = ctx.FormIfExistsErr("begin", func(value string) error {
		date, err := time.ParseInLocation("2006-01-02T15:04", value, time.Local)
		if err != nil {
			return err
		}
		params = append(params, "begin="+value)
		filters = append(filters, logs_model.QueryAfter(date.UTC()))
		return nil
	}); err != nil {
		return err
	}
	if err = ctx.FormIfExistsErr("end", func(value string) error {
		date, err := time.ParseInLocation("2006-01-02T15:04", value, time.Local)
		if err != nil {
			return err
		}
		params = append(params, "end="+value)
		filters = append(filters, logs_model.QueryBefore(date.UTC()))
		return nil
	}); err != nil {
		return err
	}
	ctx.FormIfExists("account", func(value string) {
		params = append(params, "account="+value)
		filters = append(filters, logs_model.QueryAuthor(value))
	})
	ctx.FormIfExists("ip", func(value string) {
		params = append(params, "ip="+value)
		filters = append(filters, logs_model.QueryIpAddr(value))
	})
	ctx.FormIfExists("msg", func(value string) {
		params = append(params, "msg="+value)
		filters = append(filters, logs_model.QueryMsg(value))
	})
	page := logs_model.WithPage(index, uint64(length))
	query, err := logs.Query(
		page,
		filters...,
	)
	if err != nil {
		return err
	}
	for _, datum := range query.Data {
		datum.Msg = strings.ReplaceAll(strings.ReplaceAll(datum.Msg, "<<", "<code>"), ">>", "</code>")
	}
	pageSize := query.Count / page.Length
	if pageSize*page.Length < query.Count {
		pageSize = pageSize + 1
	}
	size := int(pageSize)
	depth := 2
	pages := make([]int, 0)
	for i := 0; i < size; i++ {
		if (i >= (int(index)-depth) && i <= (int(index)+depth)) ||
			(int(index) <= depth && i <= depth*2) ||
			(int(index) >= size-depth && i >= size-depth*2-1) {
			pages = append(pages, i)
		}
	}
	return ctx.WriteComponentTemplate("logs/components/log-details", map[string]any{
		"Data":       query.Data,
		"Params":     strings.Join(params, "&"),
		"Count":      int(query.Count),
		"PageIndex":  index,
		"PageLength": length,
		"PageSize":   int(pageSize),
		"Pages":      pages,
		"HtmlId":     strconv.Itoa(funk.RandomInt(100, 200)),
	})
}
