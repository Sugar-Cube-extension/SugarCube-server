package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/MisterNorwood/SugarCube-Server/cmd"
	"github.com/MisterNorwood/SugarCube-Server/internal/utils"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var (
	SessionCtx *utils.SessionCtx
	DBClient   *mongo.Client
	server     *echo.Echo
	Context    *context.Context
)

func main() {
	session := cmd.Execute()
	SessionCtx = session
	SessionCtx.PrintEnv()
	err := Init(SessionCtx)
	if err != nil {
		os.Exit(1)
	}

	//Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Warn().Msg("Shutting down server...")

	if err := DBClient.Disconnect(context.Background()); err != nil {
		log.Error().Err(err).Msg("Error closing MongoDB connection")
	}

	log.Warn().Msg("Application exited cleanly")
}

func Init(UserSession *utils.SessionCtx) error {
	//Context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	Context = &ctx

	//Logger
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout}) //Pretty print
	log.Info().Msg("Initializing application...")

	//Mongo
	client, err := mongo.Connect(options.Client().ApplyURI(SessionCtx.GetFullUri()))
	if err != nil {
		return err
	}
	DBClient = client

	err = DBClient.Ping(ctx, nil)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to Ping the database")
		return err
	}

	log.Debug().Msg("Successfully connected to MongoDB")

	return nil

}
