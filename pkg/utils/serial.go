package utils

import (
	"context"
	"encoding/json"
	"os"
	"strings"

	"github.com/pelletier/go-toml/v2"
	"github.com/pkg/errors"
	"gopkg.d7z.net/middleware/kv"
	"gopkg.in/yaml.v3"
)

func AutoMarshal(filePath string, result any) error {
	var data []byte
	var err error
	ext := strings.ToLower(filePath[strings.LastIndex(filePath, "."):])
	switch ext {
	case ".json":
		data, err = json.Marshal(result)
	case ".yml", ".yaml":
		data, err = yaml.Marshal(result)
	case ".toml":
		data, err = toml.Marshal(result)
	default:
		return errors.New("未知文件类型 : " + filePath)
	}
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, data, 0o644)
}

func AutoUnmarshal(filePath string, result any, generate bool) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) && generate {
			if err := AutoMarshal(filePath, result); err != nil {
				return errors.Wrap(err, "文件不存在，且创建失败")
			}
			data, err = os.ReadFile(filePath)
		}
		if err != nil {
			return err
		}
	}
	ext := strings.ToLower(filePath[strings.LastIndex(filePath, "."):])
	switch ext {
	case ".json":
		return json.Unmarshal(data, result)
	case ".yml", ".yaml":
		return yaml.Unmarshal(data, result)
	case ".toml":
		return toml.Unmarshal(data, result)
	default:
		return errors.New("未知文件类型 : " + filePath)
	}
}

func AutoKVMarshal(ctx context.Context, storage kv.KV, key string, result any) error {
	data, err := json.Marshal(result)
	if err != nil {
		return err
	}
	return storage.Put(ctx, key, string(data), kv.TTLKeep)
}

func AutoKVUnmarshal(ctx context.Context, storage kv.KV, key string, result any) error {
	data, err := storage.Get(ctx, key)
	if err != nil {
		if err == kv.ErrKeyNotFound {
			return nil
		}
		return err
	}
	return json.Unmarshal([]byte(data), result)
}
