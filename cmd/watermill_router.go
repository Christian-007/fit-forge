package main

import (
	"log/slog"

	emailpubsub "github.com/Christian-007/fit-forge/internal/app/email/delivery/pubsub"
	"github.com/Christian-007/fit-forge/internal/pkg/appcontext"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-amqp/v3/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill/message"
)

func NewWatermillRouter(amqpConfig amqp.Config, watermillLogger watermill.LoggerAdapter, appCtx appcontext.AppContext) *message.Router {
	router, err := message.NewRouter(message.RouterConfig{}, watermillLogger)
	if err != nil {
		appCtx.Logger.Error("Failed to create Watermill Router",
			slog.String("error", err.Error()),
		)
		panic(err)
	}

	subscriber, err := amqp.NewSubscriber(
		amqpConfig,
		watermillLogger,
	)
	if err != nil {
		appCtx.Logger.Error("Failed to connect to RabbitMQ",
			slog.String("error", err.Error()),
		)
		panic(err)
	}

	emailpubsub.Routes(router, subscriber, appCtx)

	return router

}
