package email

import (
	"errors"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"net/http"
)

type Client struct {
	apiKey string
}

func NewClient(apiKey string) (*Client, error) {
	if len(apiKey) <= 0 {
		return nil, errors.New("Api Key should not be empty")
	}
	return &Client{apiKey }, nil
}

func (c *Client) Send(from, to, subject, body string) error {
	fromMail := mail.NewEmail("", from)
	toMail := mail.NewEmail("", to)
	message := mail.NewSingleEmail(fromMail, subject, toMail, body, body)
	client := sendgrid.NewSendClient(c.apiKey)
	response, err := client.Send(message)
	if err != nil {
		return err
	}
	if response.StatusCode != http.StatusAccepted {
		return errors.New(response.Body)
	}
	return nil
}