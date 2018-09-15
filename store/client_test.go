package store

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestBuildsPCloudURL(t *testing.T) {
	url := buildPCLoudURL("testing", url.Values{"lorem": {"ipsum"}})

	if url != "https://api.pcloud.com/testing?lorem=ipsum" {
		t.Error("Could not properly build the URL to pcloud")
	}
}

func TestAuthenticateWithPCloud(t *testing.T) {
	username, password := "fakeuser", "fakepass"

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test request parameters
		if req.URL.String() != "https://api.pcloud.com/userinfo?getauth=1&logout=1&username="+username+"&password="+password {
			t.Error("Authentication is using the wrong URL.")
		}
		// Send response to be tested
		rw.Write([]byte(`{"auth": "fakeToken"}`))
	}))

	resp, _ := authenticate(server.Client(), username, password)

	if resp.StatusCode != 200 {
		t.Error("Authentication failed. Probably using the wrong URL.")
	}
}
