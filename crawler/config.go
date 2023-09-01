package crawler

type Config struct {
	// The graphql endpoint to crawl
	Addr string

	// Headers to send with the request
	Headers map[string]string

	// Queries and mutations to ignore
	Ignore []string
}
