package config

import (
	"errors"
	"net/netip"
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
	Store     StoreOptions     `toml:"store"`
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
	CacheEnabled bool       `toml:"cache_enabled"`
	RawUpstream  string     `toml:"upstream"`
	MaxExpire    int        `toml:"max_expire"`
	Upstream     netip.Addr `toml:"-"`
	HintPath     string     `toml:"hint_path"`
	Mode         string     `toml:"mode"`
}

// FilterOptions specifies available filter config options
type FilterOptions struct {
	TTL  int    `toml:"ttl"`
	Mode string `toml:"mode"`
}

// ServerOptions specifies available server config options
type ServerOptions struct {
	CacheEnabled bool           `toml:"cache_enabled"`
	Address      string         `toml:"address"`
	AddrPort     netip.AddrPort `toml:"-"`
	Network      string         `toml:"network"`
}

type StoreOptions struct {
	Username string `toml:"username"`
	Password string `toml:"password"`
	Database string `toml:"database"`
	Backend  string `toml:"backend"`
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
}

// LogOptions specifies available log config options
type LogOptions struct {
	Enabled bool     `toml:"enabled"`
	Mode    string   `toml:"mode"`
	Level   string   `toml:"level"`
	Outputs []string `toml:"outputs"`
}

// Read reads a TOML config file and returns a Config or any error encountered while reading
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

// Write writes a TOML config to path and returns any error encountered while writing
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
			Address:      "127.0.0.1:53",
			Network:      "udp",
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

// Validate validates the config and returns an error if the config is not valid
func (c *Config) Validate() error {
	addrPort, err := netip.ParseAddrPort(c.Server.Address)
	if err != nil {
		return ErrInvalidServerAddress
	}
	c.Server.AddrPort = addrPort

	if utils.NotIn(c.Server.Network, []string{"udp", "udp4", "udp6", "tcp", "tcp4", "tcp6"}) {
		return ErrInvalidServerNetwork
	}

	if utils.NotIn(c.Resolver.Mode, []string{"r", "i", "f"}) {
		return ErrInvalidResolverMode
	}

	if c.Resolver.Mode == "f" {
		addr, err := netip.ParseAddr(c.Resolver.RawUpstream)
		if err != nil {
			return ErrInvalidResolverUpstream
		}
		c.Resolver.Upstream = addr
	}

	if c.Collector.Enabled && utils.NotIn(c.Collector.Backend, []string{"default", "mysql", "mariadb"}) {
		return ErrInvalidCollectorBackend
	}

	if c.Log.Enabled && utils.NotIn(c.Log.Mode, []string{"", "dev", "development", "prod", "production"}) {
		return ErrInvalidLogMode
	}

	return nil
}

// Defaults sets sane defaults in the config
func (c *Config) Defaults() {
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
