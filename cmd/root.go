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
			cfg, err := config.Read(c.Args[0])
			if err != nil {
				return err
			}

			err = cfg.Validate()
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
