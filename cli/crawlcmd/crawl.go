package crawlcli

import (
	"encoding/json"
	"log"
	"os"
	"strings"

	"github.com/TheLeeeo/gql-test-suite/client"
	"github.com/TheLeeeo/gql-test-suite/crawler"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	keyTarget  = "target-url"
	keyIgnore  = "ignore"
	keyHeaders = "headers"
	keyVerbose = "verbose"
)

func init() {
	CrawlCmd.AddCommand(crawlRunCmd)
	CrawlCmd.AddCommand(serverCmd)

	// add a flag for the graphql endpoint
	CrawlCmd.PersistentFlags().StringP(keyTarget, "t", "", "The graphql endpoint to crawl")
	viper.BindPFlag(keyTarget, CrawlCmd.PersistentFlags().Lookup(keyTarget))

	CrawlCmd.PersistentFlags().StringSliceP(keyIgnore, "i", []string{}, "Queries and mutations to ignore")
	viper.BindPFlag(keyIgnore, CrawlCmd.PersistentFlags().Lookup(keyIgnore))

	CrawlCmd.PersistentFlags().StringSliceP(keyHeaders, "H", []string{}, "Headers to send with the request, formatted like \"k1:v1,k2,v2\"")
	viper.BindPFlag(keyHeaders, CrawlCmd.PersistentFlags().Lookup(keyHeaders))
}

var CrawlCmd = &cobra.Command{
	Use:   "crawl",
	Short: "For crawling a graphql endpoint",
}

var crawlRunCmd = &cobra.Command{
	Use:   "run",
	Short: "Perform a crawl",
	Run: func(cmd *cobra.Command, args []string) {
		addr := viper.GetString(keyTarget)
		if addr == "" {
			log.Println("error: no graphql endpoint specified")
			os.Exit(1)
		}

		ignore := viper.GetStringSlice(keyIgnore)
		headerSlice := viper.GetStringSlice(keyHeaders)

		headers := parseHeaders(headerSlice)

		cfg := &crawler.Config{
			ClientConfig: &client.Config{
				TargetAddr: addr,
				Headers:    headers,
			},
			Ignored: ignore,
		}

		c := crawler.New(cfg)

		ops, err := c.Crawl()
		if err != nil {
			log.Println("error crawling: ", err)
			os.Exit(1)
		}

		for _, op := range ops {
			op.PrintResult()
		}

		if viper.GetBool(keyVerbose) {
			b, err := json.MarshalIndent(ops, "", "  ")
			if err != nil {
				log.Println("error marshalling operations: ", err)
				os.Exit(1)
			}

			log.Println(string(b))
		}
	},
}

func parseHeaders(headers []string) map[string]string {
	headerMap := make(map[string]string)

	for _, header := range headers {
		split := strings.Split(header, ":")
		if len(split) != 2 {
			log.Println("invalid header: ", header)
			os.Exit(1)
		}

		headerMap[split[0]] = split[1]
	}

	return headerMap
}
