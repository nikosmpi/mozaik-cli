package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var passwordCmd = &cobra.Command{
	Use:     "password",
	Aliases: []string{"pass"},
	Short:   "Generate a password",
	Long:    `Generate a password`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Generating password...")
	},
}

func init() {
	rootCmd.AddCommand(passwordCmd)
}
