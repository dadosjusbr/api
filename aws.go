package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/dadosjusbr/api/models"
	"github.com/gocarina/gocsv"
	"github.com/newrelic/go-agent/v3/newrelic"
)

type AwsSession struct {
	sess     *session.Session
	newrelic *newrelic.Application
}

func NewAwsSession(region string) (*AwsSession, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		return nil, fmt.Errorf("error creating aws session: %w", err)
	}
	return &AwsSession{sess: sess}, nil
}

func (s AwsSession) GetRemunerationsFromS3(limit, downloadLimit int, category, bucket string, results []models.SearchDetails) ([]models.SearchResult, int, error) {
	forDownload := []s3manager.BatchDownloadObject{}
	var buffer []aws.WriteAtBuffer
	var paths []string
	var numRows = 0
	var mustUnzip = true
	for _, r := range results {
		// Forma de evitar o download de arquivos que estão fora do limite.
		if mustUnzip {
			// Pegando apenas a chave do arquivo zipado.
			object := strings.Replace(r.ZipUrl, fmt.Sprintf("https://%s.s3.amazonaws.com/", bucket), "", 1)
			paths = append(paths, object)
			buffer = append(buffer, aws.WriteAtBuffer{})
		}
		/* Aqui a gente faz um "early return" se o número de resultados for maior
		que o limite de resultados da pesquisa.
		Isso evita que a gente precise processar todos os arquivos zip. */
		switch category {
		case "outras":
			numRows += r.Outras
		case "base":
			numRows += r.Base
		case "descontos":
			numRows += r.Descontos
		default:
			numRows += r.Descontos + r.Base + r.Outras
		}
		if numRows > downloadLimit {
			mustUnzip = false
		}
	}

	for i, key := range paths {
		// Criando a lista de objetos que serão baixados do S3.
		forDownload = append(forDownload, s3manager.BatchDownloadObject{
			Object: &s3.GetObjectInput{
				Bucket: aws.String(bucket),
				Key:    aws.String(key),
			},
			Writer: &buffer[i],
		})
	}

	txn := s.newrelic.StartTransaction("aws.GetRemunerations")
	defer txn.End()
	ctx := newrelic.NewContext(aws.BackgroundContext(), txn)
	// Executando o download
	downloader := s3manager.NewDownloader(s.sess)
	err := downloader.DownloadWithIterator(ctx, &s3manager.DownloadObjectsIterator{Objects: forDownload})
	if err != nil {
		return nil, 0, fmt.Errorf("error downloading files from S3: %q", err)
	}

	var searchResults []models.SearchResult
	reachedLimit := false
	for _, downloadObject := range forDownload {
		// Queremos processar apenas os dados dentro dos limites definidos.
		if reachedLimit {
			break
		}
		buf, ok := downloadObject.Writer.(*aws.WriteAtBuffer)
		if !ok {
			return nil, 0, fmt.Errorf("error converting downloaded object (%s) to WriteAtBuffer", *downloadObject.Object.Key)
		}

		zipReader, err := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(len(buf.Bytes())))
		if err != nil {
			return nil, 0, fmt.Errorf("error creating zip reader: %w", err)
		}

		fReader, err := zipReader.File[0].Open()
		if err != nil {
			return nil, 0, fmt.Errorf("error opening zip file (%s): %w", *downloadObject.Object.Key, err)
		}
		defer fReader.Close()

		var r []models.SearchResult
		if err := gocsv.Unmarshal(fReader, &r); err != nil {
			return nil, 0, fmt.Errorf("error unmarshaling remuneracoes.csv: %w", err)
		}
		/* Queremos guardar na memória apenas os resultados da categoria que o
		usuário pediu.*/
		for _, rem := range r {
			if len(searchResults) < limit {
				if category == "" || category == rem.CategoriaContracheque {
					searchResults = append(searchResults, rem)
				}
			} else {
				reachedLimit = true
				break
			}
		}
	}
	return searchResults, numRows, err
}
