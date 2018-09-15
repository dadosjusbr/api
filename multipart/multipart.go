package multipart

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
)

// UploadRequest returns the request to upload multipart file, sending params.
func UploadRequest(uri string, payloadParam, payloadName string, payload io.Reader, params map[string]string) (*http.Request, error) {
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	part, err := w.CreateFormFile(payloadParam, payloadName)
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, payload)
	if err != nil {
		return nil, err
	}
	for key, val := range params {
		if err := w.WriteField(key, val); err != nil {
			return nil, err
		}
	}
	if err := w.Close(); err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", uri, &buf)
	req.Header.Set("Content-Type", w.FormDataContentType())
	return req, nil
}
