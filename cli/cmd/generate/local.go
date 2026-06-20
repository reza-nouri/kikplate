package generate

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

func loadPlateYAML(templateDir string) (*plateYAML, error) {
	raw, err := os.ReadFile(filepath.Join(templateDir, "plate.yaml"))
	if err != nil {
		return nil, fmt.Errorf("cannot read plate.yaml in %s: %w", templateDir, err)
	}
	var py plateYAML
	if err := yaml.Unmarshal(raw, &py); err != nil {
		return nil, fmt.Errorf("cannot parse plate.yaml: %w", err)
	}
	return &py, nil
}

func generateLocal(templateDir string, py *plateYAML, values map[string]any) ([]byte, error) {
	if err := applyDefaults(py, values); err != nil {
		return nil, err
	}
	data := buildTemplateData(py, values)
	return buildZip(py, templateDir, data)
}

func resolveTemplateContent(templateDir, tmpl string) (string, error) {
	if strings.HasPrefix(tmpl, "https://") || strings.HasPrefix(tmpl, "http://") {
		var lastErr error
		for attempt := 0; attempt < 3; attempt++ {
			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, tmpl, nil)
			if err != nil {
				cancel()
				return "", fmt.Errorf("cannot build request for remote template %s: %w", tmpl, err)
			}

			resp, err := httpClient.Do(req)
			if err == nil && resp != nil && resp.StatusCode == http.StatusOK {
				raw, readErr := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
				resp.Body.Close()
				cancel()
				if readErr != nil {
					return "", readErr
				}
				return string(raw), nil
			}

			if err != nil {
				lastErr = err
			} else {
				lastErr = fmt.Errorf("status %d", resp.StatusCode)
				resp.Body.Close()
			}
			cancel()

			if attempt < 2 {
				time.Sleep(time.Duration(150*(attempt+1)) * time.Millisecond)
			}
		}
		return "", fmt.Errorf("cannot fetch remote template %s: %v", tmpl, lastErr)
	}

	if !strings.Contains(tmpl, "\n") {
		candidate := filepath.Join(templateDir, tmpl)
		if _, err := os.Stat(candidate); err == nil {
			raw, err := os.ReadFile(candidate)
			if err != nil {
				return "", err
			}
			return string(raw), nil
		}
	}

	return tmpl, nil
}
