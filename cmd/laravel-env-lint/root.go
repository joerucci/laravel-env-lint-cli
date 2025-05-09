package main

import (
	"fmt"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   			"laravel-env-lint",
	Version:       	"v1.0.0",
	Short: 			"Validate .env files against a schema",
	Long:  			`A fast CLI to lint and validate your Laravel-style .env files using a YAML schema.`,
	SilenceUsage: true,
	SilenceErrors:  true,
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("❌", err)
	}
}