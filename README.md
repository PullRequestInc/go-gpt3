# go-gpt3

A OpenAPI GPT-3 API client enabling Go/Golang programs to interact with the gpt3 APIs.

[![PkgGoDev](https://pkg.go.dev/badge/github.com/PullRequestInc/go-gpt3)](https://pkg.go.dev/github.com/PullRequestInc/go-gpt3)

## Usage

Simple usage to call the main gpt-3 API, completion:

```go
client := gpt3.NewClient(apiKey)
resp, err := client.Completion(ctx, gpt3.CompletionRequest{
    Prompt: []string{"2, 3, 5, 7, 11,"},
})

fmt.Print(resp.Choices[0].Text)
// prints " 13, 17, 19, 23, 29, 31", etc
```

### Full Examples

Try out any of these examples with putting the contents in a `main.go` and running `GO111MODULE=on go run main.go`.
You will also need to have a `.env` file that looks like:

```
API_KEY=<openAI API Key>
```

```go
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
		Prompt:    []string{"The first thing you should know about javascript is"},
		MaxTokens: gpt3.IntPtr(30),
		Stop:      []string{"."},
		Echo:      true,
	})
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(resp.Choices[0].Text)
}
```

## Support

- [x] List Engines API
- [x] Get Engine API
- [x] Completion API (this is the main gpt-3 API)
- [x] Document Search API
- [x] Overriding default url, user-agent, timeout, and other options
- [ ] TODO: Streaming on the Completion API

## Powered by

[<img src="https://www.pullrequest.com/images/pullrequest-logo.svg" width="200">](https://www.pullrequest.com)
