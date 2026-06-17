package plategenerator

import (
	"bytes"
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
		"slugify": Slugify,
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
