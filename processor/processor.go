package processor

import (
	"bytes"
	"fmt"
	"log"
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

// Process download, parse, save and publish data of one month.
func Process(url string, year, month int, emailClient *email.Client, pcloudClient *store.PCloudClient, parser *parser.ServiceClient) {
	//TODO: this function shuld return an error if something goes wrong.
	// Download files from CNJ.
	crawST := time.Now()
	results, err := crawler.Crawl(url)
	if err != nil {
		if err := emailClient.Send(emailFrom, emailTo, subject, err.Error()); err != nil {
			fmt.Println("ERROR: " + err.Error())
		}
		fmt.Println("CRAWLING ERROR: " + err.Error())
		return
	}
	fmt.Printf("Crawling OK (%d files). Took %v\n", len(results), time.Now().Sub(crawST))

	// Parsing.
	parsingST := time.Now()
	var sContents [][]byte
	var sNames []string
	for _, r := range results {
		sContents = append(sContents, r.Body)
		sNames = append(sNames, r.Name)
	}
	csv, schema, err := parser.Parse(sContents)
	if err != nil {
		// TODO: Send an email.
		log.Fatal(err)
	}
	fmt.Printf("Parsing OK. Took %v\n", time.Now().Sub(parsingST))

	// Backup.
	backupST := time.Now()
	rl, err := pcloudClient.PutZip("2018-04-raw.zip", sNames, sContents)
	if err != nil {
		if err := emailClient.Send(emailFrom, emailTo, subject, err.Error()); err != nil {
			fmt.Println("ERROR: " + err.Error())
		}
		fmt.Println("ERROR: " + err.Error())
		return
	}
	fmt.Printf("Spreadsheets backed up OK (%s). Took %v\n", rl, time.Now().Sub(backupST))

	// Packaging.
	packagingST := time.Now()
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
