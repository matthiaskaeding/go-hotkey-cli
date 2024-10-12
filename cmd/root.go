package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "sdb",
	Short: "sdb is a cli tool for quickly getting strings",
	Long:  "sdb is a cli tool for quickly getting strings, which are copied into the clipboard",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Oops. An error while executing sdb '%s'\n", err)
		os.Exit(1)
	}
}
