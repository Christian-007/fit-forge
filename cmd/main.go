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
	"github.com/Christian-007/fit-forge/internal/pkg/envvariable"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-amqp/v3/pkg/amqp"

	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/errgroup"
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
		logger.Error("Failed to connect to Postgresql",
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}
	defer pool.Close()

	// Open Redis Connection
	client, err := cache.NewRedisCache(&redis.Options{
		Addr:     envVariableService.Get("REDIS_DSN"),
		Password: envVariableService.Get("REDIS_PASSWORD"),
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
	amqpConfig := amqp.NewDurableQueueConfig(envVariableService.Get("RABBITMQ_URL"))
	amqpPublisher, err := amqp.NewPublisher(amqpConfig, watermillLogger)
	publisher := decorator.PublishWithCorrelationId{Publisher: amqpPublisher}

	if err != nil {
		logger.Error("Failed to create a publisher in RabbitMQ",
			slog.String("error", err.Error()),
		)
		panic(err)
	}
	defer publisher.Close()

	// Instantiate the all application dependencies
	appCtx := appcontext.NewAppContext(appcontext.AppContextOptions{
		Logger:             logger,
		Pool:               pool,
		RedisClient:        client,
		EnvVariableService: envVariableService,
		Publisher:          publisher,
	})

	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	errgrp, ctx := errgroup.WithContext(ctx)

	// Instantiate the PubSub router
	watermillRouter := NewWatermillRouter(amqpConfig, watermillLogger, appCtx)
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
		err = StartGrpcServer(addr, grpcServicesFn)
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
