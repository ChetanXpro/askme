/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"askme/initialization"

	utils "askme/utils"

	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/janeczku/go-spinner"
	"github.com/spf13/cobra"
	"github.com/tmc/langchaingo/embeddings/openai"
	"github.com/tmc/langchaingo/vectorstores/pinecone"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a pdf file to the vector store",
	Long:  `Add a pdf file to the vector store`,

	Run: func(cmd *cobra.Command, args []string) {

		e, err := openai.NewOpenAI()
		if err != nil {
			log.Fatal(err)
		}

		ctx := context.Background()

		home, err := os.UserHomeDir()

		if err != nil {
			log.Fatal(err)
		}

		configPath := filepath.Join(home, ".askme", "config.yaml")

		config := initialization.LoadConfig(configPath)

		path := ""
		prompt := &survey.Input{
			Message: "Enter the path of the pdf file:",
		}
		survey.AskOne(prompt, &path)

		// To extract name from pdf path
		isNameClear := strings.Contains(path, "/")
		if isNameClear {
			s := strings.Split(path, "/")
			path = s[len(s)-1]
		}

		// Create a new Pinecone vector store.
		store, err := pinecone.New(
			ctx,
			pinecone.WithProjectName(config.PineconeProjectName),
			pinecone.WithIndexName(config.PineconeIndexName),
			pinecone.WithEnvironment(config.PineconeEnvironment),
			pinecone.WithEmbedder(e),
			pinecone.WithAPIKey(config.PineconeAPIKEY),
			pinecone.WithNameSpace(path),
		)
		if err != nil {
			log.Fatal(err)
		}

		pdfdocs, err := utils.LoadPDF(path)

		if err != nil {
			log.Fatal(err)
		}

		s := spinner.StartNew("Adding document to the vector store...")

		// Add documents to the Pinecone vector store.
		err = store.AddDocuments(context.Background(),
			pdfdocs)
		if err != nil {
			log.Fatal(err)
		}
		s.Stop()
		fmt.Println("✓ Document saved")

	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
