package clone

import (
	"bytes"
	"encoding/gob"
)

func init() {
	gob.Register(map[string]interface{}{}) // descriptor.
	gob.Register([]interface{}{})          // data-package resources.
}

// Descriptor deep-copies the passed-in descriptor and returns its copy.
func Descriptor(d map[string]interface{}) (map[string]interface{}, error) {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(d); err != nil {
		return nil, err
	}
	var c map[string]interface{}
	if err := gob.NewDecoder(&buf).Decode(&c); err != nil {
		return nil, err
	}
	return c, nil
}
