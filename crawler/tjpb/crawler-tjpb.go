package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

const baseURL = "https://www.tjpb.jus.br/transparencia/gestao-de-pessoas/folha-de-pagamento-de-pessoal"

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

var netClient = &http.Client{
	Timeout: time.Second * 60,
}

func main() {
	month := flag.Int("mes", 0, "Mês a ser analisado")
	year := flag.Int("ano", 0, "Ano a ser analisado")
	flag.Parse()
	if *month == 0 || *year == 0 {
		log.Fatalf("Need all arguments to continue, please try again: \"go run crawler-tjpb.go --mes=int --ano=int\"")
	}

	links, err := links(*month, *year)
	if err != nil {
		log.Fatalf("Error trying to find links: %q", err)
	}

	for typ, url := range links {
		fileName := fmt.Sprintf("%s.pdf", typ)
		f, err := os.Create(fileName)
		if err != nil {
			log.Fatalf("error creating file(%s):%q", fileName, err)
		}
		if download(url, f); err != nil {
			f.Close()
			os.Remove(fileName)
			log.Fatalf("error while downloading content (%02d-%04d): %q", *month, *year, err)
		}
		f.Close()
		fmt.Printf("File successfully saved:%s\n", fileName)
	}
}

// Make xpath query to find links of interest for a given month and year.
func linkNodes(month, year int) ([]*html.Node, error) {
	resp, err := netClient.Get(baseURL)
	if err != nil {
		return nil, fmt.Errorf("error making GET request to %s: %q", baseURL, err)
	}
	defer resp.Body.Close()

	doc, err := htmlquery.Parse(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error loading doc (%s): %q", baseURL, err)
	}

	xpath := fmt.Sprintf("//*[@id=\"arquivos-%04d-mes-%02d\"]//a", year, month)
	if year <= 2012 {
		xpath = fmt.Sprintf("//ul[@id=\"arquivos-%04d\"]//a[contains(text(), '%s %04d')]", year, monthStr[month], year)
	}

	nodeList := htmlquery.Find(doc, xpath)
	if !(len(nodeList) > 0) {
		return nil, fmt.Errorf("couldn't find any link for specified month and year")
	}
	return nodeList, nil
}

// Generate a map of endpoints with a named reference to their content as a key.
func links(month, year int) (map[string]string, error) {
	links := make(map[string]string)
	nodeList, err := linkNodes(month, year)
	if err != nil {
		return nil, err
	}

	for _, node := range nodeList {
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

// Download from endpoint and copy content to an io.Writer.
func download(url string, w io.Writer) error {
	resp, err := netClient.Get(url)
	if err != nil {
		return fmt.Errorf("error downloading file:%q", err)
	}
	defer resp.Body.Close()
	if io.Copy(w, resp.Body); err != nil {
		return fmt.Errorf("error copying response content:%q", err)
	}
	return nil
}
