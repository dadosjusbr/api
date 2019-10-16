package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

const (
	baseURL          = "http://pitagoras.mppb.mp.br/PTMP/"
	tipoEstagiarios  = "estagiarios"
	tipoIndenizacoes = "indenizacoes"
)

var (
	tipos = map[string]int{
		"membrosAtivos":         1,
		"membrosInativos":       2,
		"servidoresAtivos":      3,
		"servidoresInativos":    4,
		"servidoresDisponiveis": 5,
		"aposentados":           6,
	}
)

func main() {
	month := flag.Int("mes", 0, "Mês a ser analisado")
	year := flag.Int("ano", 0, "Ano a ser analisado")
	flag.Parse()
	if *month == 0 || *year == 0 {
		log.Fatalf("need arguments to continue, please try again")
	}
	for typ, url := range links(baseURL, *month, *year) {
		name := fmt.Sprintf("%s-%d-%d.ods", typ, month, year)
		f, err := os.Create(name)
		if err != nil {
			log.Fatalf("error creating file(%s):%q", name, err)
		}
		if download(url, f); err != nil {
			log.Fatalf("error while downloading content: %q", err)
		}
		f.Close()
		fmt.Printf("File successfully saved:%s", name)
	}
}

// Generate endpoints able to download
func links(baseURL string, month, year int) map[string]string {
	links := make(map[string]string)
	links[tipoEstagiarios] = fmt.Sprintf("%sFolhaPagamentoEstagiarioExercicioMesOds?exercicio=%d&mes=%d", baseURL, year, month)
	links[tipoIndenizacoes] = fmt.Sprintf("%sFolhaVerbaIndenizRemTemporariaOds?mes=%d&exercicio=%d&tipo=", baseURL, month, year)
	for t, id := range tipos {
		links[t] = fmt.Sprintf("%sFolhaPagamentoExercicioMesNewOds?mes=%d&exercicio=%d&tipo=%d", baseURL, month, year, id)
	}
	return links
}

func download(url string, w io.Writer) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error downloading file:%q", err)
	}
	defer resp.Body.Close()
	if io.Copy(w, resp.Body); err != nil {
		return fmt.Errorf("error copying response content:%q", err)
	}
	return nil
}
