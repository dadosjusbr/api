package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

func main() {
	month := flag.Int("mes", 0, "Mês a ser analisado")
	year := flag.Int("ano", 0, "Ano a ser analisado")
	name := flag.String("nome", "", "Used for login purposes")
	cpf := flag.String("cpf", "", "used for login purpose. format xxx.xxx.xxx-xx")
	flag.Parse()
	if *month == 0 || *year == 0 || *cpf == "" || *name == "" {
		log.Fatalf("Need all arguments to continue, please try again\n")
	}

	acessCode, err := login(*name, *cpf)
	if err != nil {
		log.Fatalf("login error: %q", err)
	}

	data, err := queryData(acessCode, *month, *year)
	if err != nil {
		log.Fatalf("Query data error: %q", err)
	}

	dataDesc := fmt.Sprintf("remuneracoes-trepb-%02d-%04d", *month, *year)

	if err = save(dataDesc, data); err != nil {
		log.Fatalf("Error saving data to file: %q", err)
	}

}

// save downloads content from url and save it on a file.
func save(desc string, data []*html.Node) error {
	fileName := fmt.Sprintf("%s.html", desc)
	f, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("error creating file(%s):%q", fileName, err)
	}
	defer f.Close()

	for _, node := range data {
		nodeReader := strings.NewReader(htmlquery.OutputHTML(node, true))
		if io.Copy(f, nodeReader); err != nil {
			os.Remove(fileName)
			return fmt.Errorf("error copying response content to file: %q", err)
		}
	}
	return nil
}

// queryData query server for data of a specified month and year.
func queryData(acessCode string, month, year int) ([]*html.Node, error) {
	query := fmt.Sprintf(`acao=AnexoVIII&folha=&valida=true&toExcel=false&chaveDeAcesso=%s&mes=%d&ano=%04d`, acessCode, month, year)
	queryURL := fmt.Sprintf(`http://apps.tre-pb.jus.br/transparenciaDadosServidores/infoServidores?%s`, query)
	doc, err := loadURL(queryURL)
	if err != nil {
		return nil, fmt.Errorf("error while loading url: %q", err)
	}

	tables, err := htmlquery.QueryAll(doc, "//table")
	if err != nil {
		return nil, fmt.Errorf("error while making query for data tables: %q", err)
	}
	if len(tables) == 0 {
		return nil, fmt.Errorf("couldn't find any data tables")
	}
	return tables, nil
}
