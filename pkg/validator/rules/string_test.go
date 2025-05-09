package rules

import (
	"testing"

	"github.com/joerucci/laravel-env-lint-cli/pkg/schema"
)

func TestValidateString(t *testing.T) {
	tests := []struct {
		name    string
		val     string
		spec    schema.EnvSpec
		wantErr bool
	}{
		{"basic string", "hello", schema.EnvSpec{}, false},
		{"empty string", "", schema.EnvSpec{}, false},
		{"nullable null", "(null)", schema.EnvSpec{Nullable: true}, false},
		{"non-nullable null", "(null)", schema.EnvSpec{Nullable: false}, true},
		{"one_of valid", "foo", schema.EnvSpec{OneOf: []string{"foo", "bar"}}, false},
		{"one_of invalid", "baz", schema.EnvSpec{OneOf: []string{"foo", "bar"}}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateString("TEST_KEY", tt.val, tt.spec)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateString(%q, spec=%+v) error = %v, wantErr = %v",
					tt.val, tt.spec, err, tt.wantErr)
			}
		})
	}
}