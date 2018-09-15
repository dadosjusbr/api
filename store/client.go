package store

type Client struct {
}

func NewClient(username string, password string) (*Client, error) {
	return &Client{}, nil
}