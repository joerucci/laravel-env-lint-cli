package rules

import (
	"testing"

	"github.com/joerucci/laravel-env-lint-cli/pkg/schema"
)

func TestValidateFloat(t *testing.T) {
	tests := []struct {
		name     string
		val      string
		spec     schema.EnvSpec
		wantErr  bool
	}{
		{"valid float", "3.14", schema.EnvSpec{}, false},
		{"valid float with int format", "42", schema.EnvSpec{}, false},
		{"valid nullable null", "(null)", schema.EnvSpec{Nullable: true}, false},
		{"invalid nullable null", "(null)", schema.EnvSpec{Nullable: false}, true},
		{"invalid alpha", "abc", schema.EnvSpec{}, true},
		{"invalid mixed", "42abc", schema.EnvSpec{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateFloat("FLOAT_VAR", tt.val, tt.spec)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateFloat(%q) error = %v, wantErr = %v", tt.val, err, tt.wantErr)
			}
		})
	}
}