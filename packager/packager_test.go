package packager

import (
	"testing"
)

func TestSchemaDescriptor(t *testing.T) {
	sch, err := schemaDescriptor()
	if err != nil {
		t.Errorf("want:nil got:%q", err)
	}
	if _, ok := sch["fields"]; !ok {
		t.Errorf("schema expected to have fields.")
	}
}
