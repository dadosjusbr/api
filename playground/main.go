package main

import (
	"fmt"

	"github.com/dadosjusbr/remuneracao-magistrados/parser"
)

func main() {
	body, err := parser.GetSchema()
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Printf(string(body))
}
