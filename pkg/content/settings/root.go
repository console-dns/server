package settings

import (
	"github.com/console-dns/server/pkg/utils"
)

// StaticConfig 只读配置
type StaticConfig struct {
	Server  *ServerConfig  `yaml:"server" toml:"server" json:"server"`    // 服务器配置
	Auth    *AuthConfig    `yaml:"auth" toml:"auth" json:"auth"`          // 认证配置
	Storage *StorageConfig `yaml:"storage" toml:"storage" json:"storage"` // 存储配置
}

func NewStaticConfig() *StaticConfig {
	return &StaticConfig{
		Server:  NewServerConfig(),
		Auth:    NewAuthConfig(),
		Storage: NewStorageConfig(),
	}
}

func FromStaticConfig(configPath string) (*StaticConfig, error) {
	result := NewStaticConfig()
	err := utils.AutoUnmarshal(configPath, result, false)
	if err != nil {
		return nil, err
	}
	return result, nil
}
