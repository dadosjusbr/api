package main

import (
	"net/http"
	"strings"
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

// httpReq makes specified request and returns the html parsed tree.
func httpReq(req *http.Request) (*html.Node, error) {
	resp, err := netClient.Do(req)
	if err != nil {
		//return nil, fmt.Errorf("error making request to %s: %q", reqURL., err)
	}
	defer resp.Body.Close()

	doc, err := htmlquery.Parse(resp.Body)
	if err != nil {
		//return nil, fmt.Errorf("error loading doc (%s): %q", reqURL, err)
	}
	return doc, nil
}

//substringBetween returns the substring in str between before and after strings.
func substringBetween(str, before, after string) string {
	a := strings.SplitAfterN(str, before, 2)
	b := strings.SplitAfterN(a[len(a)-1], after, 2)
	if 1 == len(b) {
		return b[0]
	}
	return b[0][0 : len(b[0])-len(after)]
}
