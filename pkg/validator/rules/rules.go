package rules

import (
	"github.com/joerucci/laravel-env-lint-cli/pkg/schema"
)

type TypeValidator func(key, val string, spec schema.EnvSpec) error

var TypeValidators = map[string]TypeValidator{
	"integer": validateInteger,
	"float": validateFloat,
	"boolean": validateBoolean,
	"string":  validateString,
}