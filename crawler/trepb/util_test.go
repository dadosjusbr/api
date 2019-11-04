package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"
)

// Test if httpReq is making desired request and returning correct parsed document.
func Test_httpReq(t *testing.T) {
	htmlSample := "<html><head></head><body><div><span></span></div></body></html>"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, htmlSample)
	}))
	defer ts.Close()

	req, err := http.NewRequest("GET", ts.URL, nil)
	assert.NoError(t, err)

	doc, err := httpReq(req)
	assert.NoError(t, err)

	var buf bytes.Buffer
	assert.NoError(t, html.Render(&buf, doc))
	// HTML parser adds a \n before closing of body tag.
	assert.Equal(t, "<html><head></head><body><div><span></span></div>\n</body></html>", buf.String())
}

// Test if substringBetween is returning the desired strings.
func Test_substringBetween(t *testing.T) {
	type args struct {
		str    string
		before string
		after  string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"no splits", args{"hello", "g", "g"}, "hello"},
		{"both splits", args{"ghellog", "g", "o"}, "hell"},
		{"first split only", args{"ghello", "g", "m"}, "hello"},
		{"last split only", args{"hellog", "m", "g"}, "hello"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := substringBetween(tt.args.str, tt.args.before, tt.args.after); got != tt.want {
				t.Errorf("substringBetween() = %v, want %v", got, tt.want)
			}
		})
	}
}
