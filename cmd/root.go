package cmd

import (
	"os"

	"github.com/go-void/portal/cmd/cli"
	"github.com/go-void/portal/internal/server"
)

func Execute() error {
	app := &cli.App{
		Name:  "portal",
		Usage: "portal runs a DNS server",
		Action: func(c *cli.Context) error {
			s, err := server.New(&server.Config{
				Address: "127.0.0.1",
				Network: "udp",
				Port:    8533,
			})
			if err != nil {
				return err
			}

			return s.ListenAndServe()
		},
	}

	return app.Run(os.Args)
}
