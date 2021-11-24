package config

import (
	"errors"
	"net"

	"github.com/go-void/portal/pkg/constants"
	"github.com/go-void/portal/pkg/utils"
)

var (
	ErrInvalidResolverUpstream = errors.New("invalid resolver upstream")
	ErrInvalidCollectorBackend = errors.New("invalid collector backend")
	ErrInvalidServerAddress    = errors.New("invalid server address")
	ErrInvalidServerNetwork    = errors.New("invalid network")
	ErrInvalidResolverMode     = errors.New("no such resolver mode")
)

// NOTE (Techassi): Can we define the options in the packages itself?

type Config struct {
	Collector CollectorOptions `toml:"collector"`
	Resolver  ResolverOptions  `toml:"resolver"`
	Filter    FilterOptions    `toml:"filter"`
	Server    ServerOptions    `toml:"server"`
}

type CollectorOptions struct {
	MaxEntries int    `toml:"max_entries"`
	Anonymize  bool   `toml:"anonymize"`
	Interval   uint   `toml:"interval"`
	Enabled    bool   `toml:"enabled"`
	Backend    string `toml:"backend"`
}

type ResolverOptions struct {
	CacheEnabled bool   `toml:"cache_enabled"`
	RawUpstream  string `toml:"upstream"`
	MaxExpire    int    `toml:"max_expire"`
	Upstream     net.IP `toml:"-"`
	HintPath     string `toml:"hint_path"`
	Mode         string `toml:"mode"`
}

type FilterOptions struct {
	TTL  int    `toml:"ttl"`
	Mode string `toml:"mode"`
}

type ServerOptions struct {
	CacheEnabled bool   `toml:"cache_enabled"`
	RawAddress   string `toml:"address"`
	Address      net.IP `toml:"-"`
	Network      string `toml:"network"`
	Port         int    `toml:"port"`
}

// Default returns a config with default values
func Default() *Config {
	return &Config{
		Server: ServerOptions{
			CacheEnabled: true,
			Address:      net.ParseIP("127.0.0.1"),
			Network:      "udp",
			Port:         53,
		},
		Resolver: ResolverOptions{
			CacheEnabled: true,
			Mode:         "r",
			HintPath:     "",
			MaxExpire:    300,
		},
		Filter: FilterOptions{
			TTL:  0,
			Mode: "null",
		},
		Collector: CollectorOptions{
			MaxEntries: 1000,
			Anonymize:  false,
			Interval:   900,
			Enabled:    true,
			Backend:    "default",
		},
	}
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

	if c.Collector.Enabled && !utils.In(c.Collector.Backend, []string{"default", "mysql", "mariadb"}) {
		return ErrInvalidCollectorBackend
	}

	return nil
}

func (c *Config) Defaults() {
	if c.Server.Port <= 0 {
		c.Server.Port = 53
	}

	if c.Collector.MaxEntries <= 0 {
		c.Collector.MaxEntries = constants.CollectorDefaultMaxEntries
	}

	if c.Collector.Interval <= 0 {
		c.Collector.Interval = constants.CollectorDefaultInterval
	}
}
