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
					RawAddress: "127.0.0.1",
					Network:    "udp",
					Port:       8553,
				},
				Resolver: config.ResolverOptions{
					RawUpstream: "1.1.1.1",
					Mode:        "r",
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
				},
			}

			err := cfg.Validate()
			if err != nil {
				return err
			}

			s := server.New()
			s.Configure(cfg)

			return s.ListenAndServe()
		},
	}

	return app.Run(os.Args)
}
