package client

type Config struct {
	// The graphql endpoint to crawl
	Addr string

	// Headers to send with the request
	Headers map[string]string
}
