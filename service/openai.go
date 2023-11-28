package server

import (
	"context"
	openai "github.com/sashabaranov/go-openai"
	"log"
	"time"
)

type OpenAIResponse struct {
	Response string `json:"response"`
	Tokens   int    `json:"tokens"`
}
