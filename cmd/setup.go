package cmd

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Setup llm configuration",
	Long:  `Setup LLM configurations like openai key`,
	Run: func(cmd *cobra.Command, args []string) {
		openai_key := ""
		prompt := &survey.Input{
			Message: "Enter your openai API key",
			Help:    "",
		}
		survey.AskOne(prompt, &openai_key)
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setupCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// setupCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
