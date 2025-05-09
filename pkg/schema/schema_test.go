package schema

import (
    "testing"

    "gopkg.in/yaml.v3"
)

func TestUnmarshalSchema(t *testing.T) {
    yamlData := `
FOO:
  type: string
  required: true
  one_of: [a, b, c]
  required_when:
    BAR: local
  nullable: false

BAR:
  type: integer
  required: false
  nullable: true
`

    var s Schema
    if err := yaml.Unmarshal([]byte(yamlData), &s); err != nil {
        t.Fatalf("failed to unmarshal YAML: %v", err)
    }

    foo, ok := s["FOO"]
    if !ok {
        t.Fatal("expected key FOO in schema")
    }
    if foo.Type != "string" {
        t.Errorf("FOO.Type = %q; want \"string\"", foo.Type)
    }
    if !foo.Required {
        t.Error("FOO.Required = false; want true")
    }
    if foo.Nullable {
        t.Error("FOO.Nullable = true; want false")
    }
    if len(foo.OneOf) != 3 || foo.OneOf[0] != "a" || foo.OneOf[2] != "c" {
        t.Errorf("FOO.OneOf = %v; want [a b c]", foo.OneOf)
    }
    if v, ok := foo.RequiredWhen["BAR"]; !ok || v != "local" {
        t.Errorf("FOO.RequiredWhen[\"BAR\"] = %q; want \"local\"", v)
    }

    bar, ok := s["BAR"]
    if !ok {
        t.Fatal("expected key BAR in schema")
    }
    if bar.Type != "integer" {
        t.Errorf("BAR.Type = %q; want \"integer\"", bar.Type)
    }
    if bar.Required {
        t.Error("BAR.Required = true; want false")
    }
    if !bar.Nullable {
        t.Error("BAR.Nullable = false; want true")
    }
}