package settings

import "fmt"

type ServerConfig struct {
	Addr string `json:"addr"` // 绑定地址
	Port uint16 `json:"port"` // 绑定端口
}

func NewServerConfig() *ServerConfig {
	return &ServerConfig{
		Addr: "0.0.0.0",
		Port: 8090,
	}
}

func (c *ServerConfig) AddressPort() string {
	return fmt.Sprintf("%s:%d", c.Addr, c.Port)
}
