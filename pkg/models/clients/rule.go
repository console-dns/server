package clients

import (
	"fmt"
	"regexp"
	"slices"

	"github.com/console-dns/server/pkg/models/zones"
	"github.com/pkg/errors"
)

// ClientRule token 策略
type ClientRule struct {
	Host       string     `yaml:"host"`       // 区域 (正则匹配)
	Record     string     `yaml:"record"`     // 可更改的类型 (正则匹配)
	Permission Permission `yaml:"permission"` // 权限
}

func NewClientRule(host, record, permission string) (*ClientRule, error) {
	_, err := regexp.Compile("^" + host + "$")
	if err != nil {
		return nil, err
	}
	if !slices.Contains(zones.RecordTypes, record) && record != "*" {
		return nil, errors.New("未知记录类型")
	}
	if permission != string(RO) && permission != string(RW) && permission != string(RU) {
		return nil, errors.New("未知权限类型")
	}
	return &ClientRule{
		Host:       host,
		Record:     record,
		Permission: Permission(permission),
	}, nil
}

// Permission 权限类型
// 在设计中 none 高于 rw 高于 ro
type Permission string

var (
	NO Permission = "no" // 无权限(只在内部使用)
	RO Permission = "ro" // 只读
	RU Permission = "ru" // 只读或更新
	RW Permission = "rw" // 读写
)

func (p Permission) CanRead() bool {
	switch p {
	case NO:
		return false
	case RW, RU, RO:
		return true
	}
	return false
}

func (p Permission) CanWrite() bool {
	switch p {
	case NO, RO:
		return false
	case RW, RU:
		return true
	}
	return false
}

func (p Permission) CanDelete() bool {
	switch p {
	case NO, RO, RU:
		return false
	case RW:
		return true
	}
	return false
}

func (t *ClientRule) VerifyRule(zone, location, record string) Permission {
	targetHost := fmt.Sprintf("%s.%s", location, zone)
	matchesHost, _ := regexp.MatchString(t.Host, targetHost)
	if !matchesHost {
		return NO
	}
	if matchesHost && (t.Record == record || t.Record == "*") {
		return t.Permission
	}
	return NO
}

func (t *ClientRule) Equals(other *ClientRule) bool {
	return t.Host == other.Host && t.Record == other.Record &&
		t.Permission == other.Permission
}
