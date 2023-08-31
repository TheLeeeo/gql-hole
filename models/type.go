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

// // CompileFields builds a partial query for the fields of a type and returns the lines of the query
// func (t *Type) CompileFields() []string {
// 	if t.Kind == ListTypeKind || t.Kind == NonNullTypeKind {
// 		return []string{} //Temporary behavior
// 	}

// 	var fieldNames []string

// 	for _, f := range t.Fields {
// 		baseType := f.Type.GetBaseType()

// 		switch baseType.Kind {
// 		case EnumTypeKind:
// 			fieldNames = append(fieldNames, f.Name)
// 		case ScalarTypeKind:
// 			fieldNames = append(fieldNames, f.Name)
// 		case ObjectTypeKind:
// 			fieldNames = append(fieldNames, fmt.Sprintf("%s {", f.Name))
// 			fieldNames = append(fieldNames, baseType.CompileFields()...)
// 			fieldNames = append(fieldNames, "}")
// 		default:
// 			fmt.Printf("Unhandled type %s in type %s\n", baseType.Kind, f.Name)
// 		}
// 	}

// 	return fieldNames
// }
