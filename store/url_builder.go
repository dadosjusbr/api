package store

import (
	"net/url"
)

// buildPCLoudURL returns the full URL string formatted to consume the PCloud API.
func buildPCLoudURL(path string, values url.Values) string {
	const (
		apiScheme = "https"
		host      = "api.pcloud.com"
	)

	u := url.URL{
		Scheme:   apiScheme,
		Host:     host,
		Path:     path,
		RawQuery: values.Encode(),
	}

	return u.String()
}
