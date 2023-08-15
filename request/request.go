package request

import (
	"encoding/json"
	"fmt"

	"github.com/TheLeeeo/gql-test-suite/models"
)

type request struct {
	Request   string         `json:"query"`
	Variables map[string]any `json:"variables,omitempty"`
}

type RequestType string

const (
	Query    RequestType = "query"
	Mutation RequestType = "mutation"
)

func BuildFromString(requestString string, variables map[string]any) []byte {
	query := &request{
		Request:   requestString,
		Variables: variables,
	}

	b, err := json.Marshal(query)
	if err != nil {
		panic(err)
	}

	return b
}

func Build(requestField *models.Field, variables map[string]any, t RequestType) string {
	if t != Query && t != Mutation {
		panic(fmt.Sprintf("invalid request type: %s", t))
	}

	var input string
	if len(requestField.Args) > 0 {
		input = fmt.Sprintf(" (%s)", requestField.Args[0].Compile())
	}

	requestString := fmt.Sprintf("%s%s{\n%s\n}", t, input, requestField.Compile())

	return string(BuildFromString(requestString, variables))
}
