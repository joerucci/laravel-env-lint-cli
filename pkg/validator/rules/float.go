package rules

import (
	"fmt"
	"strconv"

	"github.com/joerucci/laravel-env-lint-cli/pkg/schema"
)

// validateFloat checks if the value is a valid floating point number,
// or accepts null if the schema marks the variable as nullable.
//
// The value is first cast using Laravel-style .env parsing rules
// (e.g., "null", "(null)", etc.), and must be a string that can be
// parsed as a float64 unless explicitly null.
//
// Example:
//   TAX_RATE=7.25          → valid
//   TAX_RATE=null        → valid if nullable: true
//   TAX_RATE=abc           → invalid
func validateFloat(key, val string, spec schema.EnvSpec) error {
	casted, err := checkNullable(key, val, spec)
	if err != nil {
		return err
	}
	
	// Exit early if null is valid
	if casted == nil {
		return nil
	}

	strVal, ok := casted.(string)
	if !ok {
		return fmt.Errorf("%s must be a string before parsing (got: %v)", key, casted)
	}

	if _, err := strconv.ParseFloat(strVal, 64); err != nil {
		return fmt.Errorf("%s must be a float (got: %s)", key, strVal)
	}

	return nil
}