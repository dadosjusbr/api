package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
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

	for typ, url := range links(*month, *year) {
		c, err := download(url)
		if err != nil {
			fmt.Printf("Error while downloading content: %q\n", err)
			continue
		}
		err = saveToOds(c, typ, *month, *year)
		if err != nil {
			fmt.Printf("Error while saving to file: %q\n", err)
		}
	}
}

// Generate endpoints able to download
func links(month, year int) map[string]string {
	baseURL := "http://pitagoras.mppb.mp.br/PTMP/"
	links := make(map[string]string, len(tipos)+1)
	links["estagio"] = fmt.Sprintf("%sFolhaPagamentoEstagiarioExercicioMesOds?exercicio=%d&mes=%d", baseURL, year, month)
	for t, id := range tipos {
		links[t] = fmt.Sprintf("%sFolhaPagamentoExercicioMesNewOds?mes=%d&exercicio=%d&tipo=%d", baseURL, month, year, id)
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

// Receive a slice of bytes after download, write, nominate file and save
func saveToOds(content []byte, typ string, monthString, yearString int) error {
	newFile, err := os.Create(fmt.Sprintf("%s-%d-%-d", typ, monthString, yearString))
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
