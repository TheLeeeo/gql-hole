package schema

import "fmt"

type Type struct {
	Kind           TypeKind     `json:"kind"`
	Name           string       `json:"name"`
	Description    string       `json:"description"`
	Fields         []Field      `json:"fields"`
	Interfaces     []Type       `json:"interfaces"`
	PossibleTypes  []Type       `json:"possibleTypes"`
	EnumValues     []EnumValue  `json:"enumValues"`
	InputFields    []InputValue `json:"inputFields"`
	OfType         *Type        `json:"ofType"`
	SpecifiedByURL string       `json:"specifiedByURL"`
}

type TypeKind string

const (
	ScalarTypeKind      TypeKind = "SCALAR"
	ObjectTypeKind      TypeKind = "OBJECT"
	InterfaceTypeKind   TypeKind = "INTERFACE"
	UnionTypeKind       TypeKind = "UNION"
	EnumTypeKind        TypeKind = "ENUM"
	InputObjectTypeKind TypeKind = "INPUT_OBJECT"
	ListTypeKind        TypeKind = "LIST"
	NonNullTypeKind     TypeKind = "NON_NULL"
)

func (t *Type) GetBaseType() *Type {
	if t.OfType == nil {
		return t
	}

	return t.OfType.GetBaseType()
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

func (t *Type) GetDefaultValue() any {
	if t.Kind == NonNullTypeKind {
		return t.OfType.GetDefaultValue()
	}

	if t.Kind == ListTypeKind {
		return []any{}
	}

	if t.Kind == ScalarTypeKind {
		if t.Name == "Boolean" {
			return true
		} else if t.Name == "String" {
			return "0"
		} else if t.Name == "Int" {
			return 0
		} else if t.Name == "Float" {
			return 0.0
		} else if t.Name == "ID" {
			return "0"
		} else if t.Name == "Time" {
			return "2019-01-01T00:00:00Z"
		} else {
			panic(fmt.Sprintf("Unhandled scalar type %s", t.Name))
		}
	}

	return nil
}
