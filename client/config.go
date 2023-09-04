package client

type Config struct {
	// The graphql endpoint to crawl
	TargetUrl string

	// Headers to send with the request
	Headers map[string]string

	PollingConfig PollingConfig
}

type PollingConfig struct {
	// Poll for chenges in the schema
	Enabled bool
	// The number of minutes to wait between polls
	Interval int
}
