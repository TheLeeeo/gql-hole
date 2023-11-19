package cli

import (
	"encoding/json"
	"log"
	"os"

	"github.com/TheLeeeo/gql-test-suite/client"
	"github.com/spf13/cobra"
)

var targetURL string
var filePath string

func init() {
	executeFileCmd.PersistentFlags().StringVarP(&targetURL, "target-url", "t", "", "The graphql endpoint to crawl")
	executeFileCmd.PersistentFlags().StringVarP(&filePath, "file", "f", "", "The file containing the query to execute")
}

var executeFileCmd = &cobra.Command{
	Use:   "execute",
	Short: "execute a graphql query from a file",
	Run: func(cmd *cobra.Command, args []string) {
		if targetURL == "" {
			log.Println("error: no graphql endpoint specified")
			os.Exit(1)
		}
		if filePath == "" {
			log.Println("error: no file specified")
			os.Exit(1)
		}

		client := client.New(targetURL)
		res, err := client.ExecuteFile(filePath)
		if err != nil {
			log.Println("error executing file: ", err)
			os.Exit(1)
		}

		b, err := json.MarshalIndent(res, "", "  ")
		if err != nil {
			log.Println("error marshalling result: ", err)
			os.Exit(1)
		}

		log.Println(string(b))
	},
}
