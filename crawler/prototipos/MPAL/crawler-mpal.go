package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/jackdanger/collectlinks"
)

var tipos = map[string]int{
	"membrosAtivos":           65,
	"membrosInativos":         66,
	"todosServidoresAtivos":   67,
	"todosServidoresInativos": 68,
	"colaboradores":           70,
	//"aposentados":            6, precisa varrer o site e baixar os links
}

var meses = map[int]string{
	1:  "janeiro",
	2:  "fevereiro",
	3:  "março",
	4:  "abril",
	5:  "maio",
	6:  "junho",
	7:  "julho",
	8:  "agosto",
	9:  "setembro",
	10: "outubro",
	11: "novembro",
	12: "dezembro",
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

// https://sistemas.mpal.mp.br/transparencia/contracheque/index/65?tipo=membrosativos&mes=8&ano=2019&busca=&download=ods
// Generate endpoints able to download
func links(month, year int) map[string]string {
	linkPensionista := ""
	response, err := http.Get(fmt.Sprint("https://sistemas.mpal.mp.br/transparencia/principal/publicacoes/69?competencia=%d&busca=%s", year, meses[month]))
	if err != nil {
		fmt.Printf("Error while downloading content: %q\n", err)
	} else {
		linksURL := collectlinks.All(response.Body)
		for id := range linksURL {
			if strings.Contains(linksURL[id], "download") {
				linkPensionista = linksURL[id]
			}
		}
	}
	typFile := "ods"
	baseURL := "https://sistemas.mpal.mp.br/transparencia/contracheque/index/"
	baseURLPensionista := "https://sistemas.mpal.mp.br/"
	links := make(map[string]string, len(tipos)+1)
	links["pensionistas"] = baseURLPensionista + linkPensionista
	for t, id := range tipos {
		links[t] = fmt.Sprintf("%s%d?tipo=%s&mes=%d&ano=%d&busca=&download=%s", baseURL, id, strings.ToLower(t), month, year, typFile)
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
	newFile, err := os.Create(fmt.Sprintf("%s-%d-%d.ods", typ, monthString, yearString))
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
