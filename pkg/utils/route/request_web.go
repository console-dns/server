package route

import (
	"net/http"
	"strconv"

	"github.com/console-dns/server/assets"
	"github.com/console-dns/server/pkg/content"
	"github.com/console-dns/server/pkg/models/auth"
	"github.com/console-dns/server/pkg/models/logs"
	"github.com/fatih/structs"
)

type WebHandler func(ctx *WebRequest) error

type WebRequest struct {
	*http.Request
	W          http.ResponseWriter
	Content    *content.Content
	LoginMeta  *auth.ConsoleLoginMeta
	selfHacker bool
	logs.LogPusher
	auth.Hacker
}

type TemplateData struct {
	Title string
	Path  string
	Data  any
}

func (c *WebRequest) WritePageTemplate(path string, args any) error {
	data := &TemplateData{
		Title: "",
		Data:  args,
		Path:  c.RequestURI,
	}
	if argMap, ok := args.(map[string]any); ok {
		if title, ok := argMap["Title"].(string); ok {
			data.Title = title
		}
	} else if args != nil {
		argMap := structs.Map(args)
		if title, ok := argMap["Title"].(string); ok {
			data.Title = title
		}
	}
	for _, p := range []string{"include/head", path, "include/body"} {
		if err := assets.WriteTemplate(c.W, p, data); err != nil {
			return err
		}
	}
	return nil
}

// WriteComponentTemplate 渲染模板内容
func (c *WebRequest) WriteComponentTemplate(path string, args any) error {
	return assets.WriteTemplate(c.W, path, args)
}

func (c *WebRequest) UIntForm(name string, def uint) (uint, error) {
	data := c.Request.FormValue(name)
	if data == "" {
		return def, nil
	}
	parseUint, err := strconv.ParseUint(data, 10, 0)
	return uint(parseUint), err
}

func (c *WebRequest) FormIfExists(name string, fallback func(value string)) {
	value := c.FormValue(name)
	if value != "" {
		fallback(value)
	}
}

func (c *WebRequest) FormIfExistsErr(name string, fallback func(value string) error) error {
	value := c.FormValue(name)
	if value != "" {
		return fallback(value)
	}
	return nil
}
