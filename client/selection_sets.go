package client

import (
	"fmt"

	"github.com/TheLeeeo/gql-test-suite/models"
)

// CompileFields builds a partial query for the fields of a type and returns the lines of the query
func (c *Client) CompileType(t *models.Type) []string {
	if t.Kind == models.ListTypeKind || t.Kind == models.NonNullTypeKind {
		return []string{} //Temporary behavior
	}

	var fieldNames []string

	for _, f := range t.Fields {
		// Field is optional
		if f.Type.Kind != models.NonNullTypeKind {
			continue
		}

		baseType := f.Type.GetBaseType()

		switch baseType.Kind {
		case models.EnumTypeKind:
			fieldNames = append(fieldNames, f.Name)
		case models.ScalarTypeKind:
			fieldNames = append(fieldNames, f.Name)
		case models.ObjectTypeKind:
			fullBaseType := c.GetType(baseType.Name)

			fieldNames = append(fieldNames, fmt.Sprintf("%s {", f.Name))
			fieldNames = append(fieldNames, c.CompileType(fullBaseType)...)
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

func (c *Client) CompileField(f *models.Field) string {
	baseType := f.Type.GetBaseType()

	var queryBody string
	if baseType.Kind == models.ScalarTypeKind {
		queryBody = ""
	} else if baseType.Kind == models.ObjectTypeKind {
		completeBaseType := c.GetType(baseType.Name)
		compiledTypeList := c.CompileType(completeBaseType)

		var fields string
		for _, l := range compiledTypeList {
			fields += fmt.Sprintf("\n%s", l)
		}

		queryBody = fmt.Sprintf("{%s\n}", fields)
	} else if baseType.Kind == models.UnionTypeKind {
		// TODO?
	} else {
		panic(fmt.Sprintf("Unhandled type %s in field %s", baseType.Kind, f.Name))
	}

	var input string
	if len(f.Args) > 0 {
		input = fmt.Sprintf(" (%s: $%s)", f.Args[0].Name, f.Args[0].Name)
	}

	return fmt.Sprintf("%s%s%s", f.Name, input, queryBody)
}
