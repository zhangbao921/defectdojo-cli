package cmd

import (
	"fmt"
	"os"

	"github.com/defectdojo-cli/pkg"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	host   string
	apiKey string
	client *pkg.Client
)

var rootCmd = &cobra.Command{
	Use:   "dd",
	Short: "DefectDojo CLI - API v2 client",
	Long: `A command-line client for DefectDojo API v2.
Manage products, engagements, tests, findings, and more.

Examples:
  dd products list
  dd findings get 42
  dd findings list --severity Critical --limit 50
  dd import-scan --product-id 1 --engagement-id 1 --scan-type "ZAP Scan" --file results.json
  dd products create --data '{"name":"MyApp","description":"test"}'
`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Skip for root help
		if cmd.Name() == "dd" && !cmd.HasParent() && len(args) == 0 {
			return nil
		}

		// Resolve host: flag > env > config > default
		if host == "" {
			host = os.Getenv("DD_HOST")
		}
		if host == "" {
			host = viper.GetString("host")
		}
		if host == "" {
			host = "https://demo.defectdojo.org"
		}

		// Resolve API key: flag > env > config
		if apiKey == "" {
			apiKey = os.Getenv("DD_API_KEY")
		}
		if apiKey == "" {
			apiKey = viper.GetString("api-key")
		}
		if apiKey == "" {
			return fmt.Errorf("API key required. Set DD_API_KEY env var, --api-key flag, or configure in ~/.dd/config.yaml")
		}

		client = pkg.NewClient(host, apiKey)
		return nil
	},
	SilenceUsage: true,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&host, "host", "", "DefectDojo host URL (env: DD_HOST)")
	rootCmd.PersistentFlags().StringVar(&apiKey, "api-key", "", "API key (env: DD_API_KEY)")
	rootCmd.PersistentFlags().StringP("output", "o", "table", "Output format: table, json, yaml")

	registerResourceCommands()
	registerScanCommands()
	registerSpecialCommands()
}

func initConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.dd/")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	viper.ReadInConfig()
}
