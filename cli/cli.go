package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/dadosjusbr/remuneracao-magistrados/db"
	"github.com/dadosjusbr/remuneracao-magistrados/parser"
	"github.com/dadosjusbr/remuneracao-magistrados/processor"
	"github.com/dadosjusbr/remuneracao-magistrados/store"
	"github.com/kelseyhightower/envconfig"
)

type config struct {
	PCloudUsername   string `envconfig:"PCLOUD_USERNAME"`
	PCloudPassword   string `envconfig:"PCLOUD_PASSWORD"`
	ParserURL        string `envconfig:"PARSER_URL"`
	Month            int    `envconfig:"MONTH"`
	Year             int    `envconfig:"YEAR"`
	SpreadsheetsPath string `envconfig:"LOCAL_SPREADSHEETS_PATH"`
	MonthURL         string `envconfig:"MONTH_URL"`
	DBUrl            string `envconfig:"DADOSJUSBR_DB_URL"`
}

func main() {
	var conf config
	err := envconfig.Process("remuneracao-magistrados", &conf)
	if err != nil {
		log.Fatal(err.Error())
	}

	pcloudClient, err := store.NewPCloudClient(conf.PCloudUsername, conf.PCloudPassword)
	if err != nil {
		log.Fatal("ERROR: ", err.Error())
	}

	parserClient := parser.NewServiceClient(conf.ParserURL)

	dbClient, err := db.NewClient(conf.DBUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer dbClient.CloseConnection()

	var indexPath string

	if conf.MonthURL != "" {
		indexPath = conf.MonthURL
	} else {
		p, err := generateIndexMock(conf.SpreadsheetsPath)
		if err != nil {
			log.Fatal(err.Error())
		}
		defer os.Remove(p)
		indexPath = fmt.Sprintf("file://%s", p)
	}
	fmt.Printf("Processing spreadshets from: %s\n", indexPath)

	err = processor.Process(indexPath, conf.Month, conf.Year, pcloudClient, parserClient, dbClient)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Month successfuly published.")
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
	return indexPath, nil
}
