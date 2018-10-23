package email

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// Client Class containing the needed methods to send an email
type Client struct {
	sendgridClient *sendgrid.Client
}

// NewClient class constructor function
func NewClient(apiKey string) (*Client, error) {
	if len(apiKey) <= 0 {
		return nil, errors.New("Api Key should not be empty")
	}

	return &Client{sendgrid.NewSendClient(apiKey)}, nil
}

// Send an email and return an error if it fails
func (c *Client) Send(from, to, subject, body string) error {
	fromMail := mail.NewEmail("", from)
	toMail := mail.NewEmail("", to)
	message := mail.NewSingleEmail(fromMail, subject, toMail, body, body)
	response, err := c.sendgridClient.Send(message)

	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusAccepted {
		return fmt.Errorf("An error has ocurred while trying to send an email.\n"+
			"Status Code: %d\n"+
			"Body: %s\n"+
			"Headers: %v",
			response.StatusCode,
			response.Body,
			response.Headers)
	}

	return nil
}
