package store

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
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
// based on: http://hassansin.github.io/Unit-Testing-http-client-in-Go
func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}

func TestBuildsPCloudURL(t *testing.T) {
	URL := buildPCLoudURL("testing", url.Values{"lorem": {"ipsum"}})

	if URL != "https://api.pcloud.com/testing?lorem=ipsum" {
		t.Error("Could not properly build the URL to pcloud")
	}
}

func TestAuthenticateWithPCloud(t *testing.T) {
	username, password := "fakeuser", "fakepass"

	client := NewTestClient(func(req *http.Request) *http.Response {
		// Test request parameters
		authURL := buildPCLoudURL("userinfo", url.Values{
			"getauth":  {"1"},
			"logout":   {"1"},
			"username": {username},
			"password": {password},
		})

		if req.URL.String() != authURL {
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

	token, err := authenticate(client, username, password)

	if err != nil {
		t.Error("Failed to build the authenticationResponse")
	}

	if token == "" {
		t.Error("Expected authentication response to not be nil")
	}

	if token != "fakeToken" {
		t.Error("Failed to get the token from the JSON response")
	}
}

func TestAuthenticateHandlesWrongCredentialsResponse(t *testing.T) {
	username, password := "wronguser", "wrongpass"

	client := NewTestClient(func(req *http.Request) *http.Response {
		// Test request parameters
		authURL := buildPCLoudURL("userinfo", url.Values{
			"getauth":  {"1"},
			"logout":   {"1"},
			"username": {username},
			"password": {password},
		})

		if req.URL.String() != authURL {
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

	token, err := authenticate(client, username, password)

	if err == nil {
		t.Error("Failed to build the authenticationResponse")
	}

	if token != "" {
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
		authURL := buildPCLoudURL("userinfo", url.Values{
			"getauth":  {"1"},
			"logout":   {"1"},
			"username": {username},
			"password": {password},
		})

		if req.URL.String() != authURL {
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

	token, err := authenticate(client, username, password)

	if err == nil {
		t.Error("Failed to build the authenticationResponse")
	}

	if token != "" {
		t.Error("Expected the authResponse to be nil")
	}

	if !strings.Contains(err.Error(), "Server responded with a non 200 (OK) status code. Response dump: ") {
		t.Error("Failed to dump the response")
	}

	if !strings.Contains(err.Error(), "HTTP/0.0 500 Internal Server Error") {
		t.Error("Failed to dump the response")
	}
}

func successUploadFileResponse(fileID int) *http.Response {
	return &http.Response{
		StatusCode: http.StatusOK,
		// Send response to be tested
		Body: ioutil.NopCloser(bytes.NewBufferString(fmt.Sprintf("{\"fileids\": [%d]}", fileID))),
		// Must be set to non-nil value or it panics
		Header: make(http.Header),
	}
}

func successUploadedFileLinkResponse() *http.Response {
	return &http.Response{
		StatusCode: http.StatusOK,
		// Send response to be tested
		Body: ioutil.NopCloser(bytes.NewBufferString(`{"link": "https://my.pcloud.com/#page=publink&code=LinkCode"}`)),
		// Must be set to non-nil value or it panics
		Header: make(http.Header),
	}
}

func TestCreatesFileOnPCloud(t *testing.T) {
	fakeToken := "fakeToken"
	filename := "loremipsum.txt"
	fileID := 123

	client := NewTestClient(func(req *http.Request) *http.Response {
		// Test request parameters
		uploadURL := buildPCLoudURL("uploadfile", url.Values{
			"auth":     {fakeToken},
			"path":     {"/"},
			"filename": {filename},
		})

		if req.URL.String() == uploadURL {
			return successUploadFileResponse(fileID)
		}

		generateLinkURL := buildPCLoudURL("getfilepublink", url.Values{
			"auth":   {fakeToken},
			"fileid": {strconv.Itoa(fileID)},
		})

		if req.URL.String() == generateLinkURL {
			return successUploadedFileLinkResponse()
		}

		t.Error("Unkown mocked request.")

		return nil
	})

	pcloud := PCloudClient{Client: client, Token: fakeToken}

	r := strings.NewReader("Works")

	URL, _ := pcloud.Put(filename, r)

	if URL != "https://my.pcloud.com/#page=publink&code=LinkCode" {
		t.Error("Expected mocked public link doesnt exist")
	}
}

func TestHandlesErrorWhenUploadFileFails(t *testing.T) {
	fakeToken := "fakeToken"
	filename := "loremipsum.txt"

	client := NewTestClient(func(req *http.Request) *http.Response {
		// Test request parameters
		uploadURL := buildPCLoudURL("uploadfile", url.Values{
			"auth":     {fakeToken},
			"path":     {"/"},
			"filename": {filename},
		})

		if req.URL.String() != uploadURL {
			t.Error("Put is using the wrong URL.")
		}

		return &http.Response{
			StatusCode: http.StatusInternalServerError,
			// Send response to be tested
			Body: ioutil.NopCloser(bytes.NewBufferString(`{
				"error": "Something went wrong."
			}`)),
			// Must be set to non-nil value or it panics
			Header: make(http.Header),
		}
	})

	pcloud := PCloudClient{Client: client, Token: fakeToken}

	r := strings.NewReader("Works")

	_, err := pcloud.Put(filename, r)

	if err == nil {
		t.Error("Expected error to be nil")
	}
}

func TestHandlesPublicLinkGenerationFails(t *testing.T) {
	fakeToken := "fakeToken"
	filename := "loremipsum.txt"
	fileID := 123

	client := NewTestClient(func(req *http.Request) *http.Response {
		// Test request parameters
		uploadURL := buildPCLoudURL("uploadfile", url.Values{
			"auth":     {fakeToken},
			"path":     {"/"},
			"filename": {filename},
		})

		if req.URL.String() == uploadURL {
			return successUploadFileResponse(fileID)
		}

		generateLinkURL := buildPCLoudURL("getfilepublink", url.Values{
			"auth":   {fakeToken},
			"fileid": {strconv.Itoa(fileID)},
		})

		if req.URL.String() == generateLinkURL {
			return &http.Response{
				StatusCode: http.StatusInternalServerError,
				// Send response to be tested
				Body: ioutil.NopCloser(bytes.NewBufferString(`{
					"error": "Something went wrong."
				}`)),
				// Must be set to non-nil value or it panics
				Header: make(http.Header),
			}
		}

		t.Error("Unkown mocked request.")

		return nil
	})

	pcloud := PCloudClient{Client: client, Token: fakeToken}

	r := strings.NewReader("Works")

	_, err := pcloud.Put(filename, r)

	if err == nil {
		t.Error("Expected error to be nil")
	}
}
