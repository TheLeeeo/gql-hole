package main

import (
	_ "embed"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/TheLeeeo/gql-test-suite/client"
	"github.com/TheLeeeo/gql-test-suite/models"
	"github.com/TheLeeeo/gql-test-suite/request"
	"github.com/TheLeeeo/gql-test-suite/utils"
	"github.com/fatih/color"
)

var addr = "http://localhost:4010/gql"

// var addr = "http://localhost:4000"

func main() {
	log.SetOutput(os.Stderr)
	log.SetFlags(0)

	cl := client.New(addr)
	err := cl.FetchSchema()
	if err != nil {
		log.Println("Error: ", err)
		os.Exit(1)
	}

	var allowedRequests []string
	for _, q := range cl.Queries {
		if q.Name == "_entities" {
			continue
		}
		vars := GenerateMinimalTestData(cl, q)
		r := cl.Build(q, vars, request.Query)
		resp, err := cl.Execute([]byte(r))
		if err != nil {
			log.Panic("Error: ", err)
		}

		if printTestResult(resp, request.Query, q.Name) {
			allowedRequests = append(allowedRequests, q.Name)
		}
	}

	for _, m := range cl.Mutations {
		vars := GenerateMinimalTestData(cl, m)
		r := cl.Build(m, vars, request.Mutation)
		resp, err := cl.Execute([]byte(r))
		if err != nil {
			log.Panic("Error: ", err)
		}

		if printTestResult(resp, request.Mutation, m.Name) {
			allowedRequests = append(allowedRequests, m.Name)
		}
	}

	log.SetOutput(os.Stdout)
	log.Print(allowedRequests)
}

func printTestResult(resp []byte, requestType request.RequestType, name string) bool {
	var resultString string

	failedToFetch := isFetchFailed(resp)
	isDenied := is401Error(resp) || is403Error(resp)

	if failedToFetch {
		resultString = color.YellowString("FAILED TO FETCH")
	} else if isDenied {
		resultString = color.GreenString("DENIED")
	} else {
		resultString = color.RedString("ALLOWED")
	}

	log.Printf("%s \"%s\": %s\n", requestType, name, resultString)
	// Got allowed
	if !failedToFetch && !isDenied {
		log.Println("	Response: ", string(resp))
		return true
	}

	return false
}

func testIntro(cl *client.Client) error {
	q := utils.LoadQuery("introspection.gql")
	req := request.BuildFromString(q, nil)
	resp, err := cl.Execute([]byte(req))
	if err != nil {
		return fmt.Errorf("error executing request: %v", err)
	}

	respType, err := utils.ParseResponse(resp)
	if err != nil {
		return err
	}
	dataMap := respType.Data["__schema"].(map[string]any)

	sch := &models.Schema{}
	err = utils.ParseMap(dataMap, sch)
	if err != nil {
		return fmt.Errorf("error parsing schema: %v", err)
	}

	for _, t := range sch.Types {
		if t.Name == "Query" || t.Name == "Mutation" {
			continue
		}

		fmt.Println("Fetched type: ", t.Name)
	}

	return nil
}

func isFetchFailed(resp []byte) bool {
	respType, err := utils.ParseResponse(resp)
	if err != nil {
		panic(err)
	}

	if len(respType.Errors) == 0 {
		return false
	}

	for _, e := range respType.Errors {
		if strings.Contains(e.Message, "HTTP fetch failed") {
			return true
		}
	}

	return false
}

func is401Error(resp []byte) bool {
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

func is403Error(resp []byte) bool {
	respType, err := utils.ParseResponse(resp)
	if err != nil {
		panic(err)
	}

	if len(respType.Errors) == 0 {
		return false
	}

	for _, e := range respType.Errors {
		if strings.Contains(e.Message, "PermissionDenied") {
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
		//Optional args :)
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
		completeBaseType := cl.GetInputType(baseType.Name)
		value = GenerateMinimalTestDataForType(cl, completeBaseType)
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
			completeBaseType := cl.GetInputType(baseType.Name)
			value = GenerateMinimalTestDataForType(cl, completeBaseType)
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
