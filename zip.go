package main

import (
	"archive/zip"
	"fmt"

	"github.com/dadosjusbr/api/models"
	"github.com/gocarina/gocsv"
)

func getRemunerationsFromZip(zipFilePath string) ([]models.SearchResult, error) {
	zipReader, err := zip.OpenReader(zipFilePath)
	if err != nil {
		return nil, fmt.Errorf("error opening zip reader(%s): %w", zipFilePath, err)
	}
	defer zipReader.Close()

	var results []models.SearchResult

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
