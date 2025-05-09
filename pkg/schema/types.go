package schema

type EnvSpec struct {
	Type         string            `yaml:"type"`
	Required     bool              `yaml:"required"`
	OneOf        []string          `yaml:"one_of"`
	RequiredWhen map[string]string `yaml:"required_when"`
	Nullable     bool              `yaml:"nullable"`
}

type Schema map[string]EnvSpec