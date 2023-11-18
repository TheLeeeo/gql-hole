package schema

// Introspection types implemented as specified by https://spec.graphql.org/October2021/#sec-Introspection

type Schema struct {
	Description      string      `json:"description"`
	Types            []Type      `json:"types"`
	QueryType        *Type       `json:"queryType"`        // Will only carry the name
	MutationType     *Type       `json:"mutationType"`     // Will only carry the name
	SubscriptionType *Type       `json:"subscriptionType"` // Will only carry the name
	Directives       []Directive `json:"directives"`
}

func (s *Schema) GetType(name string) *Type {
	for _, t := range s.Types {
		if t.Name == name {
			return &t
		}
	}

	return nil
}
