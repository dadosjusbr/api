package twitter

import (
	"fmt"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

// Client struct encapsulates a twitter Client object
// to allow access to twitter API.
type Client struct {
	client *twitter.Client
}

// NewClient builds some config objects to
// setup the twitter Client handler
func NewClient(consumerKey, consumerSecret, accessToken, accessSecret string) *Client {
	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessSecret)
	httpClient := config.Client(oauth1.NoContext, token)

	return &Client{
		client: twitter.NewClient(httpClient),
	}
}

// Post gets a message to be posted on twitter
// as a string and posts it in the account
// related to the access keys provided.
func (c *Client) Post(message string) error {
	_, _, err := c.client.Statuses.Update(message, nil)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
