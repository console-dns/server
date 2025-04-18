package logs

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type PushLog interface {
	// PushLog 推送审计日志
	PushLog(msg string, groups ...string)
}
type PushLogByUser interface {
	// PushLogByUser 推送自定义用户的审计日志
	PushLogByUser(user *Account, msg string, groups ...string)
}

// LogPusher 通用的审计处理
type LogPusher interface {
	PushLog
	PushLogByUser
}

// QueryResult 查询返回
type QueryResult struct {
	Count uint64  // 总数量
	Data  []*Meta // 数据
}

type QueryPage struct {
	Index  uint   // 索引
	Length uint64 // 数据数量
}

func WithPage(index uint, length uint64) *QueryPage {
	return &QueryPage{
		Index:  index,
		Length: length,
	}
}

type ConsoleLog interface {
	// Query 查询数据
	Query(page *QueryPage, filters ...Filter) (*QueryResult, error)
	// Append 追加日志
	Append(meta ...*Meta) error
}

type Filter func(meta *Meta) FilterStatus

// FilterStatus 过滤器状态
type FilterStatus int

var (
	STOP FilterStatus = -1 // 表明触发边界，结束查询
	OK   FilterStatus = 0  // 表明此内容符合查询
	SKIP FilterStatus = 1  // 跳过内容并继续下一条
)

func CheckFilters(filters []Filter, meta *Meta) FilterStatus {
	for _, executor := range filters {
		m := executor(meta)
		if m != OK {
			return m
		}
	}
	return OK
}

type Meta struct {
	CreatedAt time.Time // 创建时间
	Group     string    // 任务组，使用斜杠分割，一般为 group/action/value1/value2
	Author    *Account  // 任务角色
	Ip        string    // IP 地址
	Msg       string    // 任务动作
}

func NewMeta(user *Account, ip, msg string, groups ...string) *Meta {
	for i, group := range groups {
		msg = strings.ReplaceAll(msg, "$"+strconv.Itoa(i+1), fmt.Sprintf("<< %s >>", group))
	}
	return &Meta{
		CreatedAt: time.Now().UTC(),
		Group:     strings.Join(groups, "/"),
		Author:    user,
		Ip:        ip,
		Msg:       msg,
	}
}

func FromMeta(line string) (*Meta, error) {
	split := strings.SplitN(line, "|", 7)
	if len(split) != 5 {
		return nil, fmt.Errorf("invalid line format: %s", line)
	}
	t, err := strconv.ParseInt(split[0], 10, 64)
	if err != nil {
		return nil, err
	}
	account, err := FromAccount(split[2])
	if err != nil {
		return nil, err
	}
	data, err := base64.StdEncoding.DecodeString(split[4])
	if err != nil {
		return nil, err
	}
	return &Meta{
		time.UnixMilli(t),
		split[1],
		account,
		split[3],
		string(data),
	}, nil
}

func (m *Meta) String() string {
	return fmt.Sprintf("%d|%s|%s|%s|%s",
		m.CreatedAt.UnixMilli(),
		m.Group,
		m.Author,
		m.Ip,
		base64.StdEncoding.EncodeToString([]byte(m.Msg)),
	)
}

type Account struct {
	Type string
	Name string
}

func NewAccount(mType, name string) *Account {
	return &Account{
		Name: name,
		Type: mType,
	}
}

func FromAccount(str string) (*Account, error) {
	if strings.HasSuffix(str, "()") {
		return &Account{Type: strings.TrimSuffix(str, "()")}, nil
	}
	data := strings.SplitN(str, "(", 2)
	if len(data) != 2 || !strings.HasSuffix(str, ")") {
		return nil, errors.New("不是合法的用户类型: " + str)
	}
	return &Account{
		Type: data[0],
		Name: strings.TrimSuffix(data[1], ")"),
	}, nil
}

func (a *Account) String() string {
	return fmt.Sprintf("%s(%s)", a.Type, a.Name)
}
