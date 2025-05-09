package rules

import (
	"fmt"

	"github.com/joerucci/laravel-env-lint-cli/pkg/schema"
)

func isInList(val string, list []string) bool {
	for _, item := range list {
		if val == item {
			return true
		}
	}
	return false
}

// checkNullable casts the .env value using Laravel-style casting
// and verifies whether it's allowed to be null based on the schema.
// Returns the casted value (which may be nil) or an error if null is not allowed.
func checkNullable(key string, val string, spec schema.EnvSpec) (any, error) {
	casted := castLaravelEnvValue(val)

	if casted == nil {
		if spec.Nullable {
			return nil, nil
		}
		return nil, fmt.Errorf("%s cannot be null", key)
	}

	return casted, nil
}

// castLaravelEnvValue converts Laravel-style .env values
// like "(true)", "false", "null", etc., into native Go types.
func castLaravelEnvValue(val string) any {
	switch val {
	case "true", "(true)":
		return true
	case "false", "(false)":
		return false
	case "null", "(null)":
		return nil
	case "empty", "(empty)":
		return ""
	default:
		return val
	}
}