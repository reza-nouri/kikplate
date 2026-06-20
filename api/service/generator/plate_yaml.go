package generator

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type PlateYAML struct {
	Name    string                      `yaml:"name"`
	Schema  map[string]SchemaField      `yaml:"schema"`
	Modules map[string]ModuleDefinition `yaml:"modules"`
	Files   []FileEntry                 `yaml:"files"`
}

type SchemaField struct {
	Type     string   `yaml:"type"`
	Required bool     `yaml:"required"`
	Values   []string `yaml:"values"`
	Default  any      `yaml:"default"`
}

type ModuleDefinition struct {
	Enabled bool `yaml:"enabled"`
}

type FileEntry struct {
	Path      string `yaml:"path"`
	Template  string `yaml:"template"`
	Condition string `yaml:"condition"`
}

func fetchPlateYAML(repoURL, branch string) (*PlateYAML, error) {
	apiURL := repoURLToContentsURL(repoURL, branch)

	client := &http.Client{Timeout: 15 * time.Second}
	req, err := http.NewRequest(http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, ErrFetchFailed
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "kikplate-api")

	resp, err := client.Do(req)
	if err != nil {
		return nil, ErrFetchFailed
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrMissingYAML
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return nil, fmt.Errorf("%w: github returned %d (%s)", ErrFetchFailed, resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var ghResp struct {
		Content  string `json:"content"`
		Encoding string `json:"encoding"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&ghResp); err != nil {
		return nil, ErrFetchFailed
	}

	var raw []byte
	if ghResp.Encoding == "base64" {
		raw, err = base64.StdEncoding.DecodeString(ghResp.Content)
		if err != nil {
			return nil, ErrFetchFailed
		}
	} else {
		raw = []byte(ghResp.Content)
	}

	var py PlateYAML
	if err := yaml.Unmarshal(raw, &py); err != nil {
		return nil, ErrFetchFailed
	}
	return &py, nil
}

func fetchRemoteTemplate(templateURL string) (string, error) {
	client := &http.Client{Timeout: 15 * time.Second}
	var lastErr error
	for attempt := 0; attempt < 3; attempt++ {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, templateURL, nil)
		if err != nil {
			cancel()
			return "", fmt.Errorf("%w: cannot fetch template from %s", ErrFetchFailed, templateURL)
		}

		resp, err := client.Do(req)
		if err == nil && resp != nil && resp.StatusCode == http.StatusOK {
			raw, readErr := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
			resp.Body.Close()
			cancel()
			if readErr != nil {
				return "", ErrFetchFailed
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
	return "", fmt.Errorf("%w: cannot fetch template from %s (%v)", ErrFetchFailed, templateURL, lastErr)
}

func repoURLToContentsURL(repoURL, branch string) string {
	repoURL = strings.TrimSuffix(repoURL, ".git")
	repoURL = strings.TrimPrefix(repoURL, "git@github.com:")
	repoURL = strings.TrimPrefix(repoURL, "https://github.com/")

	if branch == "" {
		branch = "main"
	}
	return fmt.Sprintf("https://api.github.com/repos/%s/contents/plate.yaml?ref=%s", repoURL, branch)
}
