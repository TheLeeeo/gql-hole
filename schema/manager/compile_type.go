package manager

import (
	"fmt"

	"github.com/TheLeeeo/gql-test-suite/schema"
)

// CompileFields builds a partial query for the fields of a type and returns the lines of the query
func (m *Manager) CompileType(t schema.Type) []string {
	if t.Kind == schema.ListTypeKind || t.Kind == schema.NonNullTypeKind {
		return []string{} //Temporary behavior
	}

	var fieldNames []string

	for _, f := range t.Fields {
		// Field is optional
		if f.Type.Kind != schema.NonNullTypeKind {
			continue
		}

		baseType := f.Type.GetBaseType()

		switch baseType.Kind {
		case schema.EnumTypeKind:
			fieldNames = append(fieldNames, f.Name)
		case schema.ScalarTypeKind:
			fieldNames = append(fieldNames, f.Name)
		case schema.ObjectTypeKind:
			fullBaseType := m.Types[baseType.Name]

			fieldNames = append(fieldNames, fmt.Sprintf("%s {", f.Name))
			fieldNames = append(fieldNames, m.CompileType(fullBaseType)...)
			fieldNames = append(fieldNames, "}")
		default:
			fmt.Printf("Unhandled type %s in type %s\n", baseType.Kind, f.Name)
		}
	}

	// At least one selection field is required
	if len(fieldNames) == 0 {
		fieldNames = append(fieldNames, "__typename")
	}

	return fieldNames
}

func (m *Manager) CompileField(f schema.Field) string {
	baseType := f.Type.GetBaseType()

	var queryBody string
	if baseType.Kind == schema.ScalarTypeKind {
		queryBody = ""
	} else if baseType.Kind == schema.ObjectTypeKind {
		completeBaseType := m.Types[baseType.Name]
		compiledTypeList := m.CompileType(completeBaseType)

		var fields string
		for _, l := range compiledTypeList {
			fields += fmt.Sprintf("\n%s", l)
		}

		queryBody = fmt.Sprintf("{%s\n}", fields)
	} else {
		panic(fmt.Sprintf("Unhandled type %s in field %s", baseType.Kind, f.Name))
	}

	var input string
	if len(f.Args) > 0 {
		input = fmt.Sprintf(" (%s: $%s)", f.Args[0].Name, f.Args[0].Name)
	}

	return fmt.Sprintf("%s%s%s", f.Name, input, queryBody)
}
