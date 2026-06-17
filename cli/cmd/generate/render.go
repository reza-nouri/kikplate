package generate

import (
	"fmt"

	gen "github.com/kickplate/plategenerator"
)

func applyDefaults(py *plateYAML, values map[string]any) error {
	shared := toShared(py)
	gen.ApplyDefaults(shared, values)
	if err := gen.ValidateAndCoerce(shared, values); err != nil {
		return fmt.Errorf("%s — pass it with --set key=value or -f values.yaml", err.Error())
	}
	return nil
}

func buildTemplateData(py *plateYAML, values map[string]any) map[string]any {
	return gen.BuildTemplateData(toShared(py), values)
}

func toShared(py *plateYAML) *gen.PlateYAML {
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
