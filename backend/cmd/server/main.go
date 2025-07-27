package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/azaviyalov/null3/backend/internal/core/db"
	"github.com/azaviyalov/null3/backend/internal/core/logging"
	"github.com/azaviyalov/null3/backend/internal/core/server"
	"github.com/azaviyalov/null3/backend/internal/frontend"
	"github.com/azaviyalov/null3/backend/internal/mood"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	_ = godotenv.Load()

	// Set config defaults if not set
	if os.Getenv("PORT") == "" {
		os.Setenv("PORT", "8080")
	}
	if os.Getenv("DATABASE_URL") == "" {
		os.Setenv("DATABASE_URL", "file:null3.db?_fk=1")
	}
	if os.Getenv("ENABLE_FRONTEND_DIST") == "" {
		os.Setenv("ENABLE_FRONTEND_DIST", "false")
	}
	if os.Getenv("API_URL") == "" {
		os.Setenv("API_URL", fmt.Sprintf("http://localhost:%s/api", os.Getenv("PORT")))
	}

	logging.Setup()

	// Initialize the database connection
	slog.Info("database connecting", "database_url", os.Getenv("DATABASE_URL"), "database_type", "sqlite")
	database := db.Connect()
	db.AutoMigrate(database)

	e := server.NewEchoServer()

	mood.InitModule(e, database)

	if os.Getenv("ENABLE_FRONTEND_DIST") == "true" {
		slog.Info("frontend dist enabled, initializing")
		// Initialize frontend module last to ensure it can serve static files.
		frontend.InitModule(e)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// Start the server in a goroutine to allow graceful shutdown
	go func() {
		slog.Info("starting server",
			"port", os.Getenv("PORT"),
			"api_url", os.Getenv("API_URL"),
		)

		if err := e.Start(":" + os.Getenv("PORT")); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				slog.Error("failed to start server", "error", err)
				os.Exit(1)
			}
		}
	}()

	<-ctx.Done()
	slog.Info("shutting down server gracefully")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		slog.Error("failed to gracefully shutdown server", "error", err)
		os.Exit(1)
	}

	slog.Info("server stopped")
}
