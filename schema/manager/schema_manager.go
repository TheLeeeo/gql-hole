package manager

import (
	"fmt"
	"log"

	"github.com/TheLeeeo/gql-test-suite/client"
	"github.com/TheLeeeo/gql-test-suite/schema"
)

type Manager struct {
	Types     map[string]schema.Type
	Queries   map[string]schema.Field
	Mutations map[string]schema.Field
}

func New(s *schema.Schema) *Manager {
	m := &Manager{
		Types:     make(map[string]schema.Type),
		Queries:   make(map[string]schema.Field),
		Mutations: make(map[string]schema.Field),
	}

	for _, t := range s.Types {
		m.Types[t.Name] = t
	}

	queries, ok := m.Types["Query"]
	if ok && len(queries.Fields) > 0 {
		for _, f := range queries.Fields {
			m.Queries[f.Name] = f
		}
	}

	mutations, ok := m.Types["Mutation"]
	if ok && len(mutations.Fields) > 0 {
		for _, f := range mutations.Fields {
			m.Mutations[f.Name] = f
		}
	}

	_, ok = m.Types["Subscription"]
	if ok {
		log.Println("Schema contains subscriptions, these are not supported :(")
	}

	return m
}

func (c *Manager) Build(requestField schema.Field, t client.RequestType) string {
	if t != client.QueryRequest && t != client.MutationRequest {
		panic(fmt.Sprintf("invalid request type: %s", t))
	}

	var input string
	if len(requestField.Args) > 0 {
		input = fmt.Sprintf(" (%s)", requestField.Args[0].Compile())
	}

	requestString := fmt.Sprintf("%s%s{\n%s\n}", t, input, c.CompileField(requestField))

	return requestString
}
