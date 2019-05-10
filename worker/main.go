package main

import (
	"log"

	"github.com/dadosjusbr/remuneracao-magistrados/email"
	"github.com/dadosjusbr/remuneracao-magistrados/parser"
	"github.com/dadosjusbr/remuneracao-magistrados/processor"
	"github.com/dadosjusbr/remuneracao-magistrados/store"
	"github.com/kelseyhightower/envconfig"
)

type config struct {
	SendgridAPIKey string `envconfig:"SENDGRID_API_KEY"`
	PCloudUsername string `envconfig:"PCLOUD_USERNAME"`
	PCloudPassword string `envconfig:"PCLOUD_PASSWORD"`
	ParserURL      string `envconfig:"PARSER_URL"`
}

const (
	emailFrom = "no-reply@dadosjusbr.com"
	emailTo   = "dadosjusbrops@googlegroups.com"
	subject   = "remuneracao-magistrados error"
)

func main() {
	// TODO: Treat Signals.
	var conf config
	err := envconfig.Process("remuneracao-magistrados", &conf)
	if err != nil {
		log.Fatal(err.Error())
	}
	emailClient, err := email.NewClient(conf.SendgridAPIKey)
	if err != nil {
		log.Fatal("ERROR: ", err.Error())
	}
	pcloudClient, err := store.NewPCloudClient(conf.PCloudUsername, conf.PCloudPassword)
	if err != nil {
		log.Fatal("ERROR: ", err.Error())
	}
	processor.Process(04, 2018, emailClient, pcloudClient, parser.NewServiceClient(conf.ParserURL))
}
