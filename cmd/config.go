package cmd

import (
	"fmt"

	"github.com/nikosmpi/mozaik-cli/wpconfig"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Print the configuration",
	Long:  `Print the configuration.`,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := wpconfig.GetConfig()
		if err != nil {
			fmt.Println("Error getting config:", err)
			return
		}
		fmt.Println("Config:", config)
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
