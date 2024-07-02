package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/Christian-007/fit-forge/internal/db"
	"github.com/Christian-007/fit-forge/internal/pkg/appcontext"

	"github.com/joho/godotenv"
)

func main() {
	// Accepting `-addr="{port}"` flag via terminal
	addr := flag.String("addr", ":4000", "HTTP network address")

	// Initialize logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Load `.env` file
	err := godotenv.Load()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// Open DB connection
	pool, err := db.OpenPostgresDbPool(os.Getenv("POSTGRES_URL"))
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer pool.Close()

	// Instantiate the all application dependencies
	appCtx := appcontext.NewAppContext(appcontext.AppContextOptions{
		Logger: logger,
		Pool:   pool,
	})

	// HTTP Server configurations (Non TLS)
	server := &http.Server{
		Addr:         *addr,
		Handler:      Routes(appCtx),
		ErrorLog:     slog.NewLogLogger(appCtx.Logger.Handler(), slog.LevelError),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	logger.Info("starting server", "addr", *addr)

	err = server.ListenAndServe()
	logger.Error(err.Error())
}
