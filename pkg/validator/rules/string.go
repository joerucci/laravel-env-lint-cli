package rules

import (
	"fmt"

	"github.com/joerucci/laravel-env-lint-cli/pkg/schema"
)

// validateString checks if the value is a string and matches any one_of constraints,
// or accepts null if nullable is explicitly allowed.
func validateString(key, val string, spec schema.EnvSpec) error {
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
		return fmt.Errorf("%s must be a string (got: %v)", key, casted)
	}

	if len(spec.OneOf) > 0 && !isInList(strVal, spec.OneOf) {
		return fmt.Errorf("%s must be one of %v (got: %s)", key, spec.OneOf, strVal)
	}

	return nil
}