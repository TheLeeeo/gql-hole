package crawler

import "github.com/TheLeeeo/gql-test-suite/client"

type Config struct {
	ClientConfig client.Config

	// Operations to ignore
	Ignore []string
}
