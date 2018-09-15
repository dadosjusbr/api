package packager

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

var schemaDescriptor = make(map[string]interface{})

const schemaPath = "schema.json"

func init() {
	f, err := os.Open(schemaPath)
	if err != nil {
		log.Fatalf("Error trying to open the schema descriptor: %q", err)
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatalf("Error trying to read the schema descriptor: %q", err)
	}
	if err := json.Unmarshal(b, &schemaDescriptor); err != nil {
		log.Fatalf("Error trying to unmarshal the schema descriptor: %q", err)
	}
}
