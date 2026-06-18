package plate

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

func (s *plateService) fetchKickplateYAML(repoURL, branch string) (*KickplateYAML, error) {
	return s.fetchKickplateYAMLWithOptions(repoURL, branch, false)
}

func (s *plateService) fetchKickplateYAMLWithOptions(repoURL, branch string, forceRefresh bool) (*KickplateYAML, error) {
	manifest, err := s.fetchManifestYAML(repoURL, branch, "plate.yaml", forceRefresh)
	if err == nil && strings.TrimSpace(manifest.Owner) != "" {
		return manifest, nil
	}

	if err != nil && !errors.Is(err, ErrMissingYAML) {
		return nil, err
	}

	legacy, legacyErr := s.fetchManifestYAML(repoURL, branch, "kikplate.yaml", forceRefresh)
	if legacyErr != nil {
		if err == nil {
			return nil, ErrMissingYAML
		}
		if errors.Is(err, ErrMissingYAML) && errors.Is(legacyErr, ErrMissingYAML) {
			return nil, ErrMissingYAML
		}
		return nil, legacyErr
	}

	if strings.TrimSpace(legacy.Owner) == "" {
		return nil, ErrMissingYAML
	}

	return legacy, nil
}

func (s *plateService) fetchManifestYAML(repoURL, branch, filename string, forceRefresh bool) (*KickplateYAML, error) {
	apiURL := repoURLToContentsURL(repoURL, branch, filename)
	if forceRefresh {
		apiURL = fmt.Sprintf("%s&_nonce=%d", apiURL, time.Now().UnixNano())
	}
	req, err := http.NewRequest(http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, ErrFetchFailed
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "kickplate-api")
	if forceRefresh {
		req.Header.Set("Cache-Control", "no-cache")
		req.Header.Set("Pragma", "no-cache")
	}

	resp, err := http.DefaultClient.Do(req)
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

	var manifest KickplateYAML
	if err := yaml.Unmarshal(raw, &manifest); err != nil {
		return nil, ErrFetchFailed
	}

	return &manifest, nil
}

func repoURLToContentsURL(repoURL, branch, filename string) string {
	return fmt.Sprintf("https://api.github.com/repos/%s/contents/%s?ref=%s",
		extractRepoPath(repoURL), filename, url.QueryEscape(strings.TrimSpace(branch)))
}

func extractRepoPath(repoURL string) string {
	repoURL = strings.TrimSpace(repoURL)

	repoURL = strings.TrimPrefix(repoURL, "git@github.com:")

	for _, prefix := range []string{
		"https://github.com/",
		"http://github.com/",
		"github.com/",
	} {
		if strings.HasPrefix(repoURL, prefix) {
			repoURL = strings.TrimPrefix(repoURL, prefix)
			break
		}
	}

	repoURL = strings.TrimSuffix(repoURL, ".git")
	repoURL = strings.Trim(repoURL, "/")
	return repoURL
}
