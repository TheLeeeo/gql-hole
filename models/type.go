package models

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

type EnumValue struct {
	Name              string `json:"name"`
	Description       string `json:"description"`
	IsDeprecated      bool   `json:"isDeprecated"`
	DeprecationReason string `json:"deprecationReason"`
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
