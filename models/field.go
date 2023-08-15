package models

import "fmt"

type Field struct {
	Name              string       `json:"name"`
	Description       string       `json:"description"`
	Args              []InputValue `json:"args"`
	Type              *Type        `json:"type"`
	IsDeprecated      bool         `json:"isDeprecated"`
	DeprecationReason string       `json:"deprecationReason"`
}

func (f *Field) Compile() string {
	baseType := f.Type.GetBaseType()

	var queryBody string
	if baseType.Kind == ScalarTypeKind {
		queryBody = ""
	} else if baseType.Kind == ObjectTypeKind {
		compiledTypeList := baseType.CompileFields()

		var fields string
		for _, l := range compiledTypeList {
			fields += fmt.Sprintf("\n%s", l)
		}

		queryBody = fmt.Sprintf("{%s\n}", fields)
	} else if baseType.Kind == UnionTypeKind {

	} else {
		panic(fmt.Sprintf("Unhandled type %s in field %s", baseType.Kind, f.Name))
	}

	var input string
	if len(f.Args) > 0 {
		input = fmt.Sprintf(" (%s: $%s)", f.Args[0].Name, f.Args[0].Name)
	}

	return fmt.Sprintf("%s%s%s", f.Name, input, queryBody)
}
