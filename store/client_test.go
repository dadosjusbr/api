package store

import (
	"net/url"
	"testing"
)

func TestBuildsURL(t *testing.T) {
	url := buildURL("testing", url.Values{"lorem": {"ipsum"}})

	if url != "https://api.pcloud.com/testing?lorem=ipsum" {
		t.Error("Could not properly build the URL to pcloud")
	}
}
