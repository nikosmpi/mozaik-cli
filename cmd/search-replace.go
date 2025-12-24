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
	Short:   "Search and replace in the project",
	Long:    `Search and replace in the project by creating necessary configuration files and directories.`,
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
