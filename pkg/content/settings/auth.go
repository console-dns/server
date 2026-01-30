package settings

import "time"

type AuthConfig struct {
	Username string `yaml:"username"` // 管理员静态用户名
	Password string `yaml:"password"` // 管理员密码 TODO: 切换为 argon2

	SessionTTL  time.Duration `yaml:"session_ttl"`  // 会话保持时间
	DenyTTL     time.Duration `yaml:"deny_ttl"`     // 黑名单 IP 保留时间
	HackerTTL   time.Duration `yaml:"hacker_ttl"`   // 风险 IP 保留时间
	HackerCount int           `yaml:"hacker_count"` // 风险 IP 封禁次数
	DenyIPs     []string      `yaml:"deny_ips"`     // 静态黑名单地址
}

func NewAuthConfig() *AuthConfig {
	return &AuthConfig{
		Username:    "admin",
		Password:    "password",
		HackerCount: 5,
		SessionTTL:  30 * time.Minute,
		DenyTTL:     60 * time.Minute,
		HackerTTL:   time.Hour,
		DenyIPs:     make([]string, 0),
	}
}
