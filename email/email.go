package email

import (
	//"fmt"
	"github.com/sendgrid/sendgrid-go"
)

type Client struct {
	username string
	password string
}

func NewClient(username, password string) (*Client, error) {
	if len(username) <= 0 {
		return nil, error.New("Username should not be empty")
	}
	if len(password) <= 0 {
		return nil, error.New("Password should not be empty")
	}
	return &Client{username, password }, nil
}

//func Send(to, subject, body string) error {
//
//}