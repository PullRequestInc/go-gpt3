# go-gpt3

A OpenAPI GPT-3 API client enabling Go/Golang programs to interact with the gpt3 APIs.

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

Try out any of these examples with putting the contents in a `main.go` and running `go run main.go`

```go
package main

import (
    "fmt"
    "github.com/PullRequestInc/go-gpt3"
)

ctx := context.Background()
client := gpt3.NewClient(apiKey)

resp, err := client.Completion(ctx, gpt3.CompletionRequest{
    Prompt: []string{`
go:golang
py:python
js:`},
	MaxTokens: gpt3.IntPtr(20),
})
```

```go
package main

import (
    "fmt"
    "github.com/PullRequestInc/go-gpt3"
)

ctx := context.Background()
client := gpt3.NewClient(apiKey)

resp, err := client.Completion(ctx, gpt3.CompletionRequest{
    Prompt: []string{`
There are several things you should know about PullRequest.

1. PullRequest provides external code review as a service for all languages.
2. If you are new to Go, check out our blog posts like https://www.pullrequest.com/blog/unit-testing-in-go.
3.`},
    Stop: []string{"\n"},
})
```

## Support

- [x] List Engines API
- [x] Get Engine API
- [x] Completion API (this is the main gpt-3 API)
- [x] Document Search API
- [x] Overriding default url, user-agent, timeout, and other options
- [ ] TODO: Streaming on the Completion API
