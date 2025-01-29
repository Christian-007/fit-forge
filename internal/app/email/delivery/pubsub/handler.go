package pubsub

import (
	"github.com/Christian-007/fit-forge/internal/app/email/domains"
	"github.com/Christian-007/fit-forge/internal/app/email/services"
	"github.com/Christian-007/fit-forge/internal/app/users/dto"
)

type EmailPubSubHandler struct {
	EmailPubSubHandlerOptions
}

type EmailPubSubHandlerOptions struct {
	EmailService   services.EmailService
	MailtrapSender services.MailtrapSender
}

func NewEmailPubSubHandler(options EmailPubSubHandlerOptions) EmailPubSubHandler {
	return EmailPubSubHandler{
		options,
	}
}

func (e EmailPubSubHandler) SendEmailVerification(userResponse dto.UserResponse) error {
	verificationLink, err := e.EmailService.CreateVerificationLink(userResponse.Email)
	if err != nil {
		return err
	}

	emailRequest := domains.EmailWithTemplateRequest{
		From: domains.EmailAddressOptions{
			Email: "hello@demomailtrap.com",
			Name:  "No Reply at Fit Forge",
		},
		To: []domains.EmailAddressOptions{
			{
				Email: userResponse.Email,
				Name:  userResponse.Name,
			},
		},
		TemplateUuid: "fdbefad8-2410-45d2-bded-9d1b647ac416",
		TemplateVariables: map[string]any{
			"user_name":         userResponse.Name,
			"verification_link": verificationLink,
		},
	}
	err = e.MailtrapSender.SendWithTemplate(emailRequest)
	if err != nil {
		return err
	}

	return nil
}
