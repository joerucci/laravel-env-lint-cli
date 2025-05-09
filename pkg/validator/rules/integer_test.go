package rules

import (
	"testing"

	"github.com/joerucci/laravel-env-lint-cli/pkg/schema"
)

func TestValidateInteger(t *testing.T) {
	tests := []struct {
		name     string
		val      string
		spec     schema.EnvSpec
		wantErr  bool
	}{
		{"valid int", "123", schema.EnvSpec{}, false},
		{"valid negative int", "-456", schema.EnvSpec{}, false},
		{"invalid alpha", "abc", schema.EnvSpec{}, true},
		{"invalid decimal", "3.14", schema.EnvSpec{}, true},
		{"valid nullable null", "(null)", schema.EnvSpec{Nullable: true}, false},
		{"invalid nullable false", "(null)", schema.EnvSpec{Nullable: false}, true},
		{"empty string", "", schema.EnvSpec{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateInteger("PORT", tt.val, tt.spec)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateInteger(%q) error = %v, wantErr = %v", tt.val, err, tt.wantErr)
			}
		})
	}
}