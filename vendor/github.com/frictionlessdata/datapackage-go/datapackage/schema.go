package datapackage

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/frictionlessdata/tableschema-go/schema"
)

func loadSchema(p string) (map[string]interface{}, error) {
	var reader io.Reader
	if strings.HasPrefix(p, "http") {
		resp, err := http.Get(p)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		reader = resp.Body
	} else {
		f, err := os.Open(p)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		reader = f
	}
	buf, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	_, err = schema.Read(bytes.NewBuffer(buf))
	if err != nil {
		return nil, err
	}
	var ret map[string]interface{}
	if err := json.Unmarshal(buf, &ret); err != nil {
		return nil, err
	}
	return ret, nil
}
