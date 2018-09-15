package publisher

import (
	"bytes"
	"net/http"

	"github.com/dadosjusbr/remuneracao-magistrados/multipart"
)

const (
	path          = "https://dadosjusbr-publisher.herokuapp.com/"
	paramName     = "dataset"
	paramFileName = "datapackage.zip"
)

// Publish publishes the data content in the datahub.io/dadosjusbr.
func Publish(content []byte) error {
	req, err := multipart.UploadRequest(path, paramName, paramFileName, bytes.NewReader(content), map[string]string{})
	if err != nil {
		return nil
	}
	if _, err := http.DefaultClient.Do(req); err != nil {
		return err
	}
	return nil
}
