package cmd

import (
	"fmt"

	"github.com/nikosmpi/mozaik-cli/wpconfig"
	"github.com/nikosmpi/mozaik-cli/wpdatabase"
	"github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync database from staging to local",
	Long:  `Synchronize the staging database to your local environment using SSH dump and local MySQL import.`,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := wpconfig.GetConfig()
		if err != nil {
			fmt.Println("Error getting config:", err)
			return
		}
		if err := wpdatabase.SyncStagingToLocal(config); err != nil {
			fmt.Println("Error during sync:", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
}
