package cmd

import (
	"fmt"
	"time"

	"github.com/nikosmpi/mozaik-cli/wpconfig"
	"github.com/nikosmpi/mozaik-cli/wpdatabase"
	"github.com/spf13/cobra"
)

var allCmd = &cobra.Command{
	Use:   "all",
	Short: "Run both sync and search-replace",
	Long:  `Run both staging-to-local synchronization and database search-and-replace sequentially.`,
	Run: func(cmd *cobra.Command, args []string) {
		start := time.Now()
		config, err := wpconfig.GetConfig()
		if err != nil {
			fmt.Println("Error getting config:", err)
			return
		}

		fmt.Println("Step 1: Synchronizing from staging...")
		if err := wpdatabase.SyncStagingToLocal(config); err != nil {
			fmt.Println("Error during sync:", err)
			return
		}

		fmt.Println("\nStep 2: Running search and replace...")
		if err := wpdatabase.SearchReplace(config); err != nil {
			fmt.Println("Error during search and replace:", err)
			return
		}

		duration := time.Since(start)
		fmt.Printf("\nAll operations completed successfully in %v!\n", duration.Round(time.Second))
	},
}

func init() {
	rootCmd.AddCommand(allCmd)
}
