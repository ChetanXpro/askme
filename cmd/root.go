package cmd

import (
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "askme",
		Short: "Chat with PDF Cli app",
		Long:  `Askme is a cli app which answer users query , first a user can pdf data into vectorDB , and after that user can ask any question related to that pdf data`,
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}
