package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

var netClient = &http.Client{
	Timeout: time.Second * 60,
	CheckRedirect: func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	},
}

//loadURL loads HTML document from specified URL.
func loadURL(baseURL string) (*html.Node, error) {
	resp, err := netClient.Get(baseURL)
	if err != nil {
		return nil, fmt.Errorf("error making GET request to %s: %q", baseURL, err)
	}
	defer resp.Body.Close()

	doc, err := htmlquery.Parse(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error loading doc (%s): %q", baseURL, err)
	}
	return doc, nil
}
