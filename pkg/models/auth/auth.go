package auth

//goland:noinspection GoSnakeCaseUsage
import (
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/console-dns/server/pkg/content/settings"
	logs_model "github.com/console-dns/server/pkg/models/logs"
	"github.com/console-dns/server/pkg/utils"
	"github.com/pkg/errors"
)

type Session struct {
	localPath   string
	sessionTtl  time.Duration
	ipTtl       time.Duration
	hackerTtl   time.Duration
	hackerCount int
	Sessions    map[string]*SessionState  `toml:"sessions" yaml:"sessions" json:"sessions"`       // 会话信息
	DenyIPs     map[string]*DenyIpState   `yaml:"deny_ips" json:"deny_ips" toml:"deny_ips"`       // ip 黑名单
	HackerIPs   map[string]*HackerIpState `yaml:"hacker_ips" json:"hacker_ips" toml:"hacker_ips"` // 风险用户
	Logs        logs_model.PushLog        `json:"-" yaml:"-" toml:"-" json:"-"`
}

func FromSession(config *settings.StaticConfig) (*Session, error) {
	result := &Session{
		Sessions:  make(map[string]*SessionState),
		DenyIPs:   make(map[string]*DenyIpState),
		HackerIPs: make(map[string]*HackerIpState),
	}
	err := utils.AutoUnmarshal(config.Storage.Session, result, true)
	if err != nil {
		return nil, errors.Wrap(err, "session 配置解析失败")
	}
	result.localPath = config.Storage.Session
	result.sessionTtl = config.Auth.SessionTTL
	result.ipTtl = config.Auth.DenyTTL
	result.hackerCount = config.Auth.HackerCount
	result.hackerTtl = config.Auth.HackerTTL
	// 清理旧的静态黑名单
	var oldDenyIps []string
	for key, state := range result.DenyIPs {
		if state.Static {
			oldDenyIps = append(oldDenyIps, key)
		}
	}
	for _, key := range oldDenyIps {
		delete(result.DenyIPs, key)
	}
	for _, ip := range config.Auth.DenyIPs {
		parseIP := net.ParseIP(ip)
		if parseIP == nil {
			return nil, errors.Errorf("错误的 IP 格式: %s", ip)
		}
		// 填充静态 IP
		result.DenyIPs[parseIP.String()] = NewDenyIp(true)
	}
	return result, nil
}

func (s *Session) Refresh() error {
	removedSession := make([]string, 0)
	for session, state := range s.Sessions {
		if time.Since(state.UpdateTime) > s.sessionTtl {
			removedSession = append(removedSession, session)
		}
	}
	removedIps := make([]string, 0)
	for addr, ipState := range s.DenyIPs {
		if !ipState.Static {
			if time.Since(ipState.CreateTime) > s.ipTtl {
				removedIps = append(removedIps, addr)
			}
		}
	}
	removedHackerIPs := make([]string, 0)
	for addr, ipState := range s.HackerIPs {
		if time.Since(ipState.UpdateTime) > s.hackerTtl {
			removedHackerIPs = append(removedHackerIPs, addr)
		}
	}
	if len(removedSession) != 0 {
		msg := fmt.Sprintf("清理 << %d >> 个过期用户会话", len(removedSession))
		s.Logs.PushLog(msg, "auth", "session", "logout")
		slog.Info("清理过期用户会话", "len", len(removedSession))
	}
	if len(removedIps) != 0 {
		msg := fmt.Sprintf("清理 IP 阻止名单  << %s >>", strings.Join(removedIps, ","))
		s.Logs.PushLog(msg, "auth", "ban-ip", "remove")
		slog.Info("清理 IP 阻止名单", "len", len(removedIps))
	}
	if len(removedHackerIPs) != 0 {
		slog.Info("清理 IP 风险名单", "len", len(removedHackerIPs))
	}
	for _, session := range removedSession {
		delete(s.Sessions, session)
	}
	for _, ip := range removedIps {
		delete(s.DenyIPs, ip)
	}
	for _, ip := range removedHackerIPs {
		delete(s.HackerIPs, ip)
	}
	return nil
}

func (s *Session) Flush() error {
	return utils.AutoMarshal(s.localPath, s)
}

func (s *Session) AddDenyIp(ip string) {
	if _, ok := s.DenyIPs[ip]; !ok {
		slog.Info("添加新的 IP 封禁", "ip", ip)
		s.Logs.PushLog("添加新的 IP 封禁 $4 ", "auth", "ban-ip", "add", ip)
		s.DenyIPs[ip] = NewDenyIp(false)
	}
}

func (s *Session) AddHackerCount(ip net.IP) {
	ipStr := ip.String()
	if it, ok := s.HackerIPs[ipStr]; ok {
		slog.Debug("更新风险IP危险等级", "ip", ip, "count", it.Count)
		it.Count = it.Count + 1
		it.UpdateTime = time.Now()
	} else {
		slog.Debug("新增风险 IP", "ip", ip)
		s.HackerIPs[ipStr] = &HackerIpState{
			Count:      1,
			UpdateTime: time.Now(),
		}
	}
	hackerIp := s.HackerIPs[ipStr]
	if hackerIp.Count >= s.hackerCount {
		delete(s.HackerIPs, ipStr)
		s.AddDenyIp(ipStr)
	}
}

type SessionState struct {
	Account    string    `yaml:"account"`     // 用户
	CreateTime time.Time `yaml:"create_time"` // 创建时间
	UpdateTime time.Time `yaml:"update_time"` // 更新时间
	UserAgent  string    `yaml:"user_agent"`  // 用户 UA
	IpAddr     net.IP    `yaml:"ip_addr"`     // 用户 IP
}

func NewSessionState(r *http.Request) (*SessionState, error) {
	return &SessionState{
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
		UserAgent:  r.UserAgent(),
		IpAddr:     utils.ParseAddress(r),
	}, nil
}

type DenyIpState struct {
	Static     bool      `yaml:"static"`      // 全局黑名单
	CreateTime time.Time `yaml:"create_time"` // 黑名单创建时间
}

func NewDenyIp(static bool) *DenyIpState {
	return &DenyIpState{
		Static:     static,
		CreateTime: time.Now(),
	}
}

type HackerIpState struct {
	Count      int       `yaml:"count" toml:"count" json:"count"`
	UpdateTime time.Time `yaml:"update_time" toml:"update_time" json:"update_time"`
}
