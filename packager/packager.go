package packager

import (
	"fmt"

	"github.com/frictionlessdata/datapackage-go/datapackage"
)

// Pack creates a
func Pack(name string, content []byte) error {
	sch, err := schemaDescriptor()
	if err != nil {
		return fmt.Errorf("Error getting schema descritptor:%q", err)
	}
	fName := "data.csv"
	d := map[string]interface{}{
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
	_ := datapackage.New(d, ".")
	// TODO:
	// Create zip datapackage in memory
	// Invoke the packager service.
	return nil
}
