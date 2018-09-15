package main

import (
	"fmt"
	"log"
	"os"

	"github.com/dadosjusbr/crawler"
)

func main() {
	// Download files from CNJ.
	paths, err := crawler.Download(04, 2018)
	if err != nil {
		// TODO(janderson): Send email to ops+dadosjusbr@gmail.com
		log.Fatal(err)
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
		fmt.Println(removeErrors)
		// TODO(janderson): Send mail to ops+dadosjusbr@gmail.com
	}
}
