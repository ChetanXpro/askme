package server

import (
	"context"
	"fmt"
	"log"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/schema"
)

func OpenAi_Call(prompt string, docContext string) {
	llm, err := openai.NewChat()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	llm_prompt := fmt.Sprintf(`You are a helpful AI assistant. Use the following pieces of context to answer the question at the end.
	If you don't know the answer, just say you don't know. DO NOT try to make up an answer.
	If the question is not related to the context, politely respond that you are tuned to only answer questions that are related to the context.

	Context is Written between three dashes.

	---
	%s
	---
	
	Question: %s
	Helpful answer:`, docContext, prompt)

	completion, err := llm.Call(ctx, []schema.ChatMessage{

		schema.HumanChatMessage{Content: llm_prompt},
	}, llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
		fmt.Print(string(chunk))

		return nil
	}),
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Sprintf(completion.Content)
}
