package response

import (
	"encoding/json"
)

// Response is implemented as specified by https://spec.graphql.org/October2021/#sec-Response
type Response struct {
	Errors     []Error        `json:"errors"`
	Data       map[string]any `json:"data"`
	Extensions map[string]any `json:"extensions"`

	// This is not part of the spec, but is used to store the status code of the response
	StatusCode int `json:"-"`
}

type Error struct {
	Message    string         `json:"message"`
	Locations  []Location     `json:"locations"`
	Path       []string       `json:"path"`
	Extensions map[string]any `json:"extensions"`
}

type Location struct {
	Line   int `json:"line"`
	Column int `json:"column"`
}

func Parse(resp []byte) (*Response, error) {
	response := &Response{}
	err := json.Unmarshal(resp, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}
