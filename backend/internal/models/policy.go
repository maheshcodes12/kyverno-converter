package models

// This file defines Go structs that represent the structure of Kyverno policies.
// Using structs allows for type-safe parsing and manipulation of the YAML data.

// LegacyClusterPolicy represents the structure of a kyverno.io/v1 ClusterPolicy.
type LegacyClusterPolicy struct {
	APIVersion string     `yaml:"apiVersion"`
	Kind       string     `yaml:"kind"`
	Metadata   Metadata   `yaml:"metadata"`
	Spec       LegacySpec `yaml:"spec"`
}

type Metadata struct {
	Name        string                 `yaml:"name"`
	Annotations map[string]interface{} `yaml:"annotations,omitempty"`
}

type LegacySpec struct {
	ValidationFailureAction string       `yaml:"validationFailureAction"`
	Background              *bool        `yaml:"background,omitempty"`
	Rules                   []LegacyRule `yaml:"rules"`
}

type LegacyRule struct {
	Name     string        `yaml:"name"`
	Match    MatchBlock    `yaml:"match"`
	Exclude  MatchBlock    `yaml:"exclude,omitempty"`
	Validate ValidateBlock `yaml:"validate"`
}

type MatchBlock struct {
	Any []ResourceFilter `yaml:"any,omitempty"`
	All []ResourceFilter `yaml:"all,omitempty"`
}

type ResourceFilter struct {
	Resources struct {
		Kinds []string `yaml:"kinds"`
	} `yaml:"resources"`
}

// ValidateBlock holds the core validation logic. It can contain a pattern, a foreach, or other validation types.
// We use `interface{}` for Pattern because its structure is highly dynamic.
type ValidateBlock struct {
	Message string      `yaml:"message"`
	Pattern interface{} `yaml:"pattern,omitempty"`
	ForEach []ForEach   `yaml:"foreach,omitempty"`
	// Add other validation types like `deny` here if needed in the future
}

type ForEach struct {
	List    string      `yaml:"list"`
	Pattern interface{} `yaml:"pattern"`
}

// --- Target Policy Structures ---

// ValidatingPolicy represents the modern policies.kyverno.io/v1alpha1 policy.
type ValidatingPolicy struct {
	APIVersion string         `yaml:"apiVersion"`
	Kind       string         `yaml:"kind"`
	Metadata   Metadata       `yaml:"metadata"`
	Spec       ValidatingSpec `yaml:"spec"`
}

type ValidatingSpec struct {
	ValidationActions  []string         `yaml:"validationActions"`
	Background         *bool            `yaml:"background,omitempty"`
	MatchConstraints   MatchConstraints `yaml:"matchConstraints"`
	ExcludeConstraints MatchConstraints `yaml:"exclude,omitempty"`
	Validations        []Validation     `yaml:"validations"`
}

type MatchConstraints struct {
	ResourceRules []ResourceRule `yaml:"resourceRules,omitempty"`
}

type ResourceRule struct {
	APIGroups   []string `yaml:"apiGroups"`
	APIVersions []string `yaml:"apiVersions"`
	Operations  []string `yaml:"operations"`
	Resources   []string `yaml:"resources"`
}

type Validation struct {
	Message    string `yaml:"message"`
	Expression string `yaml:"expression"`
}
