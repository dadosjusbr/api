package processor

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dadosjusbr/crawler"
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

//
func Process(emailClient *email.Client, pcloudClient *store.PCloudClient, month, year int) {
	// Download files from CNJ.
	paths, err := crawler.Download(04, 2018)
	if err != nil {
		if err := emailClient.Send(emailFrom, emailTo, subject, err.Error()); err != nil {
			fmt.Println("ERROR: " + err.Error())
		}
		fmt.Println("ERROR: " + err.Error())
		return
	}
	defer removeFiles(paths, emailClient)
	fmt.Printf("Crawling OK. Download %d files.\n", len(paths))

	if len(paths) == 0 {
		fmt.Println("No files to download.")
		return
	}

	// Parsing.
	fmt.Println("Start parsing")
	parsingST := time.Now()
	// Create a buffer to write our archive to.
	var spreadsheetZipBuf bytes.Buffer
	spreadsheetZipWriter := zip.NewWriter(&spreadsheetZipBuf)

	var content bytes.Buffer
	for i, p := range paths {
		if i == 0 {
			// TODO: Treat middle parameter.
			data, _, err := parser.Parse(p, "planilhas", map[string]string{"formato_saida": "csv"})
			if err != nil {
				// TODO: Tweet and save error.
				return
			}
			content.Write(data)
			content.WriteRune('\n')
		} else {
			// TODO: Treat middle parameter.
			data, _, err := parser.Parse(p, "planilhas", map[string]string{"formato_saida": "headless_csv"})
			if err != nil {
				// TODO: Tweet and save error.
				fmt.Println("ERROR: " + err.Error())
				return
			}
			content.Write(data)
			if i < len(paths)-1 {
				content.WriteRune('\n')
			}
		}
		fmt.Printf("File %s parsed. %d missing.\n", p, len(paths)-i-1)
		zipFile, err := spreadsheetZipWriter.Create(filepath.Base(p))
		if err != nil {
			log.Fatal(err)
		}
		c, err := ioutil.ReadFile(p)
		if err != nil {
			// TODO: send email.
			fmt.Printf("ERROR reading spreadsheet contents (%s):%q", p, err)
			return
		}
		_, err = zipFile.Write(c)
		if err != nil {
			// TODO: send email.
			log.Fatal(err)
		}
	}
	if err := spreadsheetZipWriter.Close(); err != nil {
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
	datapackage, err := packager.Pack("2018-04", content.Bytes())
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