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
			s, err := server.New(&config.Config{
				Server: config.ServerOptions{
					RawAddress: "127.0.0.1",
					Network:    "udp",
					Port:       8533,
				},
				Resolver: config.ResolverOptions{
					Mode:      "r",
					MaxExpire: 300,
				},
			})
			if err != nil {
				return err
			}

			return s.ListenAndServe()
		},
	}

	return app.Run(os.Args)
}
