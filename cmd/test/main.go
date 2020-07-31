package main

import (
	"context"
	"log"
	"os"

	"github.com/PullRequestInc/go-gpt3"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		log.Fatalln("Missing API KEY")
	}

	ctx := context.Background()
	client := gpt3.NewClient(apiKey)

	resp, err := client.Completion(ctx, gpt3.CompletionRequest{
		Prompt: []string{
			"1\n2\n3\n4",
		},
		MaxTokens: gpt3.IntPtr(10),
	})
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("%+v\n", resp)

	resp, err = client.Completion(ctx, gpt3.CompletionRequest{
		Prompt: []string{
			"go:golang\npy:python\njs:",
		},
		Stop: []string{"\n"},
	})
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("%+v\n", resp)
}
