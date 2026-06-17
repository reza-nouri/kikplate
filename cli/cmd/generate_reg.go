package cmd

import "github.com/kickplate/cli/cmd/generate"

func init() {
	rootCmd.AddCommand(generate.NewCommand())
}
