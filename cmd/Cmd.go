package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/MisterNorwood/SugarCube-Server/internal/utils"
	"github.com/urfave/cli/v3"
)

func Execute() *utils.SessionCtx {
	var SessionCtx utils.SessionCtx
	app := &cli.Command{
		Name:    "sugarcube-server",
		Usage:   "Server for the SugarCube coupon extension",
		Version: "v1.0.0-alpha.1",
		Flags: []cli.Flag{
			&cli.UintFlag{
				Name:    "db-port",
				Value:   27017,
				Usage:   "Mongod port",
				Aliases: []string{"d"},
				Action: func(ctx context.Context, cli *cli.Command, port uint64) error {
					uid := os.Geteuid()
					if port > 65535 {
						return errors.New("Invalid port: Must be a number between 1 and 65535")
					} else if port <= 1024 && uid != 0 {
						return errors.New("Invalid port: Ports 1-1024 require root")

					}

					return nil
				},
			},
			&cli.UintFlag{
				Name:    "port",
				Value:   80,
				Usage:   "Port for the program",
				Aliases: []string{"p"},
				Action: func(ctx context.Context, cli *cli.Command, port uint64) error {
					uid := os.Geteuid()
					if port > 65535 {
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
				Usage:   "Uatabase username",
				Aliases: []string{"u"},
			},
			&cli.StringFlag{
				Name:    "db-password",
				Usage:   "Database password",
				Aliases: []string{"P"},
			},
			&cli.BoolFlag{
				Name:    "debug",
				Usage:   "enable printing information",
				Value:   false,
				Aliases: []string{"D"},
			},
		},

		Action: func(ctx context.Context, cli *cli.Command) error {

			//No, there is no better way to do this

			dbPort, err := utils.CheckForEnv(utils.EnvDBPort, cli.Uint("db-port"))
			checkEnvErr(err)
			SessionCtx.DbPort = uint16(dbPort)

			webPort, err := utils.CheckForEnv(utils.EnvPort, cli.Uint("port"))
			checkEnvErr(err)
			SessionCtx.ServerPort = uint16(webPort)

			uri, err := utils.CheckForEnv(utils.EnvDBURI, cli.String("db-uri"))
			checkEnvErr(err)
			SessionCtx.DbUri = uri

			user, err := utils.CheckForEnv(utils.EnvDBUser, cli.String("db-user"))
			checkEnvErr(err)
			SessionCtx.DbUser = user

			pass, err := utils.CheckForEnv(utils.EnvDBPassword, cli.String("db-password"))
			checkEnvErr(err)
			SessionCtx.DbPassword = pass

			debug, err := utils.CheckForEnv(utils.EnvDebug, cli.Bool("debug"))
			checkEnvErr(err)
			SessionCtx.Debug = debug

			return nil
		},
	}
	if err := app.Run(context.Background(), os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	return &SessionCtx
}

func checkEnvErr(err error) {

	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

}
