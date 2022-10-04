package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/dadosjusbr/api/models"
	"github.com/gocarina/gocsv"
)

func getRemunerationsFromZip(zipFilePath string) ([]models.SearchResult, error) {
	resp, err := http.Get(zipFilePath)
	if err != nil {
		return nil, fmt.Errorf("error requesting zip file from (%s): %w", zipFilePath, err)
	}
	defer resp.Body.Close()

	var results []models.SearchResult

	// If the status code is not 200, return an error.
	if resp.StatusCode == http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		zipReader, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
		if err != nil {
			return nil, fmt.Errorf("error opening zip reader(%s): %w", zipFilePath, err)
		}

		// All CSV zips have one file inside called "remuneracoes.csv".
		fReader, err := zipReader.File[0].Open()
		if err != nil {
			return nil, fmt.Errorf("error opening zip file(%s): %w", zipFilePath, err)
		}
		defer fReader.Close()

		if err := gocsv.Unmarshal(fReader, &results); err != nil {
			return nil, fmt.Errorf("error unmarshaling remuneracoes.csv: %w", err)
		}
		return results, nil
	}
	return nil, fmt.Errorf("error requesting zip file from (%s): %w", zipFilePath, err)
}
