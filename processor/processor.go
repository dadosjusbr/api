package processor

import (
	"bytes"
	"fmt"
	"time"

	"github.com/dadosjusbr/remuneracao-magistrados/crawler"
	"github.com/dadosjusbr/remuneracao-magistrados/db"
	"github.com/dadosjusbr/remuneracao-magistrados/packager"
	"github.com/dadosjusbr/remuneracao-magistrados/parser"
	"github.com/dadosjusbr/remuneracao-magistrados/store"
)

// Process download, parse, save and publish data of one month.
func Process(url string, month, year int, pcloudClient *store.PCloudClient, parser *parser.ServiceClient, dbClient *db.Client) error {
	//TODO: this function shuld return an error if something goes wrong.
	// Download files from CNJ.
	crawST := time.Now()
	results, err := crawler.Crawl(url)
	if err != nil {
		fmt.Println("CRAWLING ERROR: " + err.Error())
		return err
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
	csv, schema, err := parser.Parse(sContents, sNames)
	if err != nil {
		fmt.Println("PARSING ERROR: " + err.Error())
		return err
	}
	fmt.Printf("Parsing OK. Took %v\n", time.Now().Sub(parsingST))

	filePre := fmt.Sprintf("%d-%d", year, month)

	// Backup.
	backupST := time.Now()
	rl, err := pcloudClient.PutZip(fmt.Sprintf("%s-raw.zip", filePre), sNames, sContents)
	if err != nil {
		fmt.Println("BACKUP ERROR: " + err.Error())
		return err
	}
	fmt.Printf("Spreadsheets backed up OK (%s). Took %v\n", rl, time.Now().Sub(backupST))

	// Packaging.
	packagingST := time.Now()
	datapackage, err := packager.Pack(fmt.Sprintf("%s-datapackage", filePre), schema, csv)
	if err != nil {
		fmt.Println("PACKAGING ERROR: " + err.Error())
		return err
	}
	fmt.Printf("Packaging OK. Took: %s\n", time.Now().Sub(packagingST))

	// Publishing.
	publishingST := time.Now()
	dpl, err := pcloudClient.Put(fmt.Sprintf("%s-datapackage.zip", filePre), bytes.NewReader(datapackage))
	if err != nil {
		fmt.Println("PUBLISHING ERROR: " + err.Error())
		return err
	}
	fmt.Printf("Publishing OK (%s). Took %v\n", dpl, time.Now().Sub(publishingST))

	mr := db.MonthResults{Month: month, Year: year, SpreadsheetsURL: rl, DatapackageURL: dpl, Success: true}
	err = dbClient.SaveMonthResults(mr)
	if err != nil {
		return err
	}

	return nil
}
