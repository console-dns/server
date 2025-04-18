//go:build dev

package assets

import (
	"io/fs"
	"log/slog"
	"os"
	"time"
)

func init() {
	slog.Warn("使用文件系统模板")
	Cache = false
	StaticData = localFs{
		path: "assets",
	}
}

type localFs struct {
	path string
}

func (l localFs) Open(name string) (fs.File, error) {
	return os.DirFS(l.path).Open(name)
}

func (l localFs) ReadFile(name string) ([]byte, error) {
	return os.DirFS(l.path).(fs.ReadFileFS).ReadFile(name)
}

func (l localFs) ReadDir(name string) ([]fs.DirEntry, error) {
	return os.DirFS(l.path).(fs.ReadDirFS).ReadDir(name)
}

func (l localFs) ModTime(path string) (time.Time, error) {
	stat, err := os.Stat("assets" + "/" + path)
	if err != nil {
		return time.Now(), err
	}
	return stat.ModTime(), err
}
