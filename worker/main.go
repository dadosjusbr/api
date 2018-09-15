package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"github.com/dadosjusbr/remuneracao-magistrados/email"
	"github.com/dadosjusbr/crawler"
	"github.com/kelseyhightower/envconfig"
)

type config struct {
	SendgridAPIKey string `envconfig:"SENDGRID_API_KEY"`
}

const(
	emailFrom = "no-reply@dadosjusbr.com"
	emailTo   = "dadosjusbr.ops"
	subject   = "remuneracao-magistrados error"
)

func main() {
	var conf config
	err := envconfig.Process("remuneracao-magistrados", &conf)
	if err != nil {
		log.Fatal(err.Error())
	}
	emailClient, err := email.NewClient(conf.SendgridAPIKey)
	if(err != nil) {
		fmt.Println("ERROR: ", err.Error())
	}

	// Download files from CNJ.
	paths, err := crawler.Download(04, 2018)
	if err != nil {
		if err := emailClient.Send(emailFrom, emailTo, subject, err.Error()); err != nil {
			fmt.Println("ERROR: " + err.Error())
		}
		fmt.Println("ERROR: " + err.Error())
		return
	}

	if len(paths) == 0 {
		fmt.Println("No files to download.")
		return
	}

	// Removing downloaded files.
	var removeErrors []string
	for _, p := range paths {
		if err := os.Remove(p); err != nil {
			removeErrors = append(removeErrors, err.Error())
		}
	}
	if len(removeErrors) > 0 {
		joinedErrors := strings.Join(removeErrors, "\n")
		if err := emailClient.Send(emailFrom, emailTo, subject, joinedErrors);err != nil {
			fmt.Println("ERROR: " + err.Error())
		}
		fmt.Println("ERROR: " + joinedErrors)
	}
}
