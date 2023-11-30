package cmd

import (
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Setup the configuration for the app",
	Long:  `Setup the configuration for the app.`,
	Run: func(cmd *cobra.Command, args []string) {

		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		configPath := home + "/.askme"
		viper.SetConfigName("config")
		viper.AddConfigPath(configPath)
		viper.SetConfigType("yaml")

		// Check if the .askme directory exists and create it if it doesn't
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			err = os.MkdirAll(configPath, os.ModePerm)
			if err != nil {
				fmt.Printf("Unable to create config directory: %v\n", err)
				return
			}
		} else {
			fmt.Println("Config directory already exists.")

			overwrite := ""
			prompt := &survey.Input{
				Message: "Do you want to overwrite the existing config file? (y/n):",
				Default: "n",
			}
			survey.AskOne(prompt, &overwrite)

			if overwrite == "n" {
				return
			} else {
				// If user want to overwrite the existing config file, delete it

				err = os.Remove(configPath + "/config.yaml")
				if err != nil {
					fmt.Printf("Unable to delete existing config file: %v\n", err)
					return
				}
			}

		}
		all_env := [5]string{"openai_key", "pinecone_project_name", "pinecone_index_name", "pinecone_environment", "pinecone_api_key"}

		for i := 0; i < len(all_env); i++ {

			curr_env := ""
			prompt := &survey.Input{
				Message: "Enter " + all_env[i] + ":",
			}
			survey.AskOne(prompt, &curr_env)
			viper.Set(all_env[i], curr_env)

		}

		if err := viper.SafeWriteConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {

				err = viper.WriteConfigAs(configPath + "/config.yaml")

				if err != nil {
					fmt.Printf("Error occurred while writing to config file: %v\n", err)
				}
			} else {

				fmt.Printf("An error occurred: %v\n", err)
			}
		}

		fmt.Println("API key has been stored in configuration. Path: " + configPath + "/config.yaml")

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
