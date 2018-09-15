package store

import (
	"net/url"
)

func buildURL(path string, values url.Values) string {
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
