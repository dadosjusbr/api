package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// Categoria para cada tipo de membro
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

const (
	baseURL       = "http://transparencia.mpac.mp.br/categoria_arquivos/"
	fileExtension = ".ods"
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

func requireContentURL(key int, bodyStr string) string {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	body := strings.NewReader(bodyStr)
	URL := baseURL + strconv.Itoa(key)
	req, _ := http.NewRequest("POST", URL, body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "*/*")
	//dump, _ := httputil.DumpRequest(req, true)
	//fmt.Println(string(dump))

	res, _ := client.Do(req)
	defer res.Body.Close()

	if res.StatusCode == 302 {
		return res.Header.Get("Location")
	}

	return ""
}

func saveFile(aURL string, fileName string) error {
	//Download target file
	target, err := http.Get(aURL)
	defer target.Body.Close()
	if err != nil {
		fmt.Printf("Error saving file: %s", fileName)
		return err
	}
	if target.Header.Get("Content-type") == "application/pdf; charset=UTF-8" {
		fmt.Printf("Request not returning a ODS file: %s\nFile-type: %s\n", fileName, target.Header.Get("Content-type"))
		return nil
	}
	//Create a new file in the cwd
	file, err := os.Create(fileName)
	defer file.Close()
	if err != nil {
		fmt.Printf("Error creating file: %s", fileName)
		return err
	}

	//Transform body in a slice of bytes
	targetBody, err := ioutil.ReadAll(target.Body)
	if err != nil {
		fmt.Printf("Error reading body from target: %s", aURL)
		return err
	}

	//Write to file
	file.Write(targetBody)
	return nil
}

func main() {
	monthPtr := flag.Int("mes", 0, "Mês de referência")
	yearPtr := flag.Int("ano", 0, "Ano de referência")
	flag.Parse()

	missingArguments := (*monthPtr == 0) || (*yearPtr == 0)

	month := strconv.Itoa(*monthPtr)
	year := strconv.Itoa(*yearPtr)

	if missingArguments {
		fmt.Println("Need flags '--mes --ano' to work\n")
	} else {
		bodyStr := "ano=" + year + "&" + "numMes=" + month
		for key, category := range categories {
			contentURL := requireContentURL(key, bodyStr)
			if contentURL != "" {
				fileName := category + "-" + month + "-" + year + fileExtension
				_ = saveFile(contentURL, fileName)
			} else {
				fmt.Printf("Error retrieving remunerações for: %s %s/%s\n", category, year, month)
			}
		}
	}
}
