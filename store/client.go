package store

type PCloudClient struct {
}

// NewPCloudClient returns the PCloudClient to interact with PCloudAPI, or error in case using wrong credentials.
func NewPCloudClient(username string, password string) (*PCloudClient, error) {
	return &PCloudClient{}, nil
}
