package cmd

import (
	"os"

	"github.com/go-void/portal/cmd/cli"
	"github.com/go-void/portal/internal/config"
	"github.com/go-void/portal/internal/server"
)

func Execute() error {
	app := &cli.App{
		Name:  "portal",
		Usage: "portal runs a DNS server",
		Action: func(c *cli.Context) error {
			cfg := &config.Config{
				Server: config.ServerOptions{
					RawAddress: "127.0.0.1",
					Network:    "tcp",
					Port:       8553,
				},
				Resolver: config.ResolverOptions{
					RawUpstream: "1.1.1.1",
					Mode:        "f",
				},
				Filter: config.FilterOptions{
					TTL:  300,
					Mode: "null",
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
