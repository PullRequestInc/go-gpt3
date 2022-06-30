package main

import (
	"context"
	"fmt"
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

	fmt.Print("\n\nstarting stream:\n")

	request := gpt3.CompletionRequest{
		Prompt: []string{
			"One thing that you should know about golang",
		},
		MaxTokens: gpt3.IntPtr(20),
	}

	err = client.CompletionStream(ctx, request, func(resp *gpt3.CompletionResponse) {
		fmt.Println(resp.Choices[0].Text)
	})
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Print("\n\nedits API:\n")

	editsResponse, err := client.Edits(ctx, gpt3.EditsRequest{
		Model:       "text-davinci-edit-001",
		Input:       "What day of the wek is it?",
		Instruction: "Fix the spelling mistakes",
	})
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("%+v\n", editsResponse)
}
