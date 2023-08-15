package models

// Response is implemented as specified by https://spec.graphql.org/October2021/#sec-Response
type Response struct {
	Errors     []Error        `json:"errors"`
	Data       map[string]any `json:"data"`
	Extensions map[string]any `json:"extensions"`
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
