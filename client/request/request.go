package request

import (
	"encoding/json"
)

type Request struct {
	Body      string
	Variables map[string]any
	Headers   map[string]string
}

type RequestType string

const (
	Query    RequestType = "query"
	Mutation RequestType = "mutation"
)

func NewEmpty() *Request {
	return &Request{}
}

func New(body string, variables map[string]any) *Request {
	return &Request{
		Body:      body,
		Variables: variables,
	}
}

func BuildFromString(requestString string, variables map[string]any) []byte {
	query := &Request{
		Body:      requestString,
		Variables: variables,
	}

	return query.Build()
}

// Build compiles the request into a byte array that can be sent to the server.
func (r *Request) Build() []byte {
	type requestInternal struct {
		Query     string         `json:"query"`
		Variables map[string]any `json:"variables"`
	}

	req := &requestInternal{
		Query:     r.Body,
		Variables: r.Variables,
	}

	b, err := json.Marshal(req)
	if err != nil {
		panic(err)
	}

	return b
}
