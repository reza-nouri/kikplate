package plategenerator

type PlateYAML struct {
	Name    string                 `yaml:"name"`
	Schema  map[string]SchemaField `yaml:"schema"`
	Modules map[string]ModuleDef   `yaml:"modules"`
	Files   []FileEntry            `yaml:"files"`
}

type SchemaField struct {
	Type     string   `yaml:"type"`
	Required bool     `yaml:"required"`
	Values   []string `yaml:"values"`
	Default  any      `yaml:"default"`
}

type ModuleDef struct {
	Enabled bool `yaml:"enabled"`
}

type FileEntry struct {
	Path      string `yaml:"path"`
	Template  string `yaml:"template"`
	Condition string `yaml:"condition"`
}
