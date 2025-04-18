package settings

var (
	ConfigZone    = "conf/zones.yaml"
	ConfigSession = "conf/sessions.yaml"
	ConfigToken   = "conf/tokens.yaml"
	ConfigDDNS    = "conf/ddns.yaml"
	ConfigLog     = "conf/logs.csv"
)

// StorageConfig 配置持久化位置
type StorageConfig struct {
	Zone    string `json:"zone" yaml:"zone" json:"zone"`          // 区域保存位置
	Session string `json:"session" yaml:"session" json:"session"` // 会话保存位置
	Token   string `json:"token" yaml:"token" json:"token"`       // token 保存位置
	Log     string `json:"log" yaml:"log" json:"log"`             // 审计日志保存位置
}

func NewStorageConfig() *StorageConfig {
	return &StorageConfig{
		Zone:    ConfigZone,
		Session: ConfigSession,
		Token:   ConfigToken,
		Log:     ConfigLog,
	}
}
