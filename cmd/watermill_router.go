package main

import (
	"encoding/json"
	"log/slog"

	"github.com/Christian-007/fit-forge/internal/app/email/delivery/pubsub"
	emailservices "github.com/Christian-007/fit-forge/internal/app/email/services"
	"github.com/Christian-007/fit-forge/internal/app/users/dto"
	"github.com/Christian-007/fit-forge/internal/pkg/appcontext"
	"github.com/Christian-007/fit-forge/internal/pkg/security"
	"github.com/Christian-007/fit-forge/internal/pkg/topics"
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

	tokenService := security.NewTokenService(security.TokenServiceOptions{
		SecretKey: appCtx.EnvVariableService.Get("AUTH_SECRET_KEY"),
	})
	emailService := emailservices.NewEmailService(emailservices.EmailServiceOptions{
		Host:         "http://localhost:4000",
		Cache:        appCtx.RedisClient,
		TokenService: tokenService,
	})
	mailtrapSender := emailservices.NewMailtrapEmailService(emailservices.MailtrapSenderOptions{
		Host:   appCtx.EnvVariableService.Get("EMAIL_HOST"),
		ApiKey: appCtx.EnvVariableService.Get("MAILTRAP_API_KEY"),
	})

	emailPubSubHandler := pubsub.NewEmailPubSubHandler(pubsub.EmailPubSubHandlerOptions{
		EmailService:   emailService,
		MailtrapSender: mailtrapSender,
	})

	router.AddNoPublisherHandler(
		"send_email_verification",
		topics.UserRegistered,
		subscriber,
		func(msg *message.Message) error {
			var payload dto.UserResponse
			err := json.Unmarshal(msg.Payload, &payload)
			if err != nil {
				return err
			}

			appCtx.Logger.Info("Subscribing to UserRegistered",
				slog.String("UUID", msg.UUID),
				slog.Any("message", payload),
			)

			err = emailPubSubHandler.SendEmailVerification(payload)
			if err != nil {
				appCtx.Logger.Error("[send_email_verification Subscriber] Error sending an email verification",
					slog.String("error", err.Error()),
				)
				return err
			}

			appCtx.Logger.Info("Successfully send an email verification",
				slog.String("email", payload.Email),
			)
			return nil
		},
	)

	return router

}
