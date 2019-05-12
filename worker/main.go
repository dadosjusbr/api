package main

import (
	"fmt"
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

const remuneracaoPath = "http://www.cnj.jus.br/transparencia/remuneracao-dos-magistrados/remuneracao-"

var months = map[int]string{
	1:  "janeiro",
	2:  "fevereiro",
	3:  "marco",
	4:  "abril",
	5:  "maio",
	6:  "junho",
	7:  "julho",
	8:  "agosto",
	9:  "setembro",
	10: "outubro",
	11: "novembro",
	12: "dezembro",
}

const (
	month = 04
	year  = 2018
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
	parserClient := parser.NewServiceClient(conf.ParserURL)

	processor.Process(fmt.Sprintf("%s%s-%d", remuneracaoPath, months[month], year), fmt.Sprintf("%d-%d", month, year), emailClient, pcloudClient, parserClient)
}
