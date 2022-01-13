package cmd

import (
	"os"

	"github.com/go-void/portal/cmd/cli"
	"github.com/go-void/portal/pkg/config"
	"github.com/go-void/portal/pkg/server"
)

func Execute() error {
	app := &cli.App{
		Name:  "portal",
		Usage: "portal runs a DNS server",
		Action: func(c *cli.Context) error {
			cfg := &config.Config{
				Server: config.ServerOptions{
					CacheEnabled: true,
					RawAddress:   "127.0.0.1",
					Network:      "udp",
					Port:         8553,
				},
				Resolver: config.ResolverOptions{
					CacheEnabled: true,
					Mode:         "r",
				},
				Filter: config.FilterOptions{
					TTL:  300,
					Mode: "null",
				},
				Collector: config.CollectorOptions{
					Anonymize:  false,
					Enabled:    true,
					MaxEntries: 100000,
					Interval:   600,
					Backend:    "default",
				},
				Log: config.LogOptions{
					Mode:    "dev",
					Level:   "debug",
					Outputs: []string{"stdout"},
				},
			}

			err := cfg.Validate()
			if err != nil {
				return err
			}

			s := server.New(cfg)
			err = s.Run()
			if err != nil {
				return err
			}

			s.Wait()
			return nil
		},
	}

	return app.Run(os.Args)
}
