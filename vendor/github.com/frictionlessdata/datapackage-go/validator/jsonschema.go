package validator

import (
	"github.com/santhosh-tekuri/jsonschema"
)

// jsonSchemaValidator is a validator backed by JSONSchema parsing and validation.
type jsonSchema struct {
	schema *jsonschema.Schema
}

// IsValid checks the passed-in descriptor against the JSONSchema. If it returns
// false, erros can be checked calling Errors() method.
func (v *jsonSchema) Validate(descriptor map[string]interface{}) error {
	return v.schema.ValidateInterface(descriptor)
}
