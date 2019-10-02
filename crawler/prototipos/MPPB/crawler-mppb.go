package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
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
		log.Fatalf("Need arguments to continue, please try again!")
	}

	for typ, url := range links(*month, *year) {
		c, err := download(url)
		if err != nil {
			log.Fatalf("Error while downloading content: %q\n", err)
			continue
		}
		name := fmt.Sprintf("%s-%d-%d.ods", typ, *month, *year)
		if err = save(c, name); err != nil {
			log.Fatalf("Error while saving to file(%s): %q\n", name, err)
		}
		fmt.Printf("File successfully saved:%s", name)
	}
}

// Generate endpoints able to download
func links(month, year int) map[string]string {
	baseURL := "http://pitagoras.mppb.mp.br/PTMP/"
	links := make(map[string]string, len(tipos))
	links["estagio"] = fmt.Sprintf("%sFolhaPagamentoEstagiarioExercicioMesOds?exercicio=%d&mes=%d", baseURL, year, month)
	links["indenizacoes"] = fmt.Sprintf("%sFolhaVerbaIndenizRemTemporariaOds?mes=%d&exercicio=%d&tipo=", baseURL, month, year)
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
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// Receive a slice of bytes after download, write, nominate file and save
func save(content []byte, name string) error {
	newFile, err := os.Create(name)
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
