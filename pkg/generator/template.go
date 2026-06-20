package plategenerator

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"text/template"
)

var slugifyRegexp = regexp.MustCompile(`[^a-z0-9]+`)

func RenderTemplate(name, tmplText string, data map[string]any) (string, error) {
	t, err := template.New(name).Funcs(template.FuncMap{
		"lower": strings.ToLower,
		"upper": strings.ToUpper,
		"trim":  strings.TrimSpace,
		"replace": func(s, old, new string) string {
			return strings.ReplaceAll(s, old, new)
		},
		"default": func(def, val any) any {
			if IsEmptyValue(val) {
				return def
			}
			return val
		},
		"slugify":    Slugify,
		"className":  ToCamelCase,
		"camelCase":  ToCamelCase,
		"pascalCase": ToPascalCase,
		"message": func(args ...any) string {
			var result strings.Builder
			for _, arg := range args {
				result.WriteString(toString(arg))
			}
			return result.String()
		},
		"join": func(sep string, vals ...any) string {
			strs := make([]string, len(vals))
			for i, v := range vals {
				strs[i] = toString(v)
			}
			return strings.Join(strs, sep)
		},
		"has": func(s, substr string) bool {
			return strings.Contains(s, substr)
		},
	}).Option("missingkey=error").Parse(tmplText)
	if err != nil {
		return "", err
	}
	var out bytes.Buffer
	if err := t.Execute(&out, data); err != nil {
		return "", err
	}
	return out.String(), nil
}

func Slugify(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	s = slugifyRegexp.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	if s == "" {
		return "app"
	}
	return s
}

func ToCamelCase(s string) string {
	words := splitWords(s)
	if len(words) == 0 {
		return ""
	}
	result := strings.ToLower(words[0])
	for _, word := range words[1:] {
		if word != "" {
			result += strings.ToUpper(word[:1]) + strings.ToLower(word[1:])
		}
	}
	return result
}

func ToPascalCase(s string) string {
	words := splitWords(s)
	var result string
	for _, word := range words {
		if word != "" {
			result += strings.ToUpper(word[:1]) + strings.ToLower(word[1:])
		}
	}
	return result
}

func splitWords(s string) []string {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, "-", " ")
	s = strings.ReplaceAll(s, "_", " ")
	words := strings.Fields(s)
	return words
}

func toString(v any) string {
	if v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return fmt.Sprint(v)
}
