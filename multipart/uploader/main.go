package main

import (
	"fmt"
	"log"

	"github.com/dadosjusbr/remuneracao-magistrados/parser"
)

func main() {

	const path = "/home/fireman/Downloads/abril/teste.zip"

	data, _, err := parser.Parse("/home/fireman/Downloads/abril/abril-2018-0169e346edfb424b29ddb8f38ae31744.xls")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(data))
	/*
		f, err := os.Open(path)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		params := map[string]string{"formato_saida": "csv"}

		req, err := multipart.UploadRequest(url, "planilhas", "planilha.zip", f, params)
		if err != nil {
			log.Fatal(err)
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		fZip, err := ioutil.TempFile(".", ".zip")
		if err != nil {
			log.Fatal(err)
		}
		_, err = io.Copy(fZip, resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		r, err := zip.OpenReader(fZip.Name())
		if err != nil {
			log.Fatal(err)
		}
		defer r.Close()
		for _, f := range r.File {
			rc, err := f.Open()
			if err != nil {
				log.Fatal(err)
			}
			out, err := os.Create(f.Name)
			if err != nil {
				log.Fatal(err)
			}
			_, err = io.Copy(out, rc)
			if err != nil {
				log.Fatal(err)
			}
			out.Close()
			rc.Close()
		}*/
}
