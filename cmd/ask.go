/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"askme/initialization"
	utils "askme/utils"
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	service "askme/service"

	"github.com/AlecAivazis/survey/v2"
	"github.com/janeczku/go-spinner"
	"github.com/spf13/cobra"

	// "github.com/google/uuid"

	"github.com/tmc/langchaingo/embeddings/openai"
	"github.com/tmc/langchaingo/vectorstores/pinecone"
)

// askCmd represents the ask command
var askCmd = &cobra.Command{
	Use:   "ask",
	Short: "Ask a question to your pdf data",
	Long:  `Ask a question to your pdf data`,
	Run: func(cmd *cobra.Command, args []string) {

		e, err := openai.NewOpenAI()

		ctx := context.Background()

		home, err := os.UserHomeDir()

		if err != nil {
			log.Fatal(err)
		}

		configPath := filepath.Join(home, ".askme", "config.yaml")

		config := initialization.LoadConfig(configPath)

		nameSpaces := utils.GetIndexStats(utils.Config(*config))

		// i := 0
		// formatedNameSpaces := make([]string, len(nameSpaces))
		// for nameSpace := range nameSpaces {
		// 	currIndex := i + 1
		// 	formatedNameSpaces[i] = fmt.Sprint(currIndex) + ": " + nameSpaces[nameSpace]
		// 	i++
		// }

		if len(nameSpaces) == 0 {
			fmt.Println("❌ No PDFs found in the index. Please add a PDF first.")
			return
		}

		namespace := ""
		prompt := &survey.Select{
			Message: "Choose a PDF to ask question:",
			Options: nameSpaces,
		}
		survey.AskOne(prompt, &namespace)

		// Create a new Pinecone vector store.
		store, err := pinecone.New(
			ctx,
			pinecone.WithProjectName(config.PineconeProjectName),
			pinecone.WithIndexName(config.PineconeIndexName),
			pinecone.WithEnvironment(config.PineconeEnvironment),
			pinecone.WithEmbedder(e),
			pinecone.WithAPIKey(config.PineconeAPIKEY),
			pinecone.WithNameSpace(namespace),
		)

		query := ""
		queryPrompt := &survey.Input{
			Message: "Ask your question:",
		}
		survey.AskOne(queryPrompt, &query)

		// Check if the query is empty
		if query == "" {
			fmt.Println("❌ Please enter a valid query")
			return
		}

		s := spinner.StartNew("Fetching Results ... ")

		docs, err := store.SimilaritySearch(ctx, query, 1)
		if err != nil {
			log.Fatal(err)
		}
		s.Stop()

		service.OpenAi_Call(query, docs[0].PageContent)
		// fmt.Println(docs[0].PageContent)

	},
}

func init() {
	rootCmd.AddCommand(askCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// askCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// askCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
