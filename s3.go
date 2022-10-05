package main

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/dadosjusbr/api/models"
	"github.com/gocarina/gocsv"
)

func download(ctx context.Context, sess *session.Session, bucket string, objects []string) ([]models.SearchResult, error) {
	forDownload := []s3manager.BatchDownloadObject{}
	var buffer []aws.WriteAtBuffer
	for i, key := range objects {
		buffer = append(buffer, aws.WriteAtBuffer{})
		// Create batch download objects and add them to the forDownload list
		forDownload = append(forDownload, s3manager.BatchDownloadObject{
			Object: &s3.GetObjectInput{
				Bucket: aws.String(bucket),
				Key:    aws.String(key),
			},
			Writer: &buffer[i],
		})
	}

	// Run the iterator
	err := s3manager.NewDownloader(sess).DownloadWithIterator(ctx, &s3manager.DownloadObjectsIterator{Objects: forDownload})
	if err != nil {
		log.Fatal()
	}
	var results []models.SearchResult

	for _, downloadObject := range forDownload {
		buf, ok := downloadObject.Writer.(*aws.WriteAtBuffer)
		if !ok {
			continue
		}
		zipReader, err := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(len(buf.Bytes())))
		if err != nil {
			return nil, fmt.Errorf("error opening zip reader: %w", err)
		}

		fReader, err := zipReader.File[0].Open()
		if err != nil {
			return nil, fmt.Errorf("error opening zip file: %w", err)
		}
		defer fReader.Close()

		var r []models.SearchResult
		if err := gocsv.Unmarshal(fReader, &r); err != nil {
			return nil, fmt.Errorf("error unmarshaling remuneracoes.csv: %w", err)
		}
		results = append(results, r...)
	}

	return results, err
}
