package config

import (
	"errors"
	"net"
	"os"

	"github.com/go-void/portal/pkg/constants"
	"github.com/go-void/portal/pkg/utils"

	"github.com/pelletier/go-toml/v2"
)

var (
	ErrInvalidResolverUpstream = errors.New("invalid resolver upstream")
	ErrInvalidCollectorBackend = errors.New("invalid collector backend")
	ErrInvalidServerAddress    = errors.New("invalid server address")
	ErrInvalidServerNetwork    = errors.New("invalid network")
	ErrInvalidResolverMode     = errors.New("invalid resolver mode")
	ErrInvalidLogMode          = errors.New("invalid log mode")
)

// NOTE (Techassi): Can we define the options in the packages itself?

// Config specifies the global configuration options
type Config struct {
	Collector CollectorOptions `toml:"collector"`
	Resolver  ResolverOptions  `toml:"resolver"`
	Filter    FilterOptions    `toml:"filter"`
	Server    ServerOptions    `toml:"server"`
	Log       LogOptions       `toml:"log"`
}

// CollectorOptions specifies available collector config options
type CollectorOptions struct {
	MaxEntries int    `toml:"max_entries"`
	Anonymize  bool   `toml:"anonymize"`
	Interval   uint   `toml:"interval"`
	Enabled    bool   `toml:"enabled"`
	Backend    string `toml:"backend"`
}

// ResolverOptions specifies available resolver config options
type ResolverOptions struct {
	CacheEnabled bool   `toml:"cache_enabled"`
	RawUpstream  string `toml:"upstream"`
	MaxExpire    int    `toml:"max_expire"`
	Upstream     net.IP `toml:"-"`
	HintPath     string `toml:"hint_path"`
	Mode         string `toml:"mode"`
}

// FilterOptions specifies available filter config options
type FilterOptions struct {
	TTL  int    `toml:"ttl"`
	Mode string `toml:"mode"`
}

// ServerOptions specifies available server config options
type ServerOptions struct {
	CacheEnabled bool   `toml:"cache_enabled"`
	RawAddress   string `toml:"address"`
	Address      net.IP `toml:"-"`
	Network      string `toml:"network"`
	Port         int    `toml:"port"`
}

// LogOptions specifies available log config options
type LogOptions struct {
	Enabled bool     `toml:"enabled"`
	Mode    string   `toml:"mode"`
	Level   string   `toml:"level"`
	Outputs []string `toml:"outputs"`
}

// Read reads a TOML config file and returns a Config or any error encountered
// while reading
func Read(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	c := new(Config)
	d := toml.NewDecoder(f)

	err = d.Decode(c)
	return c, err
}

// Write writes a TOML config to path and returns any error encountered while
// writing
func Write(path string, c *Config) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}

	e := toml.NewEncoder(f)
	return e.Encode(c)
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
		Log: LogOptions{
			Enabled: true,
			Mode:    "production",
			Level:   "error",
			Outputs: []string{"stdout"},
		},
	}
}

func (c *Config) Validate() error {
	aip := net.ParseIP(c.Server.RawAddress)
	if aip == nil {
		return ErrInvalidServerAddress
	}
	c.Server.Address = aip

	if utils.NotIn(c.Server.Network, []string{"udp", "udp4", "udp6", "tcp", "tcp4", "tcp6"}) {
		return ErrInvalidServerNetwork
	}

	if utils.NotIn(c.Resolver.Mode, []string{"r", "i", "f"}) {
		return ErrInvalidResolverMode
	}

	if c.Resolver.Mode == "f" {
		uip := net.ParseIP(c.Resolver.RawUpstream)
		if c.Resolver.RawUpstream == "" || uip == nil {
			return ErrInvalidResolverUpstream
		}
		c.Resolver.Upstream = uip
	}

	if c.Collector.Enabled && utils.NotIn(c.Collector.Backend, []string{"default", "mysql", "mariadb"}) {
		return ErrInvalidCollectorBackend
	}

	if c.Log.Enabled && utils.NotIn(c.Log.Mode, []string{"", "dev", "development", "prod", "production"}) {
		return ErrInvalidLogMode
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

	if c.Log.Level == "" {
		c.Log.Level = "error"
	}

	if len(c.Log.Outputs) == 0 {
		c.Log.Outputs = []string{"stderr"}
	}
}
