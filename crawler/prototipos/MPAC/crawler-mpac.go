package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

// Category for each type of spreedsheet.
// Obtained from http://transparencia.mpac.mp.br/categoria/24
// Every category got it's own page http://transparencia.mpac.mp.br/categoria/${categoryNumber}
const (
	membrosAtivos = iota + 112
	membrosInativos
	servidoresAtivos
	servidoresInativos
	pensionistas
	colaboradores
	exerciciosAnteriores
	indenizacoes
)

var categories = map[int]string{
	membrosAtivos:        "MembrosAtivos",
	membrosInativos:      "MembrosInativos",
	servidoresAtivos:     "ServidoresAtivos",
	servidoresInativos:   "ServidoresInativos",
	pensionistas:         "Pensionistas",
	colaboradores:        "Colaboradores",
	exerciciosAnteriores: "ExerciciosAnteriores",
	indenizacoes:         "Indenizacoes",
}

const (
	baseURL       = "http://transparencia.mpac.mp.br/categoria_arquivos/"
	fileExtension = ".ods"
)

var client *http.Client = &http.Client{
	CheckRedirect: func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	},
}

func fetchContent(key int, year int, month int) ([]byte, error) {
	//Retrieve file location
	body := strings.NewReader(fmt.Sprintf("ano=%d&numMes=%d", year, month))
	aURL := fmt.Sprintf("%s%d", baseURL, key)
	res, err := client.Post(aURL, "application/x-www-form-urlencoded", body)

	if err != nil {
		return nil, fmt.Errorf("Error making a post request(%s): %q", aURL, err)
	}
	defer res.Body.Close()

	if res.StatusCode != 302 {
		return nil, fmt.Errorf("Resource not found(%s): statusCode(%d)", aURL, res.StatusCode)
	}

	//Download target file
	targetURL := res.Header.Get("Location")
	target, err := http.Get(targetURL)
	if err != nil {
		return nil, fmt.Errorf("Error making get request (%s): %q", targetURL, err)
	}
	defer target.Body.Close()

	if target.Header.Get("Content-type") != "application/vnd.oasis.opendocument.spreadsheet; charset=UTF-8" {
		return nil, fmt.Errorf("Request not returning an ODS file(%s): Content-type %s", targetURL, target.Header.Get("Content-type"))
	}

	//Transform body in a slice of bytes
	targetBody, err := ioutil.ReadAll(target.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading response body (%s): %q", targetURL, err)
	}

	return targetBody, nil
}

func saveFile(c []byte, year int, month int, category string) error {
	//Create a new file in the cwd
	fileName := fmt.Sprintf("%s-%d-%d%s", category, month, year, fileExtension)
	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("Error creating file(%s): %q", fileName, err)
	}
	defer file.Close()

	//Write to file
	if _, err = file.Write(c); err != nil {
		return fmt.Errorf("Error writing to file (%s): %q", fileName, err)
	}
	return nil
}

func main() {
	month := flag.Int("mes", 0, "Mês de referência")
	year := flag.Int("ano", 0, "Ano de referência")
	flag.Parse()

	if (*month == 0) || (*year == 0) {
		fmt.Println("Need flags '--mes --ano' to work")
		return
	}

	for key, category := range categories {
		c, err := fetchContent(key, *year, *month)
		if err != nil {
			fmt.Printf("Error retrieving resource: (%s %d-%d): %q\n", category, *month, *year, err)
			continue
		}
		if err = saveFile(c, *year, *month, category); err != nil {
			fmt.Printf("Error saving spreedsheet to file (%s %d-%d): %q\n", category, *month, *year, err)
		}
	}
}
