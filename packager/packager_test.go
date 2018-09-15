package packager

import (
	"testing"
)

func TestInit(t *testing.T) {
	if _, ok := schemaDescriptor["fields"]; !ok {
		t.Errorf("schema expected to have fields.")
	}
}
