package models

import "fmt"

// Introspection types implemented as specified by https://spec.graphql.org/October2021/#sec-Introspection

type Schema struct {
	Description      string       `json:"description"`
	Types            []*Type      `json:"types"`
	QueryType        *Type        `json:"queryType"`        // Will only carry the name
	MutationType     *Type        `json:"mutationType"`     // Will only carry the name
	SubscriptionType *Type        `json:"subscriptionType"` // Will only carry the name
	Directives       []*Directive `json:"directives"`
}

type InputValue struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	Type         *Type  `json:"type"`
	DefaultValue string `json:"defaultValue"`
}

func (i *InputValue) Compile() string {
	argTypeString := buildArgTypeString(i.Type)

	return fmt.Sprintf("$%s: %s", i.Name, argTypeString)
}

// Builds a string representation of the type, eg. "String", "String!", "[String]!", "[String!]!"
func buildArgTypeString(t *Type) string {
	if t.Kind == NonNullTypeKind {
		return fmt.Sprintf("%s!", buildArgTypeString(t.OfType))
	}

	if t.Kind == ListTypeKind {
		return fmt.Sprintf("[%s]", buildArgTypeString(t.OfType))
	}

	return string(t.Name)
}
