package store

type PCloudClient struct {
}

func NewPCloudClient(username string, password string) (*PCloudClient, error) {
	return &PCloudClient{}, nil
}
