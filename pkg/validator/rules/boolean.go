package rules

import (
	"fmt"

	"github.com/joerucci/laravel-env-lint-cli/pkg/schema"
)

// validateBoolean checks if the value is a valid boolean,
// or optionally allows null if the schema marks it as nullable.
func validateBoolean(key, val string, spec schema.EnvSpec) error {
	casted, err := checkNullable(key, val, spec)
	if err != nil {
		return err
	}
	
	// Exit early if null is valid
	if casted == nil {
		return nil
	}

	// Check for actual boolean
	b, ok := casted.(bool)
	if !ok {
		return fmt.Errorf("%s must be a boolean (got: %v)", key, casted)
	}
	_ = b // optional: can remove if unused

	return nil
}