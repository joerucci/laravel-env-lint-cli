package rules

import (
	"testing"
	"github.com/joerucci/laravel-env-lint-cli/pkg/schema"
)

func TestIsInList(t *testing.T) {
	tests := []struct {
		value    string
		list     []string
		expected bool
	}{
		{"foo", []string{"foo", "bar"}, true},
		{"bar", []string{"foo", "bar"}, true},
		{"baz", []string{"foo", "bar"}, false},
		{"", []string{"foo", "bar"}, false},
		{"FOO", []string{"foo", "bar"}, false}, // case-sensitive
	}

	for _, tt := range tests {
		t.Run(tt.value, func(t *testing.T) {
			result := isInList(tt.value, tt.list)
			if result != tt.expected {
				t.Errorf("isInList(%q) = %v, expected %v", tt.value, result, tt.expected)
			}
		})
	}
}

func TestCheckNullable(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		spec     schema.EnvSpec
		expected any
		wantErr  bool
	}{
		{"null allowed", "(null)", schema.EnvSpec{Nullable: true}, nil, false},
		{"null disallowed", "(null)", schema.EnvSpec{Nullable: false}, nil, true},
		{"normal string", "value", schema.EnvSpec{Nullable: true}, "value", false},
		{"empty string", "", schema.EnvSpec{Nullable: false}, "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := checkNullable("TEST_KEY", tt.value, tt.spec)

			if (err != nil) != tt.wantErr {
				t.Errorf("expected error: %v, got: %v", tt.wantErr, err)
			}

			if !tt.wantErr && result != tt.expected {
				t.Errorf("expected result: %v, got: %v", tt.expected, result)
			}
		})
	}
}

func TestCastLaravelEnvValue(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{"true", true},
		{"(true)", true},
		{"false", false},
		{"(false)", false},
		{"null", nil},
		{"(null)", nil},
		{"empty", ""},
		{"(empty)", ""},
		{"My App", "My App"},
		{"123", "123"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := castLaravelEnvValue(tt.input)

			// nil comparison is special in Go
			if result == nil && tt.expected != nil {
				t.Errorf("expected %v, got nil", tt.expected)
			} else if result != nil && tt.expected == nil {
				t.Errorf("expected nil, got %v", result)
			} else if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}