package generate

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

func promptSchema(py *plateYAML, values map[string]any) error {
	reader := bufio.NewReader(os.Stdin)

	if py.Name != "" {
		fmt.Printf("Generating: %s\n\n", py.Name)
	}

	keys := make([]string, 0, len(py.Schema))
	for k := range py.Schema {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		if _, already := values[key]; already {
			continue
		}
		field := py.Schema[key]
		fmt.Print(buildPrompt(key, field))

		line, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("input error: %w", err)
		}
		line = strings.TrimSpace(line)

		if line == "" && field.Default != nil {
			values[key] = field.Default
		} else if line != "" {
			values[key] = line
		} else if field.Required {
			return fmt.Errorf("field %q is required", key)
		}
	}

	moduleKeys := make([]string, 0, len(py.Modules))
	for k := range py.Modules {
		moduleKeys = append(moduleKeys, k)
	}
	sort.Strings(moduleKeys)

	for _, name := range moduleKeys {
		mod := py.Modules[name]
		dotKey := "modules." + name + ".enabled"
		if _, already := values[dotKey]; already {
			continue
		}
		def := "n"
		if mod.Enabled {
			def = "y"
		}
		fmt.Printf("Enable module %q? [y/n] (default: %s): ", name, def)
		line, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("input error: %w", err)
		}
		line = strings.TrimSpace(strings.ToLower(line))
		if line == "" {
			values[dotKey] = mod.Enabled
		} else {
			values[dotKey] = line == "y" || line == "yes"
		}
	}

	return nil
}

func buildPrompt(key string, field schemaField) string {
	var b strings.Builder
	b.WriteString(key)
	if field.Type != "" && field.Type != "string" {
		b.WriteString(" (")
		b.WriteString(field.Type)
		if field.Type == "enum" && len(field.Values) > 0 {
			b.WriteString(": ")
			b.WriteString(strings.Join(field.Values, "|"))
		}
		b.WriteString(")")
	}
	if field.Required {
		b.WriteString(" [required]")
	}
	if field.Default != nil {
		b.WriteString(fmt.Sprintf(" (default: %v)", field.Default))
	}
	b.WriteString(": ")
	return b.String()
}
