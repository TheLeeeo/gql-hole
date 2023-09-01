package cmd

import (
	"encoding/json"
	"log"
	"os"
	"strings"

	"github.com/TheLeeeo/gql-test-suite/crawler"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	keyAddr    = "addr"
	keyIgnore  = "ignore"
	keyHeaders = "headers"
	verbose    = "verbose"
)

func init() {
	// add a flag for the graphql endpoint
	crawlCmd.Flags().StringP(keyAddr, "a", "", "The graphql endpoint to crawl")
	viper.BindPFlag(keyAddr, crawlCmd.Flags().Lookup(keyAddr))

	crawlCmd.Flags().StringSliceP(keyIgnore, "i", []string{}, "Queries and mutations to ignore")
	viper.BindPFlag(keyIgnore, crawlCmd.Flags().Lookup(keyIgnore))

	crawlCmd.Flags().StringSliceP(keyHeaders, "H", []string{}, "Headers to send with the request, formatted like \"k1:v1,k2,v2\"")
	viper.BindPFlag(keyHeaders, crawlCmd.Flags().Lookup(keyHeaders))

	crawlCmd.Flags().BoolP(verbose, "v", false, "Verbose output")
	viper.BindPFlag(verbose, crawlCmd.Flags().Lookup(verbose))
}

var crawlCmd = &cobra.Command{
	Use:   "crawl",
	Short: "Crawl the endpoints of a graphql server",
	Long: `Crawl the endpoints of a graphql server. This will fetch the schema and
crawl all queries and mutations`,
	Run: func(cmd *cobra.Command, args []string) {
		addr := viper.GetString(keyAddr)
		if addr == "" {
			log.Println("Error: no graphql endpoint specified")
			os.Exit(1)
		}

		ignore := viper.GetStringSlice(keyIgnore)
		headerSlice := viper.GetStringSlice(keyHeaders)

		headers := parseHeaders(headerSlice)

		cfg := &crawler.Config{
			Addr:    addr,
			Headers: headers,
			Ignore:  ignore,
		}

		cl := crawler.New(cfg)
		err := cl.LoadSchema()
		if err != nil {
			log.Println("error loading schema: ", err)
			os.Exit(1)
		}

		ops := cl.Crawl()

		for _, op := range ops {
			op.PrintResult()
		}

		if viper.GetBool(verbose) {
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
