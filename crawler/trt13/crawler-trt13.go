package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

var netClient = &http.Client{
	Timeout: time.Second * 60,
}

func main() {
	month := flag.Int("mes", 0, "Mês a ser analisado")
	year := flag.Int("ano", 0, "Ano a ser analisado")
	flag.Parse()
	if *month == 0 || *year == 0 {
		log.Fatalf("Need all arguments to continue, please try again: \"go run crawler-trtpb.go --mes=int --ano=int\"")
	}
	fileName := fmt.Sprintf("remuneracoes-trt13-%02d-%04d.json", *month, *year)
	f, err := os.Create(fileName)
	if err != nil {
		log.Fatalf("Error creating file(%s):%q", fileName, err)
	}
	defer f.Close()

	reqURL := fmt.Sprintf("https://www.trt13.jus.br/transparenciars/api/anexoviii/anexoviii?mes=%02d&ano=%04d", *month, *year)
	if err = download(reqURL, f); err != nil {
		os.Remove(fileName)
		log.Fatalf("Error while downloading content (%02d-%04d): %q", *month, *year, err)
	}
}

func download(reqURL string, w io.Writer) error {
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return fmt.Errorf("error while creating *http.Request: %q", err)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := netClient.Do(req)
	if err != nil {
		return fmt.Errorf("error while making GET request to (%s): %q", reqURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code. Request: GET (%s) - Response: (%d): %s", reqURL, resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	if io.Copy(w, resp.Body); err != nil {
		return fmt.Errorf("error copying response content:%q", err)
	}
	return nil
}
