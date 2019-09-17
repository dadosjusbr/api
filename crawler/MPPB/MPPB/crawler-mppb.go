package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

var tipos = map[string]int{
	"membrosAtivos":         1,
	"membrosInativos":       2,
	"servidoresAtivos":      3,
	"servidoresInativos":    4,
	"servidoresDisponiveis": 5,
	"aposentados":           6,
}

func main() {
	month := flag.Int("mes", 0, "Mês a ser analisado")
	year := flag.Int("ano", 0, "Ano a ser analisado")
	flag.Parse()

	if *month == 0 || *year == 0 {
		fmt.Println("Need arguments to continue, please try again!")
		os.Exit(1)
	}

	monthString := strconv.Itoa(*month)
	yearString := strconv.Itoa(*year)

	links := getLinks(monthString, yearString)

	for key, value := range links {
		c, err := download(value)
		if err != nil {
			fmt.Printf("Error while downloading content: %q\n", err)
			continue
		}
		err = saveToOds(c, monthString, yearString, key)
		if err != nil {
			fmt.Printf("Error while saving to file: %q\n", err)
		}
	}
}

func getLinks(month, year string) map[string]string {
	baseURL := "http://pitagoras.mppb.mp.br/PTMP/"
	links := make(map[string]string, len(tipos)+2)
	links["anteriores"] = fmt.Sprintf("%sFolhaExercicioAnteriorMesNewOds?exercicio=%s&mes=%s", baseURL, year, month)
	links["estagio"] = fmt.Sprintf("%sFolhaPagamentoEstagiarioExercicioMesOds?exercicio=%s&mes=%s", baseURL, year, month)
	for t, id := range tipos {
		links[t] = fmt.Sprintf("%sFolhaPagamentoExercicioMesNewOds?mes=%s&exercicio=%s&tipo=%d", baseURL, month, year, id)
	}
	return links
}

func download(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bodySave, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return bodySave, nil
}

func saveToOds(content []byte, monthString, yearString, key string) error {

	newFile, err := os.Create(key + "-" + monthString + "-" + yearString)
	if err != nil {
		return err
	}
	defer newFile.Close()

	_, err = newFile.Write(content)
	if err != nil {
		return err
	}
	return nil
}
