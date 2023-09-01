package crawlserver

import (
	"github.com/TheLeeeo/gql-test-suite/crawler"
)

type Config struct {
	// The address to listen on
	HttpPort string

	CrawlerConfig *crawler.Config
}
