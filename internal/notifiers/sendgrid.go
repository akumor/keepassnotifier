package notifiers

import (
	"errors"
	"fmt"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type SendGrid struct {
	ApiKey           string
	FromEmailAddress string
	ToEmailAddress   string
}

// NewSendGrid returns a new SendGrid struct
func NewSendGrid(apiKey, fromEmail, toEmail string) (*SendGrid, error) {
	if apiKey == "" {
		return nil, errors.New("empty string provided for sendgrid api key")
	}

	if fromEmail == "" {
		return nil, errors.New("empty string provided for sendgrid from email")
	}

	if toEmail == "" {
		return nil, errors.New("empty string provided for sendgrid to email")
	}

	return &SendGrid{
		ApiKey:           apiKey,
		FromEmailAddress: fromEmail,
		ToEmailAddress:   toEmail,
	}, nil
}

func (s *SendGrid) Notify(htmlContent string) error {
	from := mail.NewEmail("Keepass Notifier", s.FromEmailAddress)
	subject := "New Notifications From Keepass"
	to := mail.NewEmail("Keepass User", s.ToEmailAddress)
	plainTextContent := ""
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(s.ApiKey)
	response, err := client.Send(message)
	if err != nil {
		return err
	}
	// check StatusCode for success
	if response.StatusCode != 202 {
		return fmt.Errorf("sendgrid api returned status code: %d\nsendgrid error body: %s", response.StatusCode, response.Body)
	}
	return nil
}
