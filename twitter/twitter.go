package twitter

import (
	"fmt"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

// Client ...
type Client struct {
	ConsumerKey    string
	ConsumerSecret string
	AccessToken    string
	AccessSecret   string
	Client         *twitter.Client
}

// NewClient ...
func NewClient(consumerKey, consumerSecret, accessToken, accessSecret string) *Client {
	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessSecret)
	httpClient := config.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)

	return &Client{
		ConsumerKey:    consumerKey,
		ConsumerSecret: consumerSecret,
		AccessToken:    accessToken,
		AccessSecret:   accessSecret,
		Client:         client,
	}
}

// Post ...
func (client *Client) Post(message string) error {
	_, _, err := client.Client.Statuses.Update(message, nil)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
