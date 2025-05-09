package main

import (
    "encoding/base64"
    "fmt"
    "strings"
    
    "github.com/spf13/cobra"
    
    "github.com/joerucci/laravel-env-lint-cli/pkg/loader"
    "github.com/joerucci/laravel-env-lint-cli/pkg/validator"
)

var (
    envFile    string
    schemaFile string
    encrypted bool
    encKey    string
)

func init() {
    validateCmd.Flags().StringVar(&envFile, "env-file", ".env", "Path to .env file to validate")
    validateCmd.Flags().StringVar(&schemaFile, "schema", "env.schema.yaml", "Path to YAML schema file")
    validateCmd.Flags().BoolVar(&encrypted, "encrypted", false, "Decrypt and validate .env.encrypted instead of .env")
    validateCmd.Flags().StringVar(&encKey, "key", "", "32-character encryption key for .env.encrypted")
    rootCmd.AddCommand(validateCmd)
}

var validateCmd = &cobra.Command{
    Use:   "validate",
    Short: "Validate a .env file (optionally encrypted) against a schema",
    RunE: func(cmd *cobra.Command, args []string) error {
        var (
            envMap map[string]string
            err    error
        )

        if encrypted {
            if encKey == "" {
                return fmt.Errorf("--key is required when using --encrypted")
            }
            keyBytes, err := normalizeKey(encKey)
            if err != nil {
                return err
            }
            envMap, err = loader.LoadEncryptedEnv(envFile+ ".encrypted", keyBytes)
        } else {
            envMap, err = loader.LoadEnv(envFile)
        }

        if err != nil {
            return fmt.Errorf("failed to load env: %w", err)
        }

        schema, err := loader.LoadSchema(schemaFile)
        if err != nil {
            return fmt.Errorf("failed to load schema: %w", err)
        }
       
        valid := validator.Validate(envMap, schema)
        if !valid {
            return fmt.Errorf("validation failed")
        }
        fmt.Println("✅ env-check complete.")
        return nil
    },
    SilenceUsage:  true,
    SilenceErrors:  true,
}

// normalizeKey parses raw or base64 Laravel-style keys.
func normalizeKey(input string) ([]byte, error) {
    if strings.HasPrefix(input, "base64:") {
        decoded, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(input, "base64:"))
        if err != nil {
            return nil, fmt.Errorf("invalid base64 key: %w", err)
        }
        if len(decoded) != 32 {
            return nil, fmt.Errorf("decoded base64 key must be 32 bytes")
        }
        return decoded, nil
    }

    if len(input) != 32 {
        return nil, fmt.Errorf("raw key must be exactly 32 characters")
    }
    return []byte(input), nil
}