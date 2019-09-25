package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const (
	baseURL = "http://transparencia.mprn.mp.br/"
)

// We're not trying to retrieve information about "pensionistas" because the data is not there. Ex: http://transparencia.mprn.mp.br/Arquivos/C0007/2019/R0086/38033.pdf?dt=25092019141321
const (
	membrosAtivos        = 82
	membrosInativos      = 83
	servidoresAtivos     = 84
	servidoresInativos   = 85
	colaboradores        = 87
	exerciciosAnteriores = 1143
)

var categories = map[int]string{
	membrosAtivos:        "MembrosAtivos",
	membrosInativos:      "MembrosInativos",
	servidoresAtivos:     "ServidoresAtivos",
	servidoresInativos:   "ServidoresInativos",
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

func link(category, month, year int) (string, error) {
	query := fmt.Sprintf("%s/home/listarAnexos?idanexo=%d&ano=%d", baseURL, category, year)
	resp, err := http.Get(query)
	if err != nil {
		return "", fmt.Errorf("Error trying to download html snippet (%s): %q", query, err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Error trying to make document from reader (%s): %q", query, err)
	}
	sel := doc.Find("tr")
	//Looking for node with a children that contains the query
	for i := range sel.Nodes {
		c := sel.Eq(i).Children()
		if strings.Contains(c.Eq(0).Children().Text(), fmt.Sprintf("%02d-%d", month, year)) {
			f := c.Eq(1).Find("a")
			if href, ok := f.Attr("href"); ok {
				return fmt.Sprintf("%s%s", baseURL, href), nil
			}
		}
	}
	return "", fmt.Errorf("Couldn't find link for the query")
}

func fetchContent(url string) ([]byte, error) {
	respFile, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("couldn't retrieve content (%s): %q", url, err)
	}
	defer respFile.Body.Close()

	if respFile.Header.Get("Content-type") != "application/oleobject" {
		return nil, fmt.Errorf("Request not returning an ODS file(%s): Content-type %s", url, respFile.Header.Get("Content-type"))
	}

	c, err := ioutil.ReadAll(respFile.Body)
	if err != nil {
		return nil, fmt.Errorf("couldn't read content to byte array (%s): %q", url, err)
	}

	return c, nil
}

func main() {
	month := flag.Int("mes", 0, "Mês de referência")
	year := flag.Int("ano", 0, "Ano de referência")

	flag.Parse()

	if *month == 0 || *year == 0 {
		log.Fatalf("Need flags: \"--month=int --year=int\"")
	}

	for catKey, category := range categories {
		link, err := link(catKey, *month, *year)
		if err != nil {
			log.Fatalf("Error retrieving content link (%s, %d-%d): %q\n", category, *month, *year, err)
		}
		c, err := fetchContent(link)
		if err != nil {
			log.Fatalf("Error fetching content (%s, %d-%d): %q\n", category, *month, *year, err)
		}
		err = saveFile(c, *month, *year, category)
		if err != nil {
			log.Fatalf("Error saving content(%s, %d-%d): %q\n", category, *month, *year, err)
		}
	}
}
