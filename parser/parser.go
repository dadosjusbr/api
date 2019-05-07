package parser

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/dadosjusbr/remuneracao-magistrados/multipart"
)

const url = "https://dadosjusbr-parser.herokuapp.com/"
const schemaResource = "schema"

// Parse parses the XLS(X) passed as parameters and returns the CSV contents, the request errors and other errors.
func Parse(path, fileNameParam string, params map[string]string) ([]byte, []interface{}, error) {
	content, err := zipFile(path)
	if err != nil {
		return nil, nil, fmt.Errorf("Error zipping spreadsheet (%s) err:%q", path, err)
	}
	req, err := multipart.UploadRequest(url, fileNameParam, "planilha.zip", bytes.NewReader(content), params)
	if err != nil {
		return nil, nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	fZip, err := os.Create("p.zip")
	if err != nil {
		return nil, nil, err
	}
	_, err = io.Copy(fZip, resp.Body)
	if err != nil {
		return nil, nil, err
	}

	r, err := zip.OpenReader(fZip.Name())
	if err != nil {
		return nil, nil, fmt.Errorf("Error opening zip file (%s): %q", fZip.Name(), err)
	}
	defer r.Close()
	var data bytes.Buffer
	var errors []interface{}
	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return nil, nil, err
		}
		switch f.Name {
		case "data.csv":
			_, err = io.Copy(&data, rc)
			if err != nil {
				return nil, nil, err
			}
		case "errors.txt":
			b, err := ioutil.ReadAll(rc)
			if err != nil {
				return nil, nil, err
			}
			if err := json.Unmarshal(b, &errors); err != nil {
				return nil, nil, err
			}
		}
		rc.Close()
	}
	return data.Bytes(), errors, nil
}

//GetSchema request schema from the Parser service
func GetSchema() (map[string]interface{}, error) {
	resp, err := http.Get(url + schemaResource)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	d := make(map[string]interface{})
	if err := json.Unmarshal(body, &d); err != nil {
		return nil, fmt.Errorf("Error trying to unmarshal the schema: %q", err)
	}
	return d, nil
}

func zipFile(path string) ([]byte, error) {
	fPath, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer fPath.Close()
	content, err := ioutil.ReadAll(fPath)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	w := zip.NewWriter(&buf)
	fZip, err := w.Create(filepath.Base(fPath.Name()))
	if err != nil {
		return nil, err
	}
	_, err = fZip.Write(content)
	if err != nil {
		return nil, err
	}
	if err := w.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
