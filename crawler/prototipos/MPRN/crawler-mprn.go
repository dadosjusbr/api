package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const baseURL = "http://transparencia.mprn.mp.br/"
const requestListURL = "/home/listarAnexos"

const (
	membrosAtivos = 82 + iota
	membrosInativos
	servidoresAtivos
	servidoresInativos
	pensionistas
	colaboradores
	exerciciosAnteriores = 1143
)

var categories = map[int]string{
	membrosAtivos:        "MembrosAtivos",
	membrosInativos:      "MembrosInativos",
	servidoresAtivos:     "ServidoresAtivos",
	servidoresInativos:   "ServidoresInativos",
	pensionistas:         "Pensionistas",
	colaboradores:        "Colaboradores",
	exerciciosAnteriores: "ExerciciosAnteriores",
}

func saveFile(c []byte, month int, year int, category string) error {
	//Create a new file in the cwd
	fileName := fmt.Sprintf("%s-%02d-%d.ods", category, month, year)
	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("Error creating file(%s): %q", fileName, err)
	}
	defer file.Close()

	//Write to file
	if _, err = file.Write(c); err != nil {
		return fmt.Errorf("Error writing to file (%s): %q", fileName, err)
	}
	return nil
}

func fetchContent(month int, year int, category int) ([]byte, error) {
	query := fmt.Sprintf("%s%s?idanexo=%d&ano=%d", baseURL, requestListURL, category, year)
	resp, err := http.Get(query)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	document, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}
	sel := document.Find("tr")
	mesAno := fmt.Sprintf("%02d-%d", month, year)
	var contentURL string
	//Looking for node with a children that contains the query
	for i := range sel.Nodes {
		nodeChildren := sel.Eq(i).Children()
		if strings.Contains(nodeChildren.Eq(0).Children().Text(), mesAno) {
			fileNode := nodeChildren.Eq(1).Find("a")
			if href, ok := fileNode.Attr("href"); ok {
				contentURL = fmt.Sprintf("%s%s", baseURL, href)
			}
		}

	}
	if contentURL == "" {
		return nil, fmt.Errorf("Couldn't find link for the query")
	}
	respFile, err := http.Get(contentURL)
	if err != nil {
		return nil, fmt.Errorf("couldn't retrieve content: %s - %d-%d", categories[category], month, year)
	}
	defer respFile.Body.Close()

	if respFile.Header.Get("Content-type") != "application/oleobject" {
		return nil, fmt.Errorf("Request not returning an ODS file(%s): Content-type %s", contentURL, resp.Header.Get("Content-type"))
	}

	c, err := ioutil.ReadAll(respFile.Body)
	if err != nil {
		return nil, fmt.Errorf("couldn't read content to byte array: %s - %d-%d", categories[category], month, year)
	}

	return c, nil
}

func main() {
	month := flag.Int("mes", 0, "Mês de referência")
	year := flag.Int("ano", 0, "Ano de referência")

	flag.Parse()

	if *month == 0 || *year == 0 {
		fmt.Println("Need flags: \"--month=int --year=int\"")
	}

	for catKey, category := range categories {
		c, err := fetchContent(*month, *year, catKey)
		if err != nil {
			fmt.Printf("Error fetching content (%s, %d-%d): %q\n", category, *month, *year, err)
			continue
		}
		err = saveFile(c, *month, *year, category)
		if err != nil {
			fmt.Printf("Error saving content(%s, %d-%d): %q\n", category, *month, *year, err)
			continue
		}
	}
}
