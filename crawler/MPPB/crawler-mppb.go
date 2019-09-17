package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

const (
	membrosAtivos         = "1"
	membrosInativos       = "2"
	servidoresAtivos      = "3"
	servidoresInativos    = "4"
	servidoresDisponiveis = "5"
	aposentados           = "6"
)

func main() {
	month := flag.Int("mes", 0, "Mês a ser analisado")
	year := flag.Int("ano", 0, "Ano a ser analisado")
	flag.Parse()

	monthString := strconv.Itoa(*month)
	yearString := strconv.Itoa(*year)

	links := getLink(monthString, yearString)
	saveToOds(links, monthString, yearString)

}

func getLink(monthString string, yearString string) map[string]string {
	baseURL := "http://pitagoras.mppb.mp.br/PTMP/"
	links := map[string]string{
		"membrosAtivos":         baseURL + "FolhaPagamentoExercicioMesNewOds?" + "mes=" + monthString + "&exercicio=" + yearString + "&tipo=" + membrosAtivos,
		"membrosInativos":       baseURL + "FolhaPagamentoExercicioMesNewOds?" + "mes=" + monthString + "&exercicio=" + yearString + "&tipo=" + membrosInativos,
		"servidoresAtivos":      baseURL + "FolhaPagamentoExercicioMesNewOds?" + "mes=" + monthString + "&exercicio=" + yearString + "&tipo=" + servidoresAtivos,
		"servidoresInativos":    baseURL + "FolhaPagamentoExercicioMesNewOds?" + "mes=" + monthString + "&exercicio=" + yearString + "&tipo=" + servidoresInativos,
		"servidoresDisponiveis": baseURL + "FolhaPagamentoExercicioMesNewOds?" + "mes=" + monthString + "&exercicio=" + yearString + "&tipo=" + servidoresDisponiveis,
		"aposentados":           baseURL + "FolhaPagamentoExercicioMesNewOds?" + "mes=" + monthString + "&exercicio=" + yearString + "&tipo=" + aposentados,
		"anteriores":            baseURL + "FolhaExercicioAnteriorMesNewOds?exercicio=" + yearString + "&mes=" + monthString,
		"estagio":               baseURL + "FolhaPagamentoEstagiarioExercicioMesOds?mes=" + monthString + "&exercicio=" + yearString,
	}

	return links

}

func saveToOds(links map[string]string, monthString string, yearString string) {

	for key, value := range links {
		resp, err := http.Get(value)
		if err != nil {
			fmt.Println("Error while getting the response")
		}

		newFile, err2 := os.Create(key + "-" + monthString + "-" + yearString)
		if err2 != nil {
			fmt.Println("Error while creating the file")
		}
		defer newFile.Close()

		bodySave, err3 := ioutil.ReadAll(resp.Body)
		if err3 != nil {
			fmt.Println("Error writing to file")
		}

		newFile.Write(bodySave)

		defer resp.Body.Close()
	}

}
