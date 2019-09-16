package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"os"
	"strconv"
)

func main() {
	month := flag.Int("mes", 0, "Mês a ser analisado")
	year := flag.Int("ano", 0, "Ano a ser analisado")
	flag.Parse()

	monthString := strconv.Itoa(*month)
	yearString := strconv.Itoa(*year)

	resp, _ := http.Get("http://pitagoras.mppb.mp.br/PTMP/FolhaPagamentoExercicioMesNewOds?mes=" + monthString + "&exercicio=" + yearString + "&tipo=1")
	dump, _ := httputil.DumpResponse(resp, false)
	fmt.Println(string(dump))
	saveToOds(resp, "teste")
	defer resp.Body.Close()
}

func saveToOds(resp *http.Response, file string) error {

	newFile, err := os.Create(file)
	defer newFile.Close()

	bodySave, err := ioutil.ReadAll(resp.Body)

	newFile.Write(bodySave)

	return err

}
