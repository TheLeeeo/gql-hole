package crawler

import "github.com/TheLeeeo/gql-test-suite/introspection"

type Config struct {
	ClientConfig introspection.Config

	// Operations to ignore
	Ignore []string
}
