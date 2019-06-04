package processor

import (
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
	"github.com/montanaflynn/stats"
)

func getMonthStatistics(dtpackageZip []byte, resource string) ([]db.Statistic, error) {
	dir, err := ioutil.TempDir("", "dadosjusbr_temp_dir")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "datapackage.zip")
	err = ioutil.WriteFile(path, dtpackageZip, 0666)
	if err != nil {
		return nil, err
	}
	pkg, err := datapackage.Load(path, validator.InMemoryLoader())
	if err != nil {
		return nil, err
	}
	res := pkg.GetResource(resource)

	iter, err := res.Iter(csv.LoadHeaders())
	if err != nil {
		return nil, err
	}
	sch, err := res.GetSchema()
	if err != nil {
		return nil, err
	}

	var totalRendimentos []float64
	var subsidio []float64
	var totalAuxilios []float64

	data := struct {
		TotalRendimentos float64 `tableheader:"total_de_rendimentos"`
		Subsidio         float64 `tableheader:"subsidio"`
	}{}

	for iter.Next() {
		sch.CastRow(iter.Row(), &data)
		totalRendimentos = append(totalRendimentos, data.TotalRendimentos)
		subsidio = append(subsidio, data.Subsidio)
		totalAuxilios = append(totalAuxilios, data.TotalRendimentos-data.Subsidio)
	}

	totalRendimentosSum, err := stats.Sum(totalRendimentos)
	if err != nil {
		return nil, err
	}
	subsidioSum, err := stats.Sum(subsidio)
	if err != nil {
		return nil, err
	}
	totalAuxiliosSum, err := stats.Sum(totalAuxilios)
	if err != nil {
		return nil, err
	}

	totalRendimentosMean, err := stats.Mean(totalRendimentos)
	if err != nil {
		return nil, err
	}
	subsidioMean, err := stats.Mean(subsidio)
	if err != nil {
		return nil, err
	}
	totalAuxiliosMean, err := stats.Mean(totalAuxilios)
	if err != nil {
		return nil, err
	}

	totalRendimentosMedian, err := stats.Median(totalRendimentos)
	if err != nil {
		return nil, err
	}
	subsidioMedian, err := stats.Median(subsidio)
	if err != nil {
		return nil, err
	}
	totalAuxiliosMedian, err := stats.Median(totalAuxilios)
	if err != nil {
		return nil, err
	}

	totalRendimentosStdDev, err := stats.StandardDeviation(totalRendimentos)
	if err != nil {
		return nil, err
	}
	subsidioStdDev, err := stats.StandardDeviation(subsidio)
	if err != nil {
		return nil, err
	}
	totalAuxiliosStdDev, err := stats.StandardDeviation(totalAuxilios)
	if err != nil {
		return nil, err
	}

	return []db.Statistic{
		{
			Name:        "Subsídios",
			Description: "Salário base do magistrado",
			Sum:         subsidioSum,
			SampleSize:  len(subsidio),
			Mean:        subsidioMean,
			Median:      subsidioMedian,
			StdDev:      subsidioStdDev,
		},
		{
			Name:        "Auxílios",
			Description: "Descreve quanto o magistrado recebeu em auxílios nesse mês",
			Sum:         totalAuxiliosSum,
			SampleSize:  len(totalAuxilios),
			Mean:        totalAuxiliosMean,
			Median:      totalAuxiliosMedian,
			StdDev:      totalAuxiliosStdDev,
		},
		{
			Name:        "Total de Rendimentos",
			Description: "Soma dos salários brutos de todos os magistrados nesse mês",
			Sum:         totalRendimentosSum,
			SampleSize:  len(totalRendimentos),
			Mean:        totalRendimentosMean,
			Median:      totalRendimentosMedian,
			StdDev:      totalRendimentosStdDev,
		},
	}, nil
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
