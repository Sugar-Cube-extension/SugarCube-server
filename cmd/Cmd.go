package cmd

import (
	"context"
	"errors"
	"fmt"
	"github.com/urfave/cli/v3"
	"os"
)

func Execute() {
	app := &cli.Command{
		Name:    "sugarcube-server",
		Usage:   "Server for the SugarCube coupon extension",
		Version: "v1.0.0-alpha.1",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:    "db-port",
				Value:   27017,
				Usage:   "Mongod port",
				Aliases: []string{"d"},
				Action: func(ctx context.Context, cli *cli.Command, port int64) error {
					uid := os.Geteuid()
					if port < 1 || port > 65535 {
						return errors.New("Invalid port: Must be a number between 1 and 65535")
					} else if port <= 1024 && uid != 0 {
						return errors.New("Invalid port: Ports 1-1024 require root")

					}

					return nil
				},
			},
			&cli.IntFlag{
				Name:    "port",
				Value:   80,
				Usage:   "Port for the program",
				Aliases: []string{"p"},
				Action: func(ctx context.Context, cli *cli.Command, port int64) error {
					uid := os.Geteuid()
					if port < 1 || port > 65535 {
						return errors.New("Invalid port: Must be a number between 1 and 65535")
					} else if port <= 1024 && uid != 0 {
						return errors.New("Invalid port: Ports 1-1024 require root")

					}

					return nil
				},
			},
			&cli.StringFlag{
				Name:    "db-uri",
				Value:   `mongodb://localhost`,
				Usage:   "Port for the program",
				Aliases: []string{"U"},
			},
			&cli.StringFlag{
				Name:    "db-user",
				Usage:   "Port for the program",
				Aliases: []string{"u"},
			},
			&cli.StringFlag{
				Name:    "db-password",
				Usage:   "Port for the program",
				Aliases: []string{"P"},
			},
		},

		Action: func(ctx context.Context, cli *cli.Command) error {
			fmt.Println(cli.Int("db-port"))

			return nil
		},
	}
	if err := app.Run(context.Background(), os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
