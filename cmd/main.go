package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/Christian-007/fit-forge/internal/db"
	"github.com/Christian-007/fit-forge/internal/pkg/appcontext"
	"github.com/Christian-007/fit-forge/internal/pkg/applog"
	"github.com/Christian-007/fit-forge/internal/pkg/cache"
	"github.com/Christian-007/fit-forge/internal/pkg/envvariable"

	"github.com/redis/go-redis/v9"
)

func main() {
	// Accepting `-addr="{port}"` flag via terminal
	addr := flag.String("addr", ":4000", "HTTP network address")

	// Initialize logger
	slogLogger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	logger := applog.NewSlogLogger(slogLogger)

	// Load `.env` file
	envVariableService := envvariable.GodotEnvVariableService{}
	err := envVariableService.Load()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// Open DB connection
	pool, err := db.OpenPostgresDbPool(envVariableService.Get("POSTGRES_URL"))
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer pool.Close()

	// Open Redis Connection
	client, err := cache.NewRedisCache(&redis.Options{
		Addr:     envVariableService.Get("REDIS_DSN"),
		Password: "",
		DB:       0,
	})
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// Instantiate the all application dependencies
	appCtx := appcontext.NewAppContext(appcontext.AppContextOptions{
		Logger:             logger,
		Pool:               pool,
		RedisClient:        client,
		EnvVariableService: envVariableService,
	})

	// HTTP Server configurations (Non TLS)
	server := &http.Server{
		Addr:         *addr,
		Handler:      Routes(appCtx),
		ErrorLog:     logger.StandardLogger(applog.LevelError),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	logger.Info("starting server", "addr", *addr)

	err = server.ListenAndServe()
	logger.Error(err.Error())
}
