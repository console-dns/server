package logs

import (
	"bufio"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
	"sync"

	"github.com/pkg/errors"
)

type FileLog struct {
	file string
	lock sync.RWMutex
}

func NewFileLog(file string) (*FileLog, error) {
	slog.Debug("新增审计持久化", "path", file)
	_, err := os.Stat(file)
	if err != nil && os.IsNotExist(err) {
		err := os.WriteFile(file, []byte(""), 0o600)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}
	openFile, err := os.OpenFile(file, os.O_RDWR, 0o600)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open file")
	}
	_ = openFile.Close()
	return &FileLog{
		file: file,
		lock: sync.RWMutex{},
	}, nil
}

func (a *FileLog) Query(page *QueryPage, filters ...Filter) (*QueryResult, error) {
	a.lock.RLock()
	defer a.lock.RUnlock()

	result := make([]*Meta, 0)
	file, err := os.Open(a.file)
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()
	reader := bufio.NewReader(file)

	seekSize := uint64(page.Index) * page.Length
	contentSize := page.Length
	count := uint64(0)
l:
	for {

		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		if len(filters) == 0 {
			continue l
		}
		if err != nil {
			return nil, err
		}
		lineStr := string(line)
		if lineStr == "" {
			break
		}
		meta, err := FromMeta(lineStr)
		if err != nil {
			return nil, err
		}
		switch CheckFilters(filters, meta) {
		case OK:
			count = count + 1
			if seekSize == 0 {
				if contentSize > 0 {
					result = append(result, meta)
					contentSize = contentSize - 1
				}
			} else {
				seekSize = seekSize - 1
			}
		case SKIP:
			continue l
		case STOP:
			break l
		}
	}
	return &QueryResult{
		Count: count,
		Data:  result,
	}, nil
}

func (a *FileLog) Append(meta ...*Meta) error {
	if len(meta) == 0 {
		return nil
	}
	a.lock.Lock()
	defer a.lock.Unlock()

	file, err := os.OpenFile(a.file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		return err
	}
	content := make([]string, len(meta))
	for i, m := range meta {
		slog.Debug("新增审计数据", "line", m)
		content[i] = m.String()
	}
	_, err = file.WriteString(fmt.Sprintf("%s\r\n", strings.Join(content, "\r\n")))
	if err != nil {
		return err
	}
	return nil
}
