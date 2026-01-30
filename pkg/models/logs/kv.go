package logs

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"sort"

	"gopkg.d7z.net/middleware/kv"
)

type KVLog struct {
	storage kv.KV
}

func NewKVLog(storage kv.KV) *KVLog {
	return &KVLog{
		storage: storage,
	}
}

func (a *KVLog) Query(page *QueryPage, filters ...Filter) (*QueryResult, error) {
	ctx := context.Background()
	// List all logs. This might be slow if there are many logs.
	list, err := a.storage.List(ctx, "logs/")
	if err != nil {
		return nil, err
	}

	keys := make([]string, 0, len(list))
	for k := range list {
		keys = append(keys, k)
	}
	// Sort by key (timestamp) ascending to match FileLog behavior (oldest first)
	sort.Strings(keys)

	allFound := make([]*Meta, 0)
	count := uint64(0)

l:
	for _, k := range keys {
		lineStr := list[k]
		meta, err := FromMeta(lineStr)
		if err != nil {
			slog.Debug("解析日志失败", "key", k, "err", err)
			continue
		}

		switch CheckFilters(filters, meta) {
		case OK:
			count = count + 1
			allFound = append(allFound, meta)
		case SKIP:
			continue l
		case STOP:
			break l
		}
	}

	// Now apply pagination on the filtered results, but we want newest first in the UI
	// So we reverse the results
	for i, j := 0, len(allFound)-1; i < j; i, j = i+1, j-1 {
		allFound[i], allFound[j] = allFound[j], allFound[i]
	}

	// Pagination
	start := uint64(page.Index) * page.Length
	end := start + page.Length
	if start > uint64(len(allFound)) {
		return &QueryResult{
			Count: count,
			Data:  []*Meta{},
		}, nil
	}
	if end > uint64(len(allFound)) {
		end = uint64(len(allFound))
	}

	return &QueryResult{
		Count: count,
		Data:  allFound[start:end],
	}, nil
}

func (a *KVLog) Append(meta ...*Meta) error {
	if len(meta) == 0 {
		return nil
	}
	ctx := context.Background()
	for _, m := range meta {
		key := fmt.Sprintf("logs/%020d_%08x", m.CreatedAt.UnixNano(), rand.Uint32())
		err := a.storage.Put(ctx, key, m.String(), kv.TTLKeep)
		if err != nil {
			return err
		}
	}
	return nil
}
