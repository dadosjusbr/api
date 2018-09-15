package packager

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/frictionlessdata/datapackage-go/datapackage"
)

// Pack creates a zip frictionless data package an returns its contents.
func Pack(name string, dataFileContents []byte) ([]byte, error) {
	fName := "data.csv"
	sch, err := schemaDescriptor()
	if err != nil {
		return fmt.Errorf("Error getting schema descritptor:%q", err)
	}
	d := map[string]interface{}{
		"name": name,
		"resources": []interface{}{
			map[string]interface{}{
				"name":    name,
				"path":    fName,
				"format":  "csv",
				"profile": "tabular-data-resource",
				"schema":  sch,
			},
		},
	}
	pkg, err := datapackage.New(d, ".")
	if err != nil {
		return nil, fmt.Errorf("Error creating data package:%q", err)
	}
	descriptorContents, err := json.Marshal(pkg.Descriptor())
	if err != nil {
		return nil, fmt.Errorf("Error getting descriptor contents: %q", err)
	}
	buf := new(bytes.Buffer)
	w := zip.NewWriter(buf)
	addFileToZip(w, fName, dataFileContents)
	addFileToZip(w, "datapackage.json", descriptorContents)
	err = w.Close()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func addFileToZip(w *zip.Writer, name string, contents []byte) error {
	f, err := w.Create(name)
	if err != nil {
		return err
	}
	_, err = f.Write(contents)
	if err != nil {
		return err
	}
	return nil
}
