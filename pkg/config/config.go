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
	TTL           int    `toml:"ttl"`
	Mode          string `toml:"mode"`
	WatchURLS     bool   `toml:"watch_urls"`
	WatchFiles    bool   `toml:"watch_files"`
	WatchInterval int    `toml:"watch_interval"`
}

// Reads the contents of the config file and returns the options as a config struct
// func Read(path string) (*Config, error) {
// 	path, err := filepath.Abs(path)
// 	if err != nil {
// 		return nil, err
// 	}

// 	c := new(Config)
// 	b, err := os.ReadFile(path)
// 	if err != nil {
// 		return nil, err
// 	}

// 	err = toml.Unmarshal(b, c)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return c, nil
// }

// // Writes the config struct as a config file onto the disk
// func Write(path string, c *Config) error {
// 	file, err := os.Create(path)
// 	if err != nil {
// 		return err
// 	}

// 	b, err := toml.Marshal(c)
// 	if err != nil {
// 		return err
// 	}

// 	_, err = file.Write(b)
// 	return err
// }

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
