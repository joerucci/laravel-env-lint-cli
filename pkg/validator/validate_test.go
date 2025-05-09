package validator

import (
	"testing"

	"github.com/joerucci/laravel-env-lint-cli/pkg/schema"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name     string
		envMap   map[string]string
		schema   schema.Schema
		wantValid bool
	}{
		{
			name: "all good",
			envMap: map[string]string{
				"FOO": "hello",
				"BAR": "123",
				"BAZ": "true",
			},
			schema: schema.Schema{
				"FOO": {Type: "string", Required: true},
				"BAR": {Type: "integer", Required: true},
				"BAZ": {Type: "boolean", Required: true},
			},
			wantValid: true,
		},
		{
			name: "missing required",
			envMap: map[string]string{
				"FOO": "hello",
			},
			schema: schema.Schema{
				"FOO": {Type: "string", Required: true},
				"BAR": {Type: "integer", Required: true},
			},
			wantValid: false,
		},
		{
			name: "conditional required_when",
			envMap: map[string]string{
				"MODE": "local",
				"DEV":  "enabled",
			},
			schema: schema.Schema{
				"MODE": {Type: "string", Required: true},
				"DEV":  {Type: "string", RequiredWhen: map[string]string{"MODE": "local"}},
			},
			wantValid: true,
		},
		{
			name: "conditional missing",
			envMap: map[string]string{
				"MODE": "local",
			},
			schema: schema.Schema{
				"MODE": {Type: "string", Required: true},
				"DEV":  {Type: "string", RequiredWhen: map[string]string{"MODE": "local"}},
			},
			wantValid: false,
		},
		{
			name: "unknown type",
			envMap: map[string]string{"X": "foo"},
			schema: schema.Schema{"X": {Type: "foobar", Required: true}},
			wantValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Validate(tt.envMap, tt.schema)
			if got != tt.wantValid {
				t.Errorf("Validate() = %v; want %v", got, tt.wantValid)
			}
		})
	}
}