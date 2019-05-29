package processor

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/dadosjusbr/remuneracao-magistrados/crawler"
	"github.com/dadosjusbr/remuneracao-magistrados/db"
	"github.com/dadosjusbr/remuneracao-magistrados/packager"
	"github.com/dadosjusbr/remuneracao-magistrados/parser"
	"github.com/dadosjusbr/remuneracao-magistrados/store"
	"github.com/frictionlessdata/datapackage-go/datapackage"
	"github.com/frictionlessdata/datapackage-go/validator"
	"github.com/frictionlessdata/tableschema-go/csv"
)

func readZipFile(zf *zip.File) (string, []byte, error) {
	f, err := zf.Open()
	if err != nil {
		return "", nil, err
	}
	defer f.Close()
	fileName := zf.FileHeader.Name
	file, err := ioutil.ReadAll(f)
	return fileName, file, err
}

func getMonthStatistics(dtpackageZip []byte, resource string) ([]db.Statistic, error) {
	dir, _ := ioutil.TempDir("", "dadosjusbr_temp_dir")
	zipReader, err := zip.NewReader(bytes.NewReader(dtpackageZip), int64(len(dtpackageZip)))
	if err != nil {
		return nil, err
	}

	// Read all the files from zip archive
	for _, zipFile := range zipReader.File {
		fileName, file, err := readZipFile(zipFile)
		if err != nil {
			return nil, err
		}
		path := filepath.Join(dir, fileName)
		ioutil.WriteFile(path, file, 0666)
	}
	descriptorPath := filepath.Join(dir, "datapackage.json")
	defer os.RemoveAll(dir)
	pkg, _ := datapackage.Load(descriptorPath, validator.InMemoryLoader())
	res := pkg.GetResource(resource)

	diarias, err := getColumnSum("diarias", res)
	if err != nil {
		return nil, err
	}

	auxAlimentacao, err := getColumnSum("auxilio_alimentacao", res)
	if err != nil {
		return nil, err
	}

	auxSaude, err := getColumnSum("auxilio_saude", res)
	if err != nil {
		return nil, err
	}

	auxMoradia, err := getColumnSum("auxilio_moradia", res)
	if err != nil {
		return nil, err
	}

	return []db.Statistic{
		db.Statistic{Name: "Diárias", Value: diarias, Description: "Total gasto com diárias nesse mês"},
		db.Statistic{Name: "Auxílio Alimentação", Value: auxAlimentacao, Description: "Total gasto com auxílio alimentação nesse mês"},
		db.Statistic{Name: "Auxílio Saúde", Value: auxSaude, Description: "Total gasto com auxílio saúde nesse mês"},
		db.Statistic{Name: "Auxílio Moradia", Value: auxMoradia, Description: "Total gasto com auxílio moradia nesse mês"},
	}, nil
}

func getColumnSum(colName string, res *datapackage.Resource) (float64, error) {
	var arr []float64
	total := 0.0
	err := res.CastColumn(colName, &arr, csv.LoadHeaders())
	if err != nil {
		return total, err
	}
	for _, value := range arr {
		total = total + value
	}
	return total, nil
}

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
	dtPackage, err := packager.Pack(fmt.Sprintf("%s-datapackage", filePre), schema, csv)
	if err != nil {
		fmt.Println("PACKAGING ERROR: " + err.Error())
		return err
	}
	fmt.Printf("Packaging OK. Took: %s\n", time.Now().Sub(packagingST))

	// Collect statistics
	statisticsST := time.Now()
	statistics, err := getMonthStatistics(dtPackage, fmt.Sprintf("%s-datapackage", filePre))
	if err != nil {
		fmt.Println("STATISTICS ERROR: " + err.Error())
		return err
	}
	fmt.Printf("Collected statistics. Took: %s\n", time.Now().Sub(statisticsST))

	// Publishing.
	publishingST := time.Now()
	dpl, err := pcloudClient.Put(fmt.Sprintf("%s-datapackage.zip", filePre), bytes.NewReader(dtPackage))
	if err != nil {
		fmt.Println("PUBLISHING ERROR: " + err.Error())
		return err
	}
	fmt.Printf("Publishing OK (%s). Took %v\n", dpl, time.Now().Sub(publishingST))

	mr := db.MonthResults{
		Month:           month,
		Year:            year,
		SpreadsheetsURL: rl,
		DatapackageURL:  dpl,
		Success:         true,
		Statistics:      statistics,
	}
	err = dbClient.SaveMonthResults(mr)
	if err != nil {
		return err
	}

	return nil
}
