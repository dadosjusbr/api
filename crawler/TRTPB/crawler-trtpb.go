package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

var statusCode = map[int]string{
	200: "OK",
	400: "Bad Request",
	404: "Records Not found!",
	500: "Internal Server Error",
}

const fileExtension = ".json"

func main() {
	month := flag.Int("mes", 0, "Mês a ser analisado")
	year := flag.Int("ano", 0, "Ano a ser analisado")
	flag.Parse()
	if *month == 0 || *year == 0 {
		log.Fatalf("need arguments to continue, please try again: \"go run crawler-trtpb.go --mes=int --ano=int\"")
	}
	fileName := fmt.Sprintf("remuneracoes-trt13-%02d-%04d%s", *month, *year, fileExtension)
	f, err := os.Create(fileName)
	if err != nil {
		log.Fatalf("error creating file(%s):%q", fileName, err)
	}

	reqURL := fmt.Sprintf("https://www.trt13.jus.br/transparenciars/api/anexoviii/anexoviii?mes=%02d&ano=%04d", *month, *year)
	err = download(reqURL, f)
	if err != nil {
		os.Remove(fileName)
		log.Fatalf("error while downloading content (%02d-%04d): %q", *month, *year, err)
	}
}

func download(reqURL string, w io.Writer) error {
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return fmt.Errorf("Error while creating *http.Request: %q", err)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("Error while making request: %q", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("Unexpected response(%d): %s", resp.StatusCode, statusCode[resp.StatusCode])
	}

	if io.Copy(w, resp.Body); err != nil {
		return fmt.Errorf("error copying response content:%q", err)
	}
	return nil
}
