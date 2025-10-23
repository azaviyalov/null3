package main

import (
	"context"
	"os"
	"strconv"

	"github.com/azaviyalov/null3/backend/internal/core/auth"
	"github.com/azaviyalov/null3/backend/internal/core/db"
	"github.com/azaviyalov/null3/backend/internal/core/frontend"
	"github.com/azaviyalov/null3/backend/internal/core/logging"
	"github.com/azaviyalov/null3/backend/internal/core/server"
	"github.com/azaviyalov/null3/backend/internal/mood"
	"github.com/joho/godotenv"
)

func main() {
	if strconv.IntSize != 64 {
		logging.Error(context.Background(), "unsupported architecture: only 64-bit systems are supported")
		os.Exit(1)
	}

	// Load environment variables from .env file
	_ = godotenv.Load()

	config, err := GetConfig()
	if err != nil {
		logging.Error(context.Background(), "failed to get configuration", "error", err)
		os.Exit(1)
	}

	logging.Setup(config.Logging)

	database, err := db.Setup(config.DB)
	if err != nil {
		logging.Error(context.Background(), "failed to setup database", "error", err)
		os.Exit(1)
	}

	e := server.NewEchoServer(config.Server)

	frontend.InitModule(e, config.Frontend)
	authModule := auth.InitModule(e, database, config.Auth, config.StubUserConfig)

	mood.InitModule(e, database, authModule)

	if err := server.StartServer(e, config.Server); err != nil {
		logging.Error(context.Background(), "general server failure", "error", err)
		os.Exit(1)
	}
}

type Config struct {
	Auth           auth.Config
	StubUserConfig auth.StubUserConfig
	DB             db.Config
	Frontend       frontend.Config
	Logging        logging.Config
	Server         server.Config
}

func GetConfig() (Config, error) {
	authConfig, err := auth.GetConfig()
	if err != nil {
		return Config{}, err
	}

	stubUserConfig := auth.GetStubUserConfig()

	dbConfig := db.GetConfig()

	frontendConfig, err := frontend.GetConfig()
	if err != nil {
		return Config{}, err
	}

	loggingConfig := logging.GetConfig()

	serverConfig, err := server.GetConfig()
	if err != nil {
		return Config{}, err
	}

	return Config{
		Auth:           authConfig,
		StubUserConfig: stubUserConfig,
		DB:             dbConfig,
		Frontend:       frontendConfig,
		Logging:        loggingConfig,
		Server:         serverConfig,
	}, nil
}
