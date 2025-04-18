package assets

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"mime"
	"path/filepath"
	"reflect"
	"sync"
	"text/template"
	"time"

	"github.com/Masterminds/sprig/v3"
)

var (
	StaticData EmbedFS
	Cache      bool
	BuildTime  = time.Now()
)

type FileLastMod interface {
	ModTime(path string) (time.Time, error)
}

type EmbedFS interface {
	FileLastMod
	fs.ReadFileFS
	fs.ReadDirFS
}

var tmplCache = new(sync.Map)

func openTemplate(path string) (*template.Template, error) {
	if data, ok := tmplCache.Load(path); ok {
		return data.(*template.Template), nil
	}
	tmplPath := fmt.Sprintf("template/%s.gohtml", path)
	tmplBuf, err := StaticData.ReadFile(tmplPath)
	if err != nil {
		return nil, err
	}
	content, err := template.New(tmplPath).Funcs(sprig.TxtFuncMap()).Funcs(
		map[string]any{
			"include": func(path string, data any) (string, error) {
				buffer := &bytes.Buffer{}
				err := WriteTemplate(buffer, path, data)
				body := buffer.String()
				return body, err
			},
			"avail": func(name string, data interface{}) bool {
				v := reflect.ValueOf(data)
				if v.Kind() == reflect.Ptr {
					v = v.Elem()
				}
				if v.Kind() == reflect.Struct {
					return v.FieldByName(name).IsValid()
				}
				if v.Kind() == reflect.Map {
					return data.(map[string]interface{})[name] != nil
				}
				return false
			},
			"foreach": func(from, to int) <-chan int {
				ch := make(chan int)
				go func() {
					for i := from; i < to; i++ {
						ch <- i
					}
					close(ch)
				}()
				return ch
			},
		}).Parse(string(tmplBuf))
	if err != nil {
		return nil, err
	}
	if Cache {
		tmplCache.Store(path, content)
	}
	return content, nil
}

func WriteTemplate(writer io.Writer, path string, data any) error {
	content, err := openTemplate(path)
	if err != nil {
		return err
	}
	return content.Execute(writer, data)
}

type FileInfo struct {
	Data    []byte
	ModTime time.Time
	Mine    string
}

func OpenResFile(path string) (*FileInfo, error) {
	data, err := StaticData.ReadFile(path)
	if err != nil {
		return nil, err
	}
	modTime, _ := StaticData.ModTime(path)
	return &FileInfo{
		Data:    data,
		ModTime: modTime,
		Mine:    mime.TypeByExtension(filepath.Ext(path)),
	}, nil
}
