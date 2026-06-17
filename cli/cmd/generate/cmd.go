package generate

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate [slug]",
		Short: "Generate a project from a plate",
		Args:  cobra.MaximumNArgs(1),
		RunE:  run,
	}

	cmd.Flags().StringP("file", "f", "", "YAML file with template values")
	cmd.Flags().StringArray("set", nil, "Set a template variable (key=value), repeatable")
	cmd.Flags().String("output-dir", "", "Output directory (default: ./<slug>)")
	cmd.Flags().String("repo", "", "Push generated project to a remote git repository")
	cmd.Flags().String("template", "", "Path to a local plate directory (bypasses server)")
	cmd.Flags().Bool("force", false, "Force push when using --repo")

	return cmd
}

func run(cmd *cobra.Command, args []string) error {
	templateDir, _ := cmd.Flags().GetString("template")
	outputDir, _ := cmd.Flags().GetString("output-dir")
	valuesFile, _ := cmd.Flags().GetString("file")
	valuesFlag, _ := cmd.Flags().GetStringArray("set")
	repoURL, _ := cmd.Flags().GetString("repo")
	force, _ := cmd.Flags().GetBool("force")

	if templateDir == "" && len(args) == 0 {
		return fmt.Errorf("provide a plate slug or --template <dir>")
	}

	values, err := loadValues(valuesFile, valuesFlag)
	if err != nil {
		return err
	}

	interactive := valuesFile == "" && len(valuesFlag) == 0

	var zipBytes []byte
	var slug string

	if templateDir != "" {
		slug = filepath.Base(templateDir)
		py, err := loadPlateYAML(templateDir)
		if err != nil {
			return err
		}
		if interactive {
			if err := promptSchema(py, values); err != nil {
				return err
			}
		}
		zipBytes, err = generateLocal(templateDir, py, values)
	} else {
		slug = args[0]
		if interactive {
			schema, err := fetchServerSchema(cmd, slug)
			if err != nil {
				return fmt.Errorf("cannot fetch plate schema: %w", err)
			}
			if err := promptSchema(schema, values); err != nil {
				return err
			}
		}
		zipBytes, err = generateFromServer(cmd, slug, values)
	}
	if err != nil {
		return err
	}

	if outputDir == "" {
		outputDir = filepath.Base(slug)
	}

	if repoURL != "" {
		return generateToRepo(zipBytes, repoURL, force)
	}
	return extractZip(zipBytes, outputDir)
}

func loadValues(valuesFile string, setPairs []string) (map[string]any, error) {
	values := map[string]any{}
	if valuesFile != "" {
		raw, err := os.ReadFile(valuesFile)
		if err != nil {
			return nil, fmt.Errorf("cannot read values file: %w", err)
		}
		if err := yaml.Unmarshal(raw, &values); err != nil {
			return nil, fmt.Errorf("cannot parse values file: %w", err)
		}
	}
	for _, p := range setPairs {
		k, v, found := strings.Cut(p, "=")
		if !found {
			return nil, fmt.Errorf("invalid --set value %q, expected key=value", p)
		}
		key := strings.TrimSpace(k)
		if key == "" {
			return nil, fmt.Errorf("invalid --set value %q, key cannot be empty", p)
		}
		values[key] = parseScalar(strings.TrimSpace(v))
	}
	return values, nil
}

func parseScalar(v string) any {
	lower := strings.ToLower(v)
	switch lower {
	case "true":
		return true
	case "false":
		return false
	case "null", "nil":
		return nil
	}

	if i, err := strconv.ParseInt(v, 10, 64); err == nil {
		return i
	}
	if f, err := strconv.ParseFloat(v, 64); err == nil {
		return f
	}

	return v
}
