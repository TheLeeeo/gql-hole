package client

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
	QueryRequest    RequestType = "query"
	MutationRequest RequestType = "mutation"
)

func NewRequest(body string, variables map[string]any) *Request {
	return &Request{
		Body:      body,
		Variables: variables,
	}
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
