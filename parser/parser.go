package parser

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

const url = "https://dadosjusbr-parser.herokuapp.com/"
const schemaResource = "schema"

// ServiceClient parses XLS/XLSX files by delegating requests to the parser service.
type ServiceClient struct {
	url    string
	client *http.Client
}

// NewServiceClient returns a parser which delegates the parsing process to
// the dadosjusbr parser microservice.
func NewServiceClient(url string) *ServiceClient {
	c := &http.Client{
		Timeout: 5 * time.Minute,
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout: 1 * time.Minute, // Important because the microsservice could be sleeping.
			}).Dial,
		},
	}
	return &ServiceClient{url, c}
}

// Parse parse the spreadsheet contents and returns an unified parsed CSV and its schema.
func (s *ServiceClient) Parse(contents [][]byte, names []string) ([]byte, map[string]interface{}, error) {
	if len(contents) != len(names) {
		return nil, nil, fmt.Errorf("error Parser: contents (%d) and names (%d) must be the same size. ", len(contents), len(names))
	}
	sch, err := s.getSchema()
	if err != nil {
		return nil, nil, err
	}
	var result bytes.Buffer
	for i, c := range contents {
		parseST := time.Now()
		filename := names[i]
		if i == 0 {
			data, err := s.request(s.url, c)
			if err != nil {
				return nil, sch, err
			}
			result.Write(data)
			result.WriteRune('\n')
			continue
		}
		data, err := s.request(fmt.Sprint(s.url, "?headless=true"), c)
		if err != nil {
			return nil, sch, err
		}
		result.Write(data)
		if i < len(contents)-1 {
			result.WriteRune('\n')
		}
		fmt.Printf("File (%s) successfuly parsed. Took %v\n", filename, time.Now().Sub(parseST))
	}
	return result.Bytes(), sch, nil
}

func (s *ServiceClient) request(url string, body []byte) ([]byte, error) {
	resp, err := s.client.Post(url, "application/octet-stream", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("error sending request do parser service(%s):%q", url, err)
	}

	fmt.Println("===================")
	fmt.Println("===================")
	fmt.Println(resp.StatusCode)
	fmt.Println("===================")
	fmt.Println("===================")

	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading parser service response(%s):%q", url, err)
	}
	return data, nil
}

func (s *ServiceClient) getSchema() (map[string]interface{}, error) {
	resp, err := s.client.Get(url + schemaResource)
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
