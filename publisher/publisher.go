package publisher

import (
	"bytes"
	"fmt"
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
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Error publishing datapackage: Invalid status code %d (%s)", resp.StatusCode, resp.Status)
	}
	return nil
}
