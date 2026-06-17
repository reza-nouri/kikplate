package generate

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	gen "github.com/kickplate/plategenerator"
)

func extractZip(zipBytes []byte, destDir string) error {
	if _, err := os.Stat(destDir); err == nil {
		return fmt.Errorf("directory %q already exists", destDir)
	}
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("cannot create output directory: %w", err)
	}
	if err := unzipInto(zipBytes, destDir); err != nil {
		return err
	}
	fmt.Printf("Project generated in: %s\n", destDir)
	return nil
}

func unzipInto(zipBytes []byte, destDir string) error {
	r, err := zip.NewReader(bytes.NewReader(zipBytes), int64(len(zipBytes)))
	if err != nil {
		return fmt.Errorf("cannot read zip: %w", err)
	}
	base := filepath.Clean(destDir) + string(os.PathSeparator)
	for _, f := range r.File {
		dest := filepath.Join(destDir, filepath.Clean(f.Name))
		if !strings.HasPrefix(dest, base) {
			return fmt.Errorf("invalid zip path: %s", f.Name)
		}
		if f.FileInfo().IsDir() {
			os.MkdirAll(dest, 0755)
			continue
		}
		if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
			return err
		}
		out, err := os.Create(dest)
		if err != nil {
			return err
		}
		rc, err := f.Open()
		if err != nil {
			out.Close()
			return err
		}
		_, copyErr := io.Copy(out, rc)
		rc.Close()
		out.Close()
		if copyErr != nil {
			return copyErr
		}
	}
	return nil
}

func buildZip(py *plateYAML, templateDir string, data map[string]any) ([]byte, error) {
	return gen.BuildZip(toShared(py), data, func(ref string) (string, error) {
		return resolveTemplateContent(templateDir, ref)
	})
}
