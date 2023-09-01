package cli

import (
	crawlcli "github.com/TheLeeeo/gql-test-suite/cli/crawlcmd"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	KeyVerbose = "verbose"
)

func init() {
	RootCmd.AddCommand(crawlcli.CrawlCmd)

	RootCmd.PersistentFlags().BoolP(KeyVerbose, "v", false, "Verbose output")
	viper.BindPFlag(KeyVerbose, RootCmd.PersistentFlags().Lookup(KeyVerbose))
}

var RootCmd = &cobra.Command{
	Use:   "gts",
	Short: "gts is a graphql test suite",
	Long: `gts is a graphql test suite. It is designed to test graphql servers
by generating queries and mutations based on the schema.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}
