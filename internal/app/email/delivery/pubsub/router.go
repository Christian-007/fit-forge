package pubsub

import (
	"encoding/json"
	"log/slog"

	emailservices "github.com/Christian-007/fit-forge/internal/app/email/services"
	"github.com/Christian-007/fit-forge/internal/app/users/dto"
	"github.com/Christian-007/fit-forge/internal/pkg/appcontext"
	"github.com/Christian-007/fit-forge/internal/pkg/security"
	"github.com/Christian-007/fit-forge/internal/pkg/topics"
	"github.com/ThreeDotsLabs/watermill-googlecloud/pkg/googlecloud"
	"github.com/ThreeDotsLabs/watermill/message"
)

func Routes(router *message.Router, subscriber *googlecloud.Subscriber, appCtx appcontext.AppContext) {
	// Instantiate dependencies
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

	// Instantiate handler
	emailPubSubHandler := NewEmailPubSubHandler(EmailPubSubHandlerOptions{
		EmailService:   emailService,
		MailtrapSender: mailtrapSender,
	})

	// Add handler into router
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
}
