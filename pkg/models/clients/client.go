package clients

import (
	"fmt"
	"net"
	"regexp"
	"time"

	"github.com/console-dns/server/pkg/content/settings"
	"github.com/console-dns/server/pkg/utils"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type Clients struct {
	localPath string
	Clients   map[string]*Client `json:"tokens" yaml:"tokens" toml:"tokens"`
}

func (c *Clients) ListByGroup(clientType ClientType) []string {
	r := make([]string, 0)
	for s, v := range c.Clients {
		if v.ClientType == clientType {
			r = append(r, s)
		}
	}
	return r
}

func (c *Clients) Get(name string) *Client {
	return c.Clients[name]
}

func (c *Clients) GetBySecret(secret string) (string, *Client) {
	if secret == "" {
		return "", nil
	}
	for n, client := range c.Clients {
		if client.Secret == secret {
			return n, client
		}
	}
	return "", nil
}

type ClientStatusResult struct {
	Name string
	At   time.Time
	IP   net.IP
	UA   string
}

func (c *Clients) GetClientStatus() []*ClientStatusResult {
	result := make([]*ClientStatusResult, 0)
	for name, client := range c.Clients {
		if client.Status != nil {
			result = append(result, &ClientStatusResult{
				Name: name,
				At:   client.Status.At,
				IP:   client.Status.IP,
				UA:   client.Status.UA,
			})
		}
	}
	return result
}

func FromClients(config *settings.StaticConfig) (*Clients, error) {
	result := &Clients{
		Clients: make(map[string]*Client),
	}
	err := utils.AutoUnmarshal(config.Storage.Token, result, true)
	if err != nil {
		return nil, errors.Wrapf(err, "解析 client 失败")
	}
	err = utils.AutoMarshal(config.Storage.Token, result)
	if err != nil {
		return nil, errors.Wrapf(err, "无法写入 client 配置")
	}
	result.localPath = config.Storage.Token
	return result, nil
}

func (c *Clients) Flush() error {
	return utils.AutoMarshal(c.localPath, c)
}

type Client struct {
	ClientType ClientType             `json:"type" yaml:"type" toml:"type"` // 客户端类型
	CreatedAt  time.Time              `json:"createdAt" yaml:"createdAt" toml:"createdAt"`
	Secret     string                 `yaml:"secret" json:"secret" toml:"secret"`                               // token 密钥
	AllowedIPs []net.IP               `yaml:"ips" json:"ips" toml:"ips"`                                        // 允许的 IP 地址 （白名单,如果为空则不判断）
	Rules      map[string]*ClientRule `yaml:"rules" json:"rules" toml:"rules"`                                  // 权限
	Status     *ClientStatus          `yaml:"status,omitempty" json:"status,omitempty" toml:"status,omitempty"` // 状态
}

type ClientStatus struct {
	At time.Time `yaml:"at" toml:"at" json:"at"`
	IP net.IP    `yaml:"ip" toml:"ip" json:"ip"`
	UA string    `yaml:"ua" toml:"ua" json:"ua"`
}

func (c *Clients) NewClient(name string, clientType ClientType) error {
	if err := utils.RegexIDName.Valid(name); err != nil {
		return err
	}
	if _, ok := c.Clients[name]; ok {
		return errors.New("client already exists")
	}
	c.Clients[name] = &Client{
		ClientType: clientType,
		CreatedAt:  time.Now(),
		Secret:     uuid.New().String(),
		AllowedIPs: make([]net.IP, 0),
		Rules:      make(map[string]*ClientRule),
	}
	return nil
}

func (c *Clients) Delete(name string) error {
	if _, ok := c.Clients[name]; !ok {
		return errors.New("client does not exist")
	}
	delete(c.Clients, name)
	return nil
}

func (c *Client) ResetToken() string {
	c.Secret = uuid.New().String()
	return c.Secret
}

func (c *Client) UpdateStatus(ip net.IP, ua string) {
	c.Status = &ClientStatus{
		At: time.Now(),
		IP: ip,
		UA: ua,
	}
}

func (c *Client) AddAllowedIp(ip net.IP) error {
	for _, p := range c.AllowedIPs {
		if p.Equal(ip) {
			return errors.New("IP 已存在")
		}
	}
	c.AllowedIPs = append(c.AllowedIPs, ip)
	return nil
}

func (c *Client) RemoveAllowedIp(ip net.IP) error {
	deleted := len(c.AllowedIPs)
	for i, p := range c.AllowedIPs {
		if p.Equal(ip) {
			deleted = i
		}
	}
	var err error
	c.AllowedIPs, err = utils.RemoveIndex(c.AllowedIPs, deleted)
	return err
}

func (c *Client) CheckAllowIP(ip net.IP) bool {
	if len(c.AllowedIPs) == 0 {
		return true
	}
	for _, allowedIP := range c.AllowedIPs {
		if allowedIP.Equal(ip) {
			return true
		}
	}
	return false
}

// VerifyRule 优化打分机制
func (c *Client) VerifyRule(zone, location, record string) Permission {
	result := NO
	targetHost := fmt.Sprintf("%s.%s", location, zone)
	for _, rule := range c.Rules {
		matchesHost, err := regexp.MatchString(rule.Host, targetHost)
		if err != nil {
			panic(err)
		}
		if matchesHost && (rule.Record == record || rule.Record == "*") {
			switch rule.Permission {
			case RO:
				if result == NO {
					result = RO
				}
			case RU:
				result = RU
			case RW:
				result = RW
				break
			}
		}
	}
	return result
}

func (c *Client) ModRule(name string, rule ClientRule) error {
	if err := utils.RegexIDName.Valid(name); err != nil {
		return err
	}
	r := &rule
	for n, clientRule := range c.Rules {
		if clientRule.Equals(r) {
			return errors.Errorf("已存在相同的策略 %s .", n)
		}
	}
	c.Rules[name] = r
	return nil
}

func (c *Client) DeleteRule(name string) error {
	if c.Rules[name] == nil {
		return errors.New("策略不存在")
	}
	delete(c.Rules, name)
	return nil
}
