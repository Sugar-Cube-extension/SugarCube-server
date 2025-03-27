package main

import (
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
	Session := cmd.Execute()
	Session.PrintEnv()
}

func Init(UserSession *utils.SessionCtx) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Info().Msg("Initializing application...")

}
