package generate

type plateYAML struct {
	Name    string                 `yaml:"name"`
	Schema  map[string]schemaField `yaml:"schema"`
	Modules map[string]moduleDef   `yaml:"modules"`
	Files   []fileEntry            `yaml:"files"`
}

type schemaField struct {
	Type     string   `yaml:"type"`
	Required bool     `yaml:"required"`
	Values   []string `yaml:"values"`
	Default  any      `yaml:"default"`
}

type moduleDef struct {
	Enabled bool `yaml:"enabled"`
}

type fileEntry struct {
	Path      string `yaml:"path"`
	Template  string `yaml:"template"`
	Condition string `yaml:"condition"`
}
