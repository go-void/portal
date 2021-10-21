package server

import "github.com/go-void/portal/internal/constants"

// Config describes available options to configure a
// DNS server instance
type Config struct {
	Address        string
	Network        string
	Port           int
	UDPMessageSize int
}

func (c *Config) Validate() error {
	if c.Address == "" {
		c.Address = "0.0.0.0"
	}

	if c.Network == "" {
		c.Network = "udp"
	}

	if c.Port == 0 {
		c.Port = 53
	}

	if c.UDPMessageSize == 0 {
		c.UDPMessageSize = constants.UDPMinMessageSize
	}

	return nil
}
