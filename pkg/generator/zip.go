package plategenerator

import (
	"archive/zip"
	"bytes"
	"fmt"
	"strings"
)

type TemplateResolver func(ref string) (string, error)

func BuildZip(py *PlateYAML, data map[string]any, resolver TemplateResolver) ([]byte, error) {
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)

	for _, fe := range py.Files {
		if fe.Condition != "" && !EvalCondition(fe.Condition, data) {
			continue
		}

		tmplText, err := resolver(fe.Template)
		if err != nil {
			return nil, fmt.Errorf("cannot load template for %s: %w", fe.Path, err)
		}

		renderedPath, err := RenderTemplate("path:"+fe.Path, fe.Path, data)
		if err != nil {
			return nil, fmt.Errorf("cannot render path %s: %w", fe.Path, err)
		}

		renderedPath, err = SanitizePath(renderedPath)
		if err != nil {
			return nil, fmt.Errorf("invalid rendered path %s: %w", fe.Path, err)
		}

		rendered, err := RenderTemplate(fe.Path, tmplText, data)
		if err != nil {
			return nil, fmt.Errorf("cannot render %s: %w", fe.Path, err)
		}

		f, err := zw.Create(renderedPath)
		if err != nil {
			return nil, fmt.Errorf("cannot create zip entry %s: %w", renderedPath, err)
		}
		if _, err := f.Write([]byte(rendered)); err != nil {
			return nil, fmt.Errorf("cannot write zip entry %s: %w", renderedPath, err)
		}
	}

	if err := zw.Close(); err != nil {
		return nil, fmt.Errorf("cannot finalise zip: %w", err)
	}
	return buf.Bytes(), nil
}

func SanitizePath(p string) (string, error) {
	p = strings.TrimSpace(strings.ReplaceAll(p, "\\", "/"))
	if p == "" {
		return "", fmt.Errorf("path cannot be empty")
	}
	if strings.HasPrefix(p, "/") {
		return "", fmt.Errorf("absolute paths are not allowed")
	}
	clean := cleanPath(p)
	if clean == "." || clean == "" {
		return "", fmt.Errorf("path cannot be current directory")
	}
	if strings.HasPrefix(clean, "../") || clean == ".." {
		return "", fmt.Errorf("path traversal is not allowed")
	}
	return clean, nil
}

func cleanPath(p string) string {
	parts := strings.Split(p, "/")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		if part == "" || part == "." {
			continue
		}
		if part == ".." {
			if len(out) > 0 {
				out = out[:len(out)-1]
			} else {
				return ".."
			}
			continue
		}
		out = append(out, part)
	}
	if len(out) == 0 {
		return "."
	}
	return strings.Join(out, "/")
}
