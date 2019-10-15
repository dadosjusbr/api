package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/antchfx/htmlquery"
)

const (
	baseURL       = "https://www.tjpb.jus.br/transparencia/gestao-de-pessoas/folha-de-pagamento-de-pessoal"
	fileExtension = ".pdf"
)

var monthStr = map[int]string{
	1:  "Janeiro",
	2:  "Fevereiro",
	3:  "Março",
	4:  "Abril",
	5:  "Maio",
	6:  "Junho",
	7:  "Julho",
	8:  "Agosto",
	9:  "Setembro",
	10: "Outubro",
	11: "Novembro",
	12: "Dezembro",
}

func main() {
	month := flag.Int("mes", 0, "Mês a ser analisado")
	year := flag.Int("ano", 0, "Ano a ser analisado")
	flag.Parse()
	if *month == 0 || *year == 0 {
		log.Fatalf("need arguments to continue, please try again")
	}

	links, err := links(baseURL, *month, *year)
	if err != nil {
		log.Fatalf("Error trying to generate links: %q", err)
	}

	for typ, url := range links {
		name := fmt.Sprintf("%s%s", typ, fileExtension)
		f, err := os.Create(name)
		if err != nil {
			log.Fatalf("error creating file(%s):%q", name, err)
		}
		if download(url, f); err != nil {
			log.Fatalf("error while downloading content: %q", err)
		}
		f.Close()
		fmt.Printf("File successfully saved:%s\n", name)
	}
}

// Generate endpoints able to download
func links(baseURL string, month, year int) (map[string]string, error) {
	links := make(map[string]string)
	doc, err := htmlquery.LoadURL(baseURL)
	if err != nil {
		return nil, fmt.Errorf("error loading doc (%s):%q", baseURL, err)
	}

	var xpath string
	if year > 2012 {
		xpath = fmt.Sprintf("//*[@id=\"arquivos-%04d-mes-%02d\"]//a", year, month)
	} else {
		xpath = fmt.Sprintf("//ul[@id=\"arquivos-%04d\"]//a[contains(text(), '%s %04d')]", year, monthStr[month], year)
	}

	list := htmlquery.Find(doc, xpath)
	if !(len(list) > 0) {
		return nil, fmt.Errorf("couldn't find any file for specified month and year")
	}

	for _, node := range list {
		href := node.Attr[0].Val //href value
		var name string
		if strings.Contains(href, "magistrados") {
			name = fmt.Sprintf("remuneracoes-magistrados-tjpb-%d-%d", month, year)
		} else if strings.Contains(href, "servidores") {
			name = fmt.Sprintf("remuneracoes-servidores-tjpb-%d-%d", month, year)
		} else {
			name = fmt.Sprintf("remuneracoes-tjpb-%d-%d", month, year)
		}
		links[name] = href
	}

	return links, nil
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
