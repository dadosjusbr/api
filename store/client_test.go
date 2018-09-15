package store

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

// RoundTripFunc .
type RoundTripFunc func(req *http.Request) *http.Response

// RoundTrip .
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

// NewTestClient returns *http.Client with Transport replaced to avoid making real calls
func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}

func TestBuildsPCloudURL(t *testing.T) {
	url := buildPCLoudURL("testing", url.Values{"lorem": {"ipsum"}})

	if url != "https://api.pcloud.com/testing?lorem=ipsum" {
		t.Error("Could not properly build the URL to pcloud")
	}
}

func TestAuthenticateWithPCloud(t *testing.T) {
	username, password := "fakeuser", "fakepass"

	client := NewTestClient(func(req *http.Request) *http.Response {
		// Test request parameters
		authUrl := buildPCLoudURL("userinfo", url.Values{
			"getauth":  {"1"},
			"logout":   {"1"},
			"username": {username},
			"password": {password},
		})

		if req.URL.String() != authUrl {
			t.Error("Authentication is using the wrong URL.")
		}

		return &http.Response{
			StatusCode: http.StatusOK,
			// Send response to be tested
			Body: ioutil.NopCloser(bytes.NewBufferString(`{"auth": "fakeToken"}`)),
			// Must be set to non-nil value or it panics
			Header: make(http.Header),
		}
	})

	resp, err := authenticate(client, username, password)

	if err != nil {
		t.Error("Failed to build the authenticationResponse")
	}

	if resp == nil {
		t.Error("Expected authentication response to not be nil")
	}

	if resp.Auth != "fakeToken" {
		t.Error("Failed to get the token from the JSON response")
	}
}

func TestAuthenticateHandlesWrongCredentialsResponse(t *testing.T) {
	username, password := "wronguser", "wrongpass"

	client := NewTestClient(func(req *http.Request) *http.Response {
		// Test request parameters
		authUrl := buildPCLoudURL("userinfo", url.Values{
			"getauth":  {"1"},
			"logout":   {"1"},
			"username": {username},
			"password": {password},
		})

		if req.URL.String() != authUrl {
			t.Error("Authentication is using the wrong URL.")
		}

		return &http.Response{
			StatusCode: http.StatusOK,
			// Send response to be tested
			Body: ioutil.NopCloser(bytes.NewBufferString(`{"error": "Failed to login."}`)),
			// Must be set to non-nil value or it panics
			Header: make(http.Header),
		}
	})

	pcloud, err := authenticate(client, username, password)

	if err == nil {
		t.Error("Failed to build the authenticationResponse")
	}

	if pcloud != nil {
		t.Error("Expected the authResponse to be nil")
	}

	if err.Error() != "Failed to parse the authentication. Response was: \n\n{\"error\": \"Failed to login.\"}" {
		t.Error("Error response should contains the payload.")
	}
}

func TestAuthenticationFailsWhenResponseIsNot200OK(t *testing.T) {
	username, password := "wronguser", "wrongpass"

	client := NewTestClient(func(req *http.Request) *http.Response {
		// Test request parameters
		authUrl := buildPCLoudURL("userinfo", url.Values{
			"getauth":  {"1"},
			"logout":   {"1"},
			"username": {username},
			"password": {password},
		})

		if req.URL.String() != authUrl {
			t.Error("Authentication is using the wrong URL.")
		}

		return &http.Response{
			StatusCode: http.StatusInternalServerError,
			// Send response to be tested
			Body: ioutil.NopCloser(bytes.NewBufferString(``)),
			// Must be set to non-nil value or it panics
			Header: make(http.Header),
		}
	})

	pcloud, err := authenticate(client, username, password)

	if err == nil {
		t.Error("Failed to build the authenticationResponse")
	}

	if pcloud != nil {
		t.Error("Expected the authResponse to be nil")
	}

	if !strings.Contains(err.Error(), "Server responded with a non 200 (OK) status code. Response dump: ") {
		t.Error("Failed to dump the response")
	}

	if !strings.Contains(err.Error(), "HTTP/0.0 500 Internal Server Error") {
		t.Error("Failed to dump the response")
	}
}
