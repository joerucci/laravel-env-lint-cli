package rules

import (
	"testing"

	"github.com/joerucci/laravel-env-lint-cli/pkg/schema"
)

func TestValidateBoolean(t *testing.T) {
	tests := []struct {
		name     string
		val      string
		spec     schema.EnvSpec
		wantErr  bool
	}{
		{"valid true", "true", schema.EnvSpec{}, false},
		{"valid (true)", "(true)", schema.EnvSpec{}, false},
		{"valid false", "false", schema.EnvSpec{}, false},
		{"valid (false)", "(false)", schema.EnvSpec{}, false},
		{"nullable (null)", "(null)", schema.EnvSpec{Nullable: true}, false},
		{"non-nullable (null)", "(null)", schema.EnvSpec{Nullable: false}, true},
		{"invalid boolean", "yes", schema.EnvSpec{}, true},
		{"invalid numeric", "0", schema.EnvSpec{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateBoolean("BOOLEAN_VAR", tt.val, tt.spec)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateBoolean(%q) error = %v, wantErr = %v", tt.val, err, tt.wantErr)
			}
		})
	}
}