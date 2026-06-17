package generator

import (
	"fmt"
	"strings"

	gen "github.com/kickplate/plategenerator"
)

func renderProject(py *PlateYAML, values map[string]any) ([]byte, error) {
	shared := toShared(py)

	gen.ApplyDefaults(shared, values)

	if err := gen.ValidateAndCoerce(shared, values); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrInvalidInput, err.Error())
	}

	data := gen.BuildTemplateData(shared, values)

	return gen.BuildZip(shared, data, func(ref string) (string, error) {
		if strings.HasPrefix(ref, "https://") || strings.HasPrefix(ref, "http://") {
			content, err := fetchRemoteTemplate(ref)
			if err != nil {
				return "", fmt.Errorf("%w: %v", ErrFetchFailed, err)
			}
			return content, nil
		}
		return ref, nil
	})
}

func toShared(py *PlateYAML) *gen.PlateYAML {
	schema := make(map[string]gen.SchemaField, len(py.Schema))
	for k, f := range py.Schema {
		schema[k] = gen.SchemaField{
			Type:     f.Type,
			Required: f.Required,
			Values:   f.Values,
			Default:  f.Default,
		}
	}

	modules := make(map[string]gen.ModuleDef, len(py.Modules))
	for k, m := range py.Modules {
		modules[k] = gen.ModuleDef{Enabled: m.Enabled}
	}

	files := make([]gen.FileEntry, len(py.Files))
	for i, f := range py.Files {
		files[i] = gen.FileEntry{
			Path:      f.Path,
			Template:  f.Template,
			Condition: f.Condition,
		}
	}

	return &gen.PlateYAML{
		Name:    py.Name,
		Schema:  schema,
		Modules: modules,
		Files:   files,
	}
}
