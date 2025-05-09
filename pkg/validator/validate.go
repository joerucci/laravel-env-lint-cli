package validator

import (
	"fmt"
	
	"github.com/joerucci/laravel-env-lint-cli/pkg/validator/rules"
	"github.com/joerucci/laravel-env-lint-cli/pkg/schema"
)

func Validate(envMap map[string]string, schema schema.Schema) bool {
	hasError := false

	for key, spec := range schema {
		if err := validateKey(key, spec, envMap); err != nil {
			fmt.Println("❌", err.Error())
			hasError = true
		}
	}

	return !hasError
}

// validateKey checks one key against the spec and returns an error 
// if the value is missing or invalid.
func validateKey(key string, spec schema.EnvSpec, env map[string]string) error {
	val, exists := env[key]

	if isConditionallyRequired(spec, env) && !exists {
		return fmt.Errorf("%s is required by condition", key)
	}

	if spec.Required && !exists {
		return fmt.Errorf("%s is required but missing", key)
	}

	if exists {
		return validateTypeAndConstraints(key, val, spec)
	}

	return nil
}

// isConditionallyRequired returns true if the given environment variable
// should be considered required, based on `required_when` conditions
// defined in the schema.EnvSpec and the current environment state.
//
// For example, if a variable is required when APP_ENV = "local", this
// function returns true when that condition is met.
func isConditionallyRequired(spec schema.EnvSpec, env map[string]string) bool {
	for depKey, depVal := range spec.RequiredWhen {
		if env[depKey] == depVal {
			return true
		}
	}
	return false
}

// validateTypeAndConstraints checks the value using a type-specific validator.
func validateTypeAndConstraints(key, val string, spec schema.EnvSpec) error {
	validatorFunc, ok := rules.TypeValidators[spec.Type]
	if !ok {
		return fmt.Errorf("%s has unknown type: %s", key, spec.Type)
	}
	return validatorFunc(key, val, spec)
}