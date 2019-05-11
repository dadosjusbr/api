package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/dadosjusbr/remuneracao-magistrados/crawler"

	"github.com/kelseyhightower/envconfig"
)

type config struct {
	SpreadsheetsPath string `envconfig:"LOCAL_SPREADSHEETS_PATH"`
	Month            string `envconfig:"MONTH"`
	Year             string `envconfig:"YEAR"`
}

func main() {
	var conf config
	err := envconfig.Process("remuneracao-magistrados", &conf)
	if err != nil {
		log.Fatal(err.Error())
	}

	indexPath, err := generateIndexMock(conf.SpreadsheetsPath)
	if err != nil {
		log.Fatal(err.Error())
	}

	results, err := crawler.Crawl(indexPath)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println("RESULTS NAMES: ")
	for _, result := range results {
		fmt.Println(result.Name)
	}

	fmt.Println(indexPath)
}

// generateIndexMock create a index.html with the local paths of the files inside the given directory path
// so that the crawler can get then.
func generateIndexMock(spreadsheetsPath string) (string, error) {
	filesInfo, err := ioutil.ReadDir(spreadsheetsPath)
	if err != nil {
		return "", err
	}
	var files []template.URL

	for _, file := range filesInfo {
		path := template.URL(fmt.Sprintf("file://%s/%s", spreadsheetsPath, file.Name()))
		files = append(files, path)
	}

	const tpl = `
		<!DOCTYPE html>
		<html>
			<head>
				<meta charset="UTF-8">
				<title>Any title</title>
			</head>
			<body>
				<table>
					<tr>
					{{range .}}
						<td>
							<a href="{{ . }}" target="_blank" rel="alternate noopener">any text</a>
						</td>
					{{end}}
					</tr>
				</table>
			</body>
		</html>`

	t, err := template.New("webpage").Parse(tpl)
	if err != nil {
		return "", err
	}

	f, err := os.Create("index.html")
	if err != nil {
		return "", err
	}

	err = t.Execute(f, files)
	if err != nil {
		return "", err
	}

	indexPath, err := filepath.Abs("./index.html")
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("file://%s", indexPath), nil
}
