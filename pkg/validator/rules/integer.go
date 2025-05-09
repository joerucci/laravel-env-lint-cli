package rules

import (
	"fmt"
	"strconv"

	"github.com/joerucci/laravel-env-lint-cli/pkg/schema"
)

// validateInteger checks if the value is a valid integer,
// or accepts null if nullable is explicitly allowed.
func validateInteger(key, val string, spec schema.EnvSpec) error {
	casted, err := checkNullable(key, val, spec)
	if err != nil {
		return err
	}
	
	// Exit early if null is valid
	if casted == nil {
		return nil
	}

	// We expect casted to still be a string at this point
	strVal, ok := casted.(string)
	if !ok {
		return fmt.Errorf("%s must be a string before parsing (got: %v)", key, casted)
	}

	if _, err := strconv.Atoi(strVal); err != nil {
		return fmt.Errorf("%s must be an integer (got: %s)", key, strVal)
	}

	return nil
}