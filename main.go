package main

import (
	_ "embed"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/TheLeeeo/gql-test-suite/client"
	"github.com/TheLeeeo/gql-test-suite/models"
	"github.com/TheLeeeo/gql-test-suite/request"
	"github.com/TheLeeeo/gql-test-suite/utils"
	"github.com/fatih/color"
)

var addr = "http://localhost:1337/gql"

func main() {
	cl := client.New(addr)
	_ = cl.FetchType("_Entity")
	err := cl.FetchSchema()
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	for _, q := range cl.Queries {
		if q.Name == "_entities" {
			continue
		}
		vars := GenerateMinimalTestData(cl, q)
		r := request.Build(q, vars, request.Query)
		resp, err := cl.Execute([]byte(r))
		if err != nil {
			fmt.Println("Error: ", err)
			os.Exit(1)
		}

		is403 := is403Error(resp)
		var is403string string
		if is403 {
			is403string = color.GreenString("true")
		} else {
			is403string = color.RedString("false")
		}

		fmt.Printf("Query \"%s\" got denied: %s\n", q.Name, is403string)
		if !is403 {
			fmt.Println("	Response: ", string(resp))
		}
	}

	for _, m := range cl.Mutations {
		vars := GenerateMinimalTestData(cl, m)
		r := request.Build(m, vars, request.Mutation)
		resp, err := cl.Execute([]byte(r))
		if err != nil {
			fmt.Println("Error: ", err)
			os.Exit(1)
		}

		is403 := is403Error(resp)
		var is403string string
		if is403 {
			is403string = color.GreenString("true")
		} else {
			is403string = color.RedString("false")
		}

		fmt.Printf("Mutation \"%s\" got denied: %s\n", m.Name, is403string)
		if !is403 {
			fmt.Println("	Response: ", string(resp))
		}
	}
}

func is403Error(resp []byte) bool {
	respType, err := utils.ParseResponse(resp)
	if err != nil {
		panic(err)
	}

	if len(respType.Errors) == 0 {
		return false
	}

	for _, e := range respType.Errors {
		if strings.Contains(e.Message, "Unauthenticated") {
			return true
		}
	}

	return false
}

func GenerateMinimalTestData(cl *client.Client, f *models.Field) map[string]any {
	if len(f.Args) == 0 {
		return nil
	}

	vars := make(map[string]any)

	arg := f.Args[0]

	if arg.Type.Kind != models.NonNullTypeKind {
		fmt.Println("Optional args :)")
		return nil
	}

	baseType := arg.Type.GetBaseType()

	var value any
	switch baseType.Kind {
	case models.EnumTypeKind:
		value = f.Type.EnumValues
	case models.ScalarTypeKind:
		if baseType.Name == "Boolean" {
			value = true
		} else if baseType.Name == "String" {
			value = "0"
		} else if baseType.Name == "Int" {
			value = 0
		} else if baseType.Name == "Float" {
			value = 0.0
		} else if baseType.Name == "ID" {
			value = "0"
		} else if baseType.Name == "Time" {
			value = time.Now()
		} else {
			panic(fmt.Sprintf("Unhandled scalar type %s", baseType.Name))
		}
	case models.InputObjectTypeKind:
		value = GenerateMinimalTestDataForType(cl, baseType)
	default:
		panic(fmt.Sprintf("Unimplemented variable kind %s", f.Type.Kind))
	}

	vars[arg.Name] = value

	return vars
}

func GenerateMinimalTestDataForType(cl *client.Client, t *models.Type) map[string]any {
	vars := make(map[string]any)

	for _, f := range t.InputFields {
		if f.Type.Kind != models.NonNullTypeKind {
			continue
		}

		baseType := f.Type.GetBaseType()

		var value any
		switch baseType.Kind {
		case models.EnumTypeKind:
			value = cl.EnumTypes[baseType.Name][0].Name
		case models.ScalarTypeKind:
			if baseType.Name == "Boolean" {
				value = true
			} else if baseType.Name == "String" {
				value = "0"
			} else if baseType.Name == "Int" {
				value = 0
			} else if baseType.Name == "Float" {
				value = 0.0
			} else if baseType.Name == "ID" {
				value = "0"
			} else if baseType.Name == "Time" {
				value = time.Now()
			} else {
				panic(fmt.Sprintf("Unhandled scalar type %s", baseType.Name))
			}
		case models.InputObjectTypeKind:
			value = GenerateMinimalTestDataForType(cl, baseType)
		default:
			panic(fmt.Sprintf("Unimplemented variable kind %s", f.Type.Kind))
		}

		vars[f.Name] = value
	}

	return vars
}

func LoadExpectedResponse(file string) string {
	expectedBytes, err := os.ReadFile("./expected.json")
	if err != nil {
		panic(fmt.Sprintf("Error reading .json file: File: %s, Error: %s", file, err.Error()))
	}
	// expected := map[string]any{}
	// err = json.Unmarshal(expectedBytes, &expected)
	// if err != nil {
	// 	panic(fmt.Sprintf("Error unmarshaling .json file: Error: %s", err.Error()))
	// }

	return string(expectedBytes)
}
