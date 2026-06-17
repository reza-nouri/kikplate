package plategenerator

import "fmt"

func ApplyDefaults(py *PlateYAML, values map[string]any) {
	for key, field := range py.Schema {
		if _, ok := values[key]; !ok && field.Default != nil {
			values[key] = field.Default
		}
	}
}

func ValidateAndCoerce(py *PlateYAML, values map[string]any) error {
	for key, field := range py.Schema {
		val, ok := values[key]
		if !ok || IsEmptyValue(val) {
			if field.Required {
				return fmt.Errorf("missing required field %q", key)
			}
			continue
		}

		coerced, err := CoerceByType(key, field.Type, val)
		if err != nil {
			return err
		}
		values[key] = coerced

		if field.Type == "enum" && len(field.Values) > 0 {
			strVal := fmt.Sprintf("%v", coerced)
			found := false
			for _, allowed := range field.Values {
				if allowed == strVal {
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("%q must be one of %v", key, field.Values)
			}
		}
	}
	return nil
}

func BuildTemplateData(py *PlateYAML, values map[string]any) map[string]any {
	data := make(map[string]any, len(values)+1)
	for k, v := range values {
		data[k] = v
	}

	modules := make(map[string]any, len(py.Modules))
	for name, mod := range py.Modules {
		enabled := mod.Enabled
		if v, ok := values["modules."+name+".enabled"]; ok {
			if b, ok := AsBool(v); ok {
				enabled = b
			}
		}
		if nestedModules, ok := values["modules"].(map[string]any); ok {
			if rawModule, ok := nestedModules[name].(map[string]any); ok {
				if b, ok := AsBool(rawModule["enabled"]); ok {
					enabled = b
				}
			}
		}
		modules[name] = map[string]any{"enabled": enabled}
	}
	data["modules"] = modules
	return data
}
