package config

import (
	"errors"
	"net"

	"github.com/go-void/portal/pkg/utils"
)

var (
	ErrInvalidResolverUpstream = errors.New("invalid resolver upstream")
	ErrInvalidServerAddress    = errors.New("invalid server address")
	ErrInvalidServerNetwork    = errors.New("invalid network")
	ErrInvalidResolverMode     = errors.New("no such resolver mode")
)

type Config struct {
	Server   ServerOptions   `toml:"server"`
	Resolver ResolverOptions `toml:"resolver"`
	Filter   FilterOptions   `toml:"filter"`
}

type ServerOptions struct {
	RawAddress string `toml:"address"`
	Address    net.IP `toml:"-"`
	Network    string `toml:"network"`
	Port       int    `toml:"port"`
}

type ResolverOptions struct {
	RawUpstream string `toml:"upstream"`
	Upstream    net.IP `toml:"-"`
	Mode        string `toml:"mode"`
	HintPath    string `toml:"hint_path"`
	MaxExpire   int    `toml:"max_expire"`
}

type FilterOptions struct {
	TTL  int    `toml:"ttl"`
	Mode string `toml:"mode"`
}

func (c *Config) Validate() error {
	aip := net.ParseIP(c.Server.RawAddress)
	if aip == nil {
		return ErrInvalidServerAddress
	}
	c.Server.Address = aip

	if !utils.In(c.Server.Network, []string{"udp", "udp4", "udp6", "tcp", "tcp4", "tcp6"}) {
		return ErrInvalidServerNetwork
	}

	if !utils.In(c.Resolver.Mode, []string{"r", "i", "f"}) {
		return ErrInvalidResolverMode
	}

	if c.Resolver.Mode == "f" {
		uip := net.ParseIP(c.Resolver.RawUpstream)
		if c.Resolver.RawUpstream == "" || uip == nil {
			return ErrInvalidResolverUpstream
		}
		c.Resolver.Upstream = uip
	}

	return nil
}

func (c *Config) Defaults() {
	if c.Server.Port <= 0 {
		c.Server.Port = 53
	}
}
