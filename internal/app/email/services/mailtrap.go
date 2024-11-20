package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Christian-007/fit-forge/internal/app/email/domains"
)

type MailtrapSender struct {
	MailtrapSenderOptions
}

type MailtrapSenderOptions struct {
	Host   string
	ApiKey string
}

func NewMailtrapEmailService(options MailtrapSenderOptions) MailtrapSender {
	return MailtrapSender{
		options,
	}
}

func (m MailtrapSender) SendWithTemplate(reqBody domains.EmailWithTemplateRequest) error {
	apiUrl := m.Host + "/api/send/3274815"

	reqBodyJson, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", apiUrl, bytes.NewBuffer(reqBodyJson))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+m.ApiKey)
	req.Header.Set("Content-Type", "application/json")

	httpClient := &http.Client{}
	response, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusAccepted {
		fmt.Printf("failed to send email, status code: %d", response.StatusCode)
		return nil
	}

	fmt.Println("Email sent successfully")
	return nil
}
