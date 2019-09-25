package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
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
		log.Fatalf("Need arguments to continue, please try again!")
	}

	// First, crawl everything which follows the same standard (a.k.a. everyone but pensionistas).
	for typ, url := range links(*month, *year) {
		c, err := download(url)
		if err != nil {
			log.Fatalf("Error while downloading content: %q\n", err)
		}
		fmt.Printf("File successfully downloaded:%s", url)
		name := fmt.Sprintf("%s-%d-%d.ods", typ, *month, *year)
		if err = save(c, name); err != nil {
			log.Fatalf("Error while saving to file(%s): %q\n", name, err)
		}
		fmt.Printf("File successfully saved:%s", name)
	}
	// Then we crawl pensionistas.
	pURL, err := pensionistaLink(*month, *year)
	c, err := download(pURL)
	if err != nil {
		log.Fatalf("Error while downloading content: %q\n", err)
	}
	fmt.Printf("File successfully downloaded:%s", pURL)
	name := fmt.Sprintf("pensionista-%d-%d.xls", *month, *year)
	if err = save(c, name); err != nil {
		log.Fatalf("Error while saving to file: %q\n", err)
	}
	fmt.Printf("File successfully saved:%s", name)
}

func pensionistaLink(month, year int) (string, error) {
	url := fmt.Sprintf("https://sistemas.mpal.mp.br/transparencia/principal/publicacoes/69?competencia=%d&busca=%s", year, meses[month])
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("error downloading URL(%s): %q", url, err)
	}
	defer resp.Body.Close()
	linksURL := collectlinks.All(resp.Body)
	for id := range linksURL {
		if strings.Contains(linksURL[id], "download") {
			return fmt.Sprintf("https://sistemas.mpal.mp.br/%s", linksURL[id]), nil
		}
	}
	return "", fmt.Errorf("could not find download link for pensionistas in page %s", url)
}

// https://sistemas.mpal.mp.br/transparencia/contracheque/index/65?tipo=membrosativos&mes=8&ano=2019&busca=&download=ods
// Generate endpoints able to download
func links(month, year int) map[string]string {
	links := make(map[string]string)
	for t, id := range tipos {
		// There are three download options: ods, json, csv.
		links[t] = fmt.Sprintf("https://sistemas.mpal.mp.br/transparencia/contracheque/index/%d?tipo=%s&mes=%d&ano=%d&busca=&download=ods", id, strings.ToLower(t), month, year)
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
