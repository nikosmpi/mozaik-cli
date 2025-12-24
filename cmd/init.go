package cmd

import (
	"fmt"

	"github.com/nikosmpi/mozaik-cli/wpconfig"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the project",
	Long:  `Initialize the project by creating necessary configuration files and directories.`,
	Run: func(cmd *cobra.Command, args []string) {
		tf, err := wpconfig.SaveConfig()
		if err != nil {
			fmt.Println("Error saving config:", err)
			return
		}
		if tf {
			fmt.Println("moz-config.json created")
			if err := wpconfig.AddGitignore(); err != nil {
				fmt.Println("Error adding gitignore:", err)
				return
			}
			fmt.Println("moz-config.json added to gitignore")
			return
		}
		fmt.Println("moz-config.json already exists")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
