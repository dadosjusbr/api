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

func getMonthPreviewData(dtpackageZip []byte, resource string) error {
	dir, _ := ioutil.TempDir("", "dadosjusbr_temp_dir")
	zipReader, err := zip.NewReader(bytes.NewReader(dtpackageZip), int64(len(dtpackageZip)))
	if err != nil {
		return err
	}

	// Read all the files from zip archive
	for _, zipFile := range zipReader.File {
		fmt.Println("Reading file:", zipFile.Name)
		fileName, file, err := readZipFile(zipFile)
		if err != nil {
			return err
		}
		path := filepath.Join(dir, fileName)
		ioutil.WriteFile(path, file, 0666)
	}
	descriptorPath := filepath.Join(dir, "datapackage.json")
	defer os.RemoveAll(dir)
	pkg, _ := datapackage.Load(descriptorPath, validator.InMemoryLoader())
	res := pkg.GetResource(resource)
	/**
	people := []struct {
		Name string `tableheader:"nome"`
	}{}
	res.Cast(&people)
	fmt.Printf("%+v", people)
	**/
	var diarias []float64
	res.CastColumn("diarias", &diarias, csv.LoadHeaders())
	sum(diarias, "diarias")

	var auxAlimentacao []float64
	res.CastColumn("auxilio_alimentacao", &auxAlimentacao, csv.LoadHeaders())
	sum(auxAlimentacao, "auxilio alimentacao")

	var auxSaude []float64
	res.CastColumn("auxilio_saude", &auxSaude, csv.LoadHeaders())
	sum(auxSaude, "auxilio saude")

	var auxMoradia []float64
	res.CastColumn("auxilio_moradia", &auxMoradia, csv.LoadHeaders())
	sum(auxMoradia, "auxilio moradia")

	var auxPreEscolar []float64
	res.CastColumn("auxilio_pre_escolar", &auxPreEscolar, csv.LoadHeaders())
	sum(auxPreEscolar, "auxilio pre escolar")

	var ajudaDeCusto []float64
	res.CastColumn("ajuda_de_custo", &ajudaDeCusto, csv.LoadHeaders())
	sum(ajudaDeCusto, "ajuda de custo")

	return nil
}

func sum(arr []float64, label string) {
	total := 0.0
	for _, value := range arr {
		total = total + value
	}
	fmt.Printf("Total %s: %.2f\n", label, total)
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

	getMonthPreviewData(dtPackage, fmt.Sprintf("%s-datapackage", filePre))

	// Publishing.
	publishingST := time.Now()
	dpl, err := pcloudClient.Put(fmt.Sprintf("%s-datapackage.zip", filePre), bytes.NewReader(dtPackage))
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
