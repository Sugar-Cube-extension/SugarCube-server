package main

import (
	"os"

	"github.com/MisterNorwood/SugarCube-Server/cmd"
	"github.com/MisterNorwood/SugarCube-Server/internal/utils"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var (
	SessionCtx *utils.SessionCtx
	db         *mongo.Client
	server     *echo.Echo
)

func main() {
	session := cmd.Execute()
	SessionCtx = session
	SessionCtx.PrintEnv()
	err := Init(SessionCtx)
	if err != nil {
		os.Exit(1)
	}
}

func Init(UserSession *utils.SessionCtx) error {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout}) //Pretty print
	log.Info().Msg("Initializing application...")
	return nil

}
