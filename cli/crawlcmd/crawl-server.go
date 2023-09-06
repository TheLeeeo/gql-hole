package crawlcli

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/TheLeeeo/gql-test-suite/client"
	"github.com/TheLeeeo/gql-test-suite/crawler"
	crawlserver "github.com/TheLeeeo/gql-test-suite/crawler/server.go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	keyHttpPort        = "http-port"
	keyEnablePolling   = "enable-polling"
	keyPollingInterval = "polling-interval"
)

func init() {
	serverCmd.AddCommand(startCmd)

	startCmd.Flags().StringP(keyHttpPort, "l", ":8080", "The port to listen for http traffic on")
	viper.BindPFlag(keyHttpPort, startCmd.Flags().Lookup(keyHttpPort))

	startCmd.Flags().Bool(keyEnablePolling, false, "Enable polling for changes to the target graphql schema")
	viper.BindPFlag(keyEnablePolling, startCmd.Flags().Lookup(keyEnablePolling))

	startCmd.Flags().Int(keyPollingInterval, 10, "The interval in minutes to poll for changes to the target graphql schema")
	viper.BindPFlag(keyPollingInterval, startCmd.Flags().Lookup(keyPollingInterval))
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Control the crawl server",
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the crawl server",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := crawlserver.Config{
			HttpPort: viper.GetString(keyHttpPort),

			CrawlerConfig: crawler.Config{
				ClientConfig: client.Config{
					TargetUrl: viper.GetString(keyTarget),
					Headers:   parseHeaders(viper.GetStringSlice(keyHeaders)),
					PollingConfig: client.PollingConfig{
						Enabled:  viper.GetBool(keyEnablePolling),
						Interval: viper.GetInt(keyPollingInterval),
					},
				},
				Ignore: viper.GetStringSlice(keyIgnore),
			},
		}

		err := validateConfig(&cfg)
		if err != nil {
			fmt.Println("error validating config: ", err)
			os.Exit(1)
		}

		log.Printf("Starting crawl server with config: %+v", cfg)

		s := crawlserver.New(cfg)
		err = s.Run()
		if err != nil {
			log.Println("error running server: ", err)
			os.Exit(1)
		}
	},
}

func validateConfig(cfg *crawlserver.Config) error {
	if cfg.HttpPort == "" {
		return errors.New("no http-port specified")
	}

	if cfg.CrawlerConfig.ClientConfig.PollingConfig.Enabled && cfg.CrawlerConfig.ClientConfig.PollingConfig.Interval < 1 {
		return errors.New("polling interval must be greater than 0")
	}

	return nil
}
