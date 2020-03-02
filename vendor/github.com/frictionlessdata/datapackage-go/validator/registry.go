package validator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/frictionlessdata/datapackage-go/validator/profile_cache"
	"github.com/santhosh-tekuri/jsonschema"

	_ "github.com/santhosh-tekuri/jsonschema/httploader" // This import alows jsonschema to load urls.
	_ "github.com/santhosh-tekuri/jsonschema/loader"     // This import alows jsonschema to load filepaths.
)

// RegistryLoader loads a registry.
type RegistryLoader func() (Registry, error)

// Registry represents a set of registered validators, which could be loaded locally or remotelly.
type Registry interface {
	GetValidator(profile string) (DescriptorValidator, error)
}

type profileSpec struct {
	ID            string `json:"id,omitempty"`
	Title         string `json:"title,omitempty"`
	Schema        string `json:"schema,omitempty"`
	SchemaPath    string `json:"schema_path,omitempty"`
	Specification string `json:"specification,omitempty"`
}

type localRegistry struct {
	registry     map[string]profileSpec
	inMemoryOnly bool
}

func (local *localRegistry) GetValidator(profile string) (DescriptorValidator, error) {
	spec, ok := local.registry[profile]
	if !ok {
		return nil, fmt.Errorf("Invalid profile:%s", profile)
	}
	b, err := profile_cache.FSByte(!local.inMemoryOnly, spec.Schema)
	if err != nil {
		return nil, err
	}
	c := jsonschema.NewCompiler()
	c.AddResource(profile, bytes.NewReader(b)) // Adding in-memory resource.
	schema, err := c.Compile(profile)
	if err != nil {
		return nil, err
	}
	return &jsonSchema{schema: schema}, nil
}

// LocalRegistryLoader creates a new registry, which is based on the local file system (or in-memory cache)
// to locate json schema profiles. Setting inMemoryOnly to true will make sure only the in-memory
// cache (registry_cache Go package) is accessed, thus avoiding access the filesystem.
func LocalRegistryLoader(localRegistryPath string, inMemoryOnly bool) RegistryLoader {
	return func() (Registry, error) {
		buf, err := profile_cache.FSByte(!inMemoryOnly, localRegistryPath)
		if err != nil {
			return nil, err
		}
		m, err := unmarshalRegistryContents(buf)
		if err != nil {
			return nil, err
		}
		return &localRegistry{registry: m, inMemoryOnly: inMemoryOnly}, nil
	}
}

type remoteRegistry struct {
	registry map[string]profileSpec
}

func (remote *remoteRegistry) GetValidator(profile string) (DescriptorValidator, error) {
	spec, ok := remote.registry[profile]
	if !ok {
		return nil, fmt.Errorf("Invalid profile:%s", profile)
	}
	c := jsonschema.NewCompiler()
	c.AddResource(profile, strings.NewReader(spec.Schema)) // Adding in-memory resource.
	schema, err := c.Compile(spec.Schema)
	if err != nil {
		return nil, err
	}
	return &jsonSchema{schema: schema}, nil
}

// RemoteRegistryLoader loads the schema registry map from the passed-in URL.
func RemoteRegistryLoader(url string) RegistryLoader {
	return func() (Registry, error) {
		resp, err := http.Get(url)
		if err != nil {
			return nil, fmt.Errorf("error fetching remote profile cache registry from %s. Err:%q\n", url, err)
		}
		defer resp.Body.Close()
		buf, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("error reading remote profile cache registry from %s. Err:%q\n", url, err)
		}
		m, err := unmarshalRegistryContents(buf)
		if err != nil {
			return nil, err
		}
		return &remoteRegistry{registry: m}, nil
	}
}

// FallbackRegistryLoader returns the first passed-in registry loaded successfully.
// It returns an error if there is no successfully loaded registry.
func FallbackRegistryLoader(loaders ...RegistryLoader) RegistryLoader {
	return func() (Registry, error) {
		if len(loaders) == 0 {
			return nil, fmt.Errorf("there should be at least one registry loader to fallback")
		}
		var registry Registry
		var errors []error
		for _, loader := range loaders {
			reg, err := loader()
			if err != nil {
				errors = append(errors, err)
				continue
			}
			registry = reg
			break
		}
		if registry == nil {
			var erroMsg string
			for _, err := range errors {
				erroMsg += fmt.Sprintln(err.Error())
			}
			return nil, fmt.Errorf(erroMsg)
		}
		return registry, nil
	}
}

func unmarshalRegistryContents(buf []byte) (map[string]profileSpec, error) {
	var specs []profileSpec
	if err := json.Unmarshal(buf, &specs); err != nil {
		return nil, fmt.Errorf("error parsing profile cache registry. Contents:\"%s\". Err:\"%q\"\n", string(buf), err)
	}
	m := make(map[string]profileSpec, len(specs))
	for _, s := range specs {
		m[s.ID] = s
	}
	return m, nil
}
