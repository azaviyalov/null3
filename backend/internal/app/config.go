package app

import (
	"github.com/azaviyalov/null3/backend/internal/core/auth"
	"github.com/azaviyalov/null3/backend/internal/core/db"
	"github.com/azaviyalov/null3/backend/internal/core/frontend"
	"github.com/azaviyalov/null3/backend/internal/core/logging"
	"github.com/azaviyalov/null3/backend/internal/core/server"
)

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
