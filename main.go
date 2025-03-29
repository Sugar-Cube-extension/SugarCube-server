package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/MisterNorwood/SugarCube-Server/cmd"
	"github.com/MisterNorwood/SugarCube-Server/internal/middleware"
	"github.com/MisterNorwood/SugarCube-Server/internal/utils"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var (
	SessionCtx     *utils.SessionCtx
	DBClient       *mongo.Client
	WebServer      *echo.Echo
	ProgramContext *context.Context
)

func main() {
	session := cmd.Execute()
	SessionCtx = session
	SessionCtx.PrintEnv()
	err := Init(SessionCtx)
	if err != nil {
		os.Exit(1)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Kill)

	<-quit
	log.Warn().Msg("Shutting down server...")

	// Shutdown logic
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//FIXME: Yeah I don't know whats wrong atm
	// Gracefully shutdown Echo server
	if err := WebServer.Shutdown(shutdownCtx); err != nil {
		log.Error().Err(err).Msg("Error shutting down Echo server")
	} else {
		log.Info().Msg("Echo server shut down gracefully")
	}

	// Gracefully disconnect from MongoDB
	if err := DBClient.Disconnect(shutdownCtx); err != nil {
		log.Error().Err(err).Msg("Error closing MongoDB connection")
	} else {
		log.Info().Msg("MongoDB connection closed gracefully")
	}

	log.Warn().Msg("Application exited cleanly")
}

func Init(UserSession *utils.SessionCtx) error {
	// Context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ProgramContext = &ctx

	// Logger Setup
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout}) // Pretty print in dev

	log.Info().Str("version", "1.0.0").Str("hostname", getHostname()).Msg("Initializing application...")

	// MongoDB Setup
	client, err := mongo.Connect(options.Client().ApplyURI(SessionCtx.GetFullUri()))
	if err != nil {
		log.Error().Err(err).Msg("Failed to create MongoDB client")
		return err
	}
	DBClient = client

	err = DBClient.Ping(ctx, nil)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to Ping the database")
		return err
	}
	log.Debug().Msg("Successfuly connected to MongoDB server")

	// Echo Server Setup
	e := echo.New()
	e.Use(middleware.ZeroLogMiddleware)

	// Set up routes
	e.GET("/", func(c echo.Context) error {
		log.Info().Str("path", "/").Msg("Received request")
		return c.String(http.StatusOK, "Hello, World!")
	})

	port := strconv.FormatUint(uint64(SessionCtx.ServerPort), 10)
	go func() {
		log.Info().Str("port", port).Msg("Starting Echo web server")
		if err := e.Start(":" + port); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Failed to start Echo server")
		}
	}()

	WebServer = e
	return nil
}
func getHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	return hostname
}
