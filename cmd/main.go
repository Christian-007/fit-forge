package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/Christian-007/fit-forge/internal/api/domains"
	"github.com/Christian-007/fit-forge/internal/api/routers"
)

func main() {
	// Accepting `-addr="{port}"` flag via terminal
	addr := flag.String("addr", ":4000", "HTTP network address")

	// Instantiate the all application dependencies
	appCtx := domains.AppContext{
		Logger: slog.New(slog.NewTextHandler(os.Stdout, nil)),
	}

	server := &http.Server{
		Addr:         *addr,
		Handler:      routers.Routes(appCtx),
		ErrorLog:     slog.NewLogLogger(appCtx.Logger.Handler(), slog.LevelError),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	appCtx.Logger.Info("starting server", "addr", *addr)

	err := server.ListenAndServe()
	appCtx.Logger.Error(err.Error())
}
