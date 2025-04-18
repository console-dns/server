package utils

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/pelletier/go-toml/v2"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

func AutoMarshal(filePath string, result any) error {
	var data []byte
	var err error
	if strings.HasSuffix(filePath, ".json") {
		data, err = json.Marshal(result)
	} else if strings.HasSuffix(filePath, ".yml") || strings.HasSuffix(filePath, ".yaml") {
		data, err = yaml.Marshal(result)
	} else if strings.HasSuffix(filePath, ".toml") {
		data, err = toml.Marshal(result)
	}
	if err != nil {
		return err
	}
	err = os.WriteFile(filePath, data, 0o644)
	if err != nil {
		return err
	}
	return nil
}

func AutoUnmarshal(filePath string, result any, generate bool) error {
	data, err := os.ReadFile(filePath)
	if err != nil && os.IsNotExist(err) && generate {
		err := AutoMarshal(filePath, result)
		if err != nil {
			return errors.Wrap(err, "文件不存在，且创建失败")
		}
	} else if err != nil {
		return err
	}
	if strings.HasSuffix(filePath, ".json") {
		return json.Unmarshal(data, result)
	} else if strings.HasSuffix(filePath, ".yml") || strings.HasSuffix(filePath, ".yaml") {
		return yaml.Unmarshal(data, result)
	} else if strings.HasSuffix(filePath, ".toml") {
		return toml.Unmarshal(data, result)
	}
	return errors.New("未知文件类型 : " + filePath)
}
