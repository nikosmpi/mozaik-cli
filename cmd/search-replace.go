package cmd

import (
	"fmt"

	"github.com/nikosmpi/mozaik-cli/wpconfig"
	"github.com/nikosmpi/mozaik-cli/wpdatabase"
	"github.com/spf13/cobra"
)

var searchReplaceCmd = &cobra.Command{
	Use:     "search-replace",
	Aliases: []string{"sr"},
	Short:   "Run search and replace on the local database",
	Long:    `Iterate through all text columns in the local database and perform search and replace based on the configuration list.`,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := wpconfig.GetConfig()
		if err != nil {
			fmt.Println("Error getting config:", err)
			return
		}
		if err := wpdatabase.SearchReplace(config); err != nil {
			fmt.Println("Error during search and replace:", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(searchReplaceCmd)
}
