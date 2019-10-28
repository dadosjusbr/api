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

var netClient = &http.Client{
	Timeout: time.Second * 60,
}

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
		log.Fatalf("Need all arguments to continue, please try again: \"go run crawler-tjpb.go --mes=int --ano=int\"")
	}

	doc, err := loadURL(baseURL)
	if err != nil {
		log.Fatalf("Error trying to load URL: %q", err)
	}

	nodes, err := findInterestNodes(doc, *month, *year)
	if err != nil {
		log.Fatalf("Error trying to find link nodes of interest: %q", err)
	}

	for typ, url := range fileURLsFromInterestNodes(nodes, *month, *year) {
		if err = save(typ, url); err != nil {
			log.Fatalf("Error trying to save file: %q", err)
		}
	}
}

// loadURL loads the HTML document from the specified URL, parses it and returns the root node of the document.
func loadURL(baseURL string) (*html.Node, error) {
	resp, err := netClient.Get(baseURL)
	if err != nil {
		return nil, fmt.Errorf("error making GET request to %s: %q", baseURL, err)
	}
	defer resp.Body.Close()

	doc, err := htmlquery.Parse(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error loading doc (%s): %q", baseURL, err)
	}
	return doc, nil
}

// fileURLsFromInterestNodes generates a map of endpoints with a named reference to their content as a key.
func fileURLsFromInterestNodes(nodes []*html.Node, month, year int) map[string]string {
	links := make(map[string]string)
	for _, node := range nodes {
		val := node.Attr[0].Val
		name := fileName(val, month, year)
		links[name] = val
	}
	return links
}

// findInterestNodes makes a xpath query to find link nodes of interest for a given month and year.
func findInterestNodes(doc *html.Node, month, year int) ([]*html.Node, error) {
	//Sets xpath for interest nodes depending on year and month.
	xpath := fmt.Sprintf("//*[@id=\"arquivos-%04d-mes-%02d\"]//a", year, month)
	if year <= 2012 {
		xpath = fmt.Sprintf("//ul[@id=\"arquivos-%04d\"]//a[contains(text(), '%s %04d')]", year, monthStr[month], year)
	}

	nodeList := htmlquery.Find(doc, xpath)
	if len(nodeList) == 0 {
		return nil, fmt.Errorf("couldn't find any link for %02d-%04d", month, year)
	}
	return nodeList, nil
}

// fileName generates name for href content.
func fileName(href string, month, year int) string {
	if strings.Contains(href, "magistrados") {
		return fmt.Sprintf("remuneracoes-magistrados-tjpb-%02d-%04d", month, year)
	} else if strings.Contains(href, "servidores") {
		return fmt.Sprintf("remuneracoes-servidores-tjpb-%02d-%04d", month, year)
	}
	return fmt.Sprintf("remuneracoes-tjpb-%02d-%04d", month, year)
}

// save downloads content from url and save it on a file.
func save(typ string, url string) error {
	fileName := fmt.Sprintf("%s.pdf", typ)
	f, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("error creating file(%s):%q", fileName, err)
	}
	defer f.Close()

	if err = download(url, f); err != nil {
		os.Remove(fileName)
		return fmt.Errorf("error while downloading content: %q", err)
	}
	return nil
}

// download gets content from url and copy it to an io.Writer.
func download(url string, w io.Writer) error {
	resp, err := netClient.Get(url)
	if err != nil {
		return fmt.Errorf("error making get request to (%s): %q", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected response status code (%s): %d - %s", url, resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	if io.Copy(w, resp.Body); err != nil {
		return fmt.Errorf("error copying response content to file: %q", err)
	}
	return nil
}
