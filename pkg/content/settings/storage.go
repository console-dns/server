package settings

// StorageConfig 配置持久化位置
type StorageConfig struct {
	DSN string `json:"dsn" yaml:"dsn" toml:"dsn"` // 存储 DSN
}

func NewStorageConfig() *StorageConfig {
	return &StorageConfig{
		DSN: "memory://",
	}
}
