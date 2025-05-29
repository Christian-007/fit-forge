package main

import (
	"context"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Christian-007/fit-forge/internal/db"
	"github.com/Christian-007/fit-forge/internal/pkg/appcontext"
	"github.com/Christian-007/fit-forge/internal/pkg/applog"
	"github.com/Christian-007/fit-forge/internal/pkg/cache"
	"github.com/Christian-007/fit-forge/internal/pkg/decorator"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-googlecloud/pkg/googlecloud"

	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/errgroup"
)

func main() {
	// Accepting `-addr="{port}"` flag via terminal
	addr := flag.String("addr", ":4000", "HTTP network address")

	// Initialize logger
	slogLogger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	logger := applog.NewSlogLogger(slogLogger)

	// Open DB connection
	pool, err := db.OpenPostgresDbPool(os.Getenv("POSTGRES_URL"))
	if err != nil {
		logger.Error("Failed to connect to Postgresql",
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}
	defer pool.Close()

	// Open Redis Connection
	client, err := cache.NewRedisCache(&redis.Options{
		Addr:     os.Getenv("REDIS_DSN"),
		Username: os.Getenv("REDIS_USERNAME"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})
	if err != nil {
		logger.Error("Failed to connect to Redis",
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}

	// Create a connection to a Message Broker
	watermillLogger := watermill.NewStdLogger(false, false)
	gcloudPublisher, err := googlecloud.NewPublisher(
		googlecloud.PublisherConfig{
			ProjectID: os.Getenv("PUBSUB_PROJECT_ID"),
		},
		watermillLogger,
	)
	if err != nil {
		logger.Error("Failed to create a publisher in Google Pub/Sub",
			slog.String("error", err.Error()),
		)
		panic(err)
	}
	publisher := decorator.PublishWithCorrelationId{Publisher: gcloudPublisher}
	defer publisher.Close()

	// Instantiate the all application dependencies
	appCtx := appcontext.NewAppContext(appcontext.AppContextOptions{
		Logger:      logger,
		Pool:        pool,
		RedisClient: client,
		Publisher:   publisher,
	})

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	errgrp, ctx := errgroup.WithContext(ctx)

	// Instantiate the PubSub router
	watermillRouter := NewWatermillRouter(watermillLogger, appCtx)
	errgrp.Go(func() error {
		logger.Info("starting PubSub router...")
		err = watermillRouter.Run(ctx) // Starting the PubSub router in a Goroutine
		if err != nil {
			logger.Error("Failed to start PubSub router",
				slog.String("error", err.Error()),
			)

			return err
		}
		return nil
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

	errgrp.Go(func() error {
		// We don't want to start the HTTP server before Watermill router (so service won't be healthy before it's ready)
		<-watermillRouter.Running()

		logger.Info("starting server", "addr", *addr)

		err = server.ListenAndServe()
		if err != nil {
			logger.Error(err.Error())
			return err
		}

		return nil
	})

	// Start gRPC Server
	errgrp.Go(func() error {
		addr := ":50051"
		logger.Info("starting gRPC server", "addr", addr)

		grpcServicesFn := InitGrpcServices(appCtx)
		err = StartGrpcServer(ctx, addr, grpcServicesFn)
		if err != nil {
			logger.Error(err.Error())
			return err
		}

		return nil
	})

	errgrp.Go(func() error {
		<-ctx.Done()
		return server.Shutdown(ctx)
	})

	err = errgrp.Wait()
	if err != nil {
		panic(err)
	}
}
