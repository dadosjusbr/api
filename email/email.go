package email

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"text/template"

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

//SendFailMail sends an email with the standart template for showing an error
func (c *Client) SendFailMail(from, to string, month, year int, err error) error {
	successTemplate := `
		<!DOCTYPE html>
		<html>
			<head>
				<meta charset="UTF-8">
			</head>
			<body>
				The pipeline for {{.Month}}/{{.Year}} failed with the following error: 
				<br>
				{{.Error}}
			</body>
		</html>
	`

	data := struct {
		Month int
		Year  int
		Error string
	}{month, year, strings.Replace(err.Error(), "\\n", "<br>", -1)}

	tmpl, err := template.New("email").Parse(successTemplate)
	if err != nil {
		return err
	}
	var tpl bytes.Buffer
	if err := tmpl.Execute(&tpl, data); err != nil {
		return err
	}
	body := tpl.String()
	c.Send(from, to, fmt.Sprintf("remuneracao-magistrados: The pipeline for %d/%d has failed", month, year), body)
	return nil
}

//SendSuccessMail sends an email notifying the success of the month processing
func (c *Client) SendSuccessMail(from, to string, month, year int) error {
	successTemplate := `
		<!DOCTYPE html>
		<html>
			<head>
				<meta charset="UTF-8">
			</head>
			<body>
				The pipeline for {{.Month}}/{{.Year}} was able to publish successfuly the data!! \o/ 
			</body>
		</html>
	`

	data := struct {
		Month int
		Year  int
	}{month, year}

	tmpl, err := template.New("email").Parse(successTemplate)
	if err != nil {
		return err
	}
	var tpl bytes.Buffer
	if err := tmpl.Execute(&tpl, data); err != nil {
		return err
	}
	body := tpl.String()
	c.Send(from, to, fmt.Sprintf("remuneracao-magistrados: The data from %d/%d was successfuly published", month, year), body)
	return nil
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
