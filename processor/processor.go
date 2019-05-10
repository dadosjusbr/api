package processor

import (
	"archive/zip"
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/dadosjusbr/remuneracao-magistrados/crawler"
	"github.com/dadosjusbr/remuneracao-magistrados/email"
	"github.com/dadosjusbr/remuneracao-magistrados/packager"
	"github.com/dadosjusbr/remuneracao-magistrados/parser"
	"github.com/dadosjusbr/remuneracao-magistrados/store"
)

const (
	emailFrom = "no-reply@dadosjusbr.com"
	emailTo   = "dadosjusbrops@googlegroups.com"
	subject   = "remuneracao-magistrados error"
)

const remuneracaoPath = "http://www.cnj.jus.br/transparencia/remuneracao-dos-magistrados/remuneracao-"

var months = map[int]string{
	1:  "janeiro",
	2:  "fevereiro",
	3:  "marco",
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

// Process download, parse, save and publish data of one month.
func Process(month, year int, emailClient *email.Client, pcloudClient *store.PCloudClient, parser *parser.ServiceClient) {
	//TODO: this function shuld return an error if something goes wrong.
	// Download files from CNJ.
	results, err := crawler.Crawl(fmt.Sprintf("%s%s-%d", remuneracaoPath, months[month], year))
	if err != nil {
		if err := emailClient.Send(emailFrom, emailTo, subject, err.Error()); err != nil {
			fmt.Println("ERROR: " + err.Error())
		}
		fmt.Println("CRAWLING ERROR: " + err.Error())
		return
	}

	// Parsing.
	parsingST := time.Now()

	// Create a buffer to write our archive to.
	var spreadsheetZipBuf bytes.Buffer
	spreadsheetZipWriter := zip.NewWriter(&spreadsheetZipBuf)

	var spreadsheetContents [][]byte
	for _, r := range results {
		zipFile, err := spreadsheetZipWriter.Create(r.Name)
		if err != nil {
			log.Fatal(err)
		}
		_, err = zipFile.Write(r.Body)
		if err != nil {
			// TODO: send email.
			log.Fatal(err)
		}
		spreadsheetContents = append(spreadsheetContents, r.Body)
	}
	if err := spreadsheetZipWriter.Close(); err != nil {
		log.Fatal(err)
	}
	csv, schema, err := parser.Parse(spreadsheetContents)
	if err != nil {
		// TODO: Send an email.
		log.Fatal(err)
	}

	rl, err := pcloudClient.Put("2018-04-raw.zip", &spreadsheetZipBuf)
	if err != nil {
		if err := emailClient.Send(emailFrom, emailTo, subject, err.Error()); err != nil {
			fmt.Println("ERROR: " + err.Error())
		}
		fmt.Println("ERROR: " + err.Error())
		return
	}
	fmt.Printf("Parsing OK (%s). Took %v\n", rl, time.Now().Sub(parsingST))

	// Packaging.
	fmt.Println("Start packaging")
	packagingST := time.Now()
	// TODO: Remove this hardcoded package name. Should be based on the worker selected work (timestamp or past).
	datapackage, err := packager.Pack(fmt.Sprintf("%d-%d", year, month), schema, csv)
	if err != nil {
		if err := emailClient.Send(emailFrom, emailTo, subject, err.Error()); err != nil {
			fmt.Println("ERROR: " + err.Error())
		}
		fmt.Println("ERROR: " + err.Error())
		return
	}
	fmt.Printf("Packaging OK. Took: %s\n", time.Now().Sub(packagingST))
	// Publishing.
	publishingST := time.Now()
	dpl, err := pcloudClient.Put("2018-04-datapackage.zip", bytes.NewReader(datapackage))
	if err != nil {
		if err := emailClient.Send(emailFrom, emailTo, subject, err.Error()); err != nil {
			fmt.Println("ERROR: " + err.Error())
		}
		fmt.Println("ERROR: " + err.Error())
		return
	}
	fmt.Printf("Publishing OK (%s). Took %v\n", dpl, time.Now().Sub(publishingST))
}

func removeFiles(paths []string, emailClient *email.Client) {
	var removeErrors []string
	for _, p := range paths {
		if err := os.Remove(p); err != nil {
			removeErrors = append(removeErrors, err.Error())
		}
	}
	if len(removeErrors) > 0 {
		joinedErrors := strings.Join(removeErrors, "\n")
		if err := emailClient.Send(emailFrom, emailTo, subject, joinedErrors); err != nil {
			fmt.Println("ERROR: " + err.Error())
		}
		fmt.Println("ERROR: " + joinedErrors)
	}
}
