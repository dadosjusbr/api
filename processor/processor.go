package processor

import (
	"fmt"
	"github.com/dadosjusbr/remuneracao-magistrados/email"
	"github.com/dadosjusbr/remuneracao-magistrados/store"
)

func Process(emailClient *email.Client, pcloudClient *store.PCloudClient, month, year int) {
	fmt.Printf("Email client:%s\n", emailClient)
	fmt.Printf("PCloud client:%s\n", pcloudClient)
	fmt.Printf("Month:%d Year: %d\n", month, year)
}