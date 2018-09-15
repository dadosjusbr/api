package store

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
)

// PCloudClient represents the PCloud client instance to interact with PCLoud API.
type PCloudClient struct {
	Client *http.Client
	Token  string
}

// authResponse; internal representation of the auth response. Will be used to Unmarshal the response
type authResponse struct {
	Auth string
}

// authenticate; sends the HTTP request to authenticate with PCloud using provided credentials.
func authenticate(c *http.Client, username string, password string) (string, error) {
	resp, err := c.Get(buildPCLoudURL("userinfo", url.Values{
		"getauth":  {"1"},
		"logout":   {"1"},
		"username": {username},
		"password": {password},
	}))

	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		dump, _ := httputil.DumpResponse(resp, true)

		return "", fmt.Errorf("Server responded with a non 200 (OK) status code. Response dump: \n\n%s", string(dump))
	}

	// Converting the JSON response to bytes.
	data, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}

	// We are going to use this struct to Unmarshal the JSON response from PCloud.
	jsonResponse := authResponse{}

	if err := json.Unmarshal(data, &jsonResponse); err != nil {
		return "", err
	}

	if jsonResponse.Auth == "" {
		return "", fmt.Errorf("Failed to parse the authentication. Response was: \n\n%s", string(data))
	}

	return jsonResponse.Auth, err
}

// NewPCloudClient returns the PCloudClient to interact with PCloudAPI, or error in case using wrong credentials.
func NewPCloudClient(username string, password string) (*PCloudClient, error) {
	c := &http.Client{}

	// We are hitting the PCloud API when to create the instance.
	token, err := authenticate(c, username, password)

	if err != nil {
		return nil, err
	}

	// PCloudClient the instance needs an HTTP client and the token.
	return &PCloudClient{Client: c, Token: token}, nil
}
