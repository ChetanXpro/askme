package main

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/tmc/langchaingo/embeddings/openai"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/vectorstores"
	"github.com/tmc/langchaingo/vectorstores/pinecone"
)

func createDocuments(field1Values []string, field2Values []string) []schema.Document {
	var documents []schema.Document

	for _, field1Value := range field1Values {
		documents = append(documents, schema.Document{
			PageContent: field1Value,
		})
	}

	return documents
}

func main() {
	// Create an embeddings client using the OpenAI API. Requires environment variable OPENAI_API_KEY to be set.
	e, err := openai.NewOpenAI()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	// Create a new Pinecone vector store.
	store, err := pinecone.New(
		ctx,
		pinecone.WithProjectName("d94161e"),
		pinecone.WithIndexName("go"),
		pinecone.WithEnvironment("us-central1-gcp"),
		pinecone.WithEmbedder(e),
		pinecone.WithAPIKey(""),
		pinecone.WithNameSpace(uuid.New().String()),
	)

	if err != nil {
		log.Fatal(err)
	}

	// Add documents to the Pinecone vector store.
	err = store.AddDocuments(context.Background(), []schema.Document{
		{
			PageContent: "Tokyo",
		},
		{
			PageContent: "Paris",
		},
		{
			PageContent: "London",
		},
		{
			PageContent: "Santiago",
		},
		{
			PageContent: "Buenos Aires",
		},
		{
			PageContent: "Rio de Janeiro",
		},
		{
			PageContent: "Sao Paulo",
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	docs, err := store.SimilaritySearch(ctx, "japan", 1)
	fmt.Println(docs)

	// Search for similar documents using score threshold.
	docs, err = store.SimilaritySearch(ctx, "only cities in south america", 10, vectorstores.WithScoreThreshold(0.80))
	fmt.Println(docs)

	// Search for similar documents using score threshold and metadata filter.
	filter := map[string]interface{}{
		"$and": []map[string]interface{}{
			{
				"area": map[string]interface{}{
					"$gte": 1000,
				},
			},
			{
				"population": map[string]interface{}{
					"$gte": 15.5,
				},
			},
		},
	}

	docs, err = store.SimilaritySearch(ctx, "only cities in south america",
		10,
		vectorstores.WithScoreThreshold(0.80),
		vectorstores.WithFilters(filter))
	fmt.Println(docs)
}
