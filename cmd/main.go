package main

import (
	"context"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/Christian-007/fit-forge/internal/api/domains"
	"github.com/Christian-007/fit-forge/internal/api/routers"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func main() {
	// Accepting `-addr="{port}"` flag via terminal
	addr := flag.String("addr", ":4000", "HTTP network address")

	// Instantiate the all application dependencies
	appCtx := domains.AppContext{
		Logger: slog.New(slog.NewTextHandler(os.Stdout, nil)),
	}

	// Load `.env` file
	err := godotenv.Load()
	if err != nil {
		appCtx.Logger.Error(err.Error())
		os.Exit(1)
	}

	// Open DB connection
	conn, err := openDB()
	if err != nil {
		appCtx.Logger.Error(err.Error())
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	// HTTP Server configurations (Non TLS)
	server := &http.Server{
		Addr:         *addr,
		Handler:      routers.Routes(appCtx),
		ErrorLog:     slog.NewLogLogger(appCtx.Logger.Handler(), slog.LevelError),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	appCtx.Logger.Info("starting server", "addr", *addr)

	err = server.ListenAndServe()
	appCtx.Logger.Error(err.Error())
}

func openDB() (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), os.Getenv("POSTGRES_URL"))
	if err != nil {
		return nil, err
	}

	err = conn.Ping(context.Background())
	if err != nil {
		conn.Close(context.Background())
		return nil, err
	}

	return conn, nil
}
