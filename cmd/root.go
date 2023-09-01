package cmd

import "github.com/spf13/cobra"

func init() {
	RootCmd.AddCommand(crawlCmd)
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
