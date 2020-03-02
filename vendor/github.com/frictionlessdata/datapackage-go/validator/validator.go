package validator

import (
	"fmt"
	"strings"

	"github.com/santhosh-tekuri/jsonschema"
)

// DescriptorValidator validates a Data-Package or Resource descriptor.
type DescriptorValidator interface {
	Validate(map[string]interface{}) error
}

const localRegistryPath = "/registry.json"
const remoteRegistryURL = "http://frictionlessdata.io/schemas/registry.json"

// NewRegistry returns a registry where users could get descriptor validators.
func NewRegistry(loaders ...RegistryLoader) (Registry, error) {
	// Default settings.
	if len(loaders) == 0 {
		loaders = append(
			loaders,
			InMemoryLoader(),
			LocalRegistryLoader(localRegistryPath, false /* inMemoryOnly*/),
			RemoteRegistryLoader(remoteRegistryURL))
	}
	registry, err := FallbackRegistryLoader(loaders...)()
	if err != nil {
		return nil, fmt.Errorf("could not load registry:%q", err)
	}
	return registry, nil
}

// New returns a new descriptor validator for the passed-in profile.
func New(profile string, loaders ...RegistryLoader) (DescriptorValidator, error) {
	// If it is a third-party schema. Directly referenced from the internet or local file.
	if strings.HasPrefix(profile, "http") || strings.HasPrefix(profile, "file") {
		schema, err := jsonschema.Compile(profile)
		if err != nil {
			return nil, err
		}
		return &jsonSchema{schema: schema}, nil
	}
	registry, err := NewRegistry(loaders...)
	if err != nil {
		return nil, err
	}
	return registry.GetValidator(profile)
}

// Validate checks whether the descriptor the descriptor is valid against the passed-in profile/registry.
// If the validation process generates multiple errors, their messages are coalesced.
// It is a syntax-sugar around getting the validator from the registry and coalescing errors.
func Validate(descriptor map[string]interface{}, profile string, registry Registry) error {
	validator, err := registry.GetValidator(profile)
	if err != nil {
		return fmt.Errorf("Invalid Schema (Profile:%s):%q", profile, err)
	}
	return validator.Validate(descriptor)
}

// MustInMemoryRegistry returns the local cache registry, which is shipped with the library.
// It panics if there are errors retrieving the registry.
func MustInMemoryRegistry() Registry {
	reg, err := InMemoryLoader()()
	if err != nil {
		panic(err)
	}
	return reg
}

// InMemoryLoader returns a loader which points tothe local cache registry.
func InMemoryLoader() RegistryLoader {
	return LocalRegistryLoader(localRegistryPath, true /* in memory only*/)
}
