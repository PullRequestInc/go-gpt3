package gpt3

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Engine Types
const (
	AdaEngine     = "ada"
	BabbageEngine = "babbage"
	CurieEngine   = "curie"
	DavinciEngine = "davinci"
	DefaultEngine = DavinciEngine
)

const (
	defaultBaseURL        = "https://api.openai.com/v1"
	defaultUserAgent      = "go-gpt3"
	defaultTimeoutSeconds = 30
)

func getEngineURL(engine string) string {
	return fmt.Sprintf("%s/engines/%s/completions", defaultBaseURL, engine)
}

// A Client is an API client to communicate with the OpenAI gpt-3 APIs
type Client interface {
	// Engines lists the currently available engines, and provides basic information about each
	// option such as the owner and availability.
	Engines(ctx context.Context) (*EnginesResponse, error)

	// Engine retrieves an engine instance, providing basic information about the engine such
	// as the owner and availability.
	Engine(ctx context.Context, engine string) (*EngineObject, error)

	// Completion creates a completion with the default engine. This is the main endpoint of the API.
	// Returns new text as well as, if requested, the probabilities over each alternative token at
	// each position.
	Completion(ctx context.Context, request CompletionRequest) (*CompletionResponse, error)

	// CompletionWithEngine creates a completion with the specified engine. This is the main endpoint
	// of the API. Returns new text as well as, if requested, the probabilities over each alternative
	// token at each position.
	CompletionWithEngine(ctx context.Context, engine string, request CompletionRequest) (*CompletionResponse, error)

	// Search performs a semantic search over a list of documents with the default engine.
	Search(ctx context.Context, request SearchRequest) (*SearchResponse, error)

	// SearchWithEngine performs a semantic search over a list of documents with the specified engine.
	SearchWithEngine(ctx context.Context, engine string, request SearchRequest) (*SearchResponse, error)
}

type client struct {
	baseURL       string
	apiKey        string
	userAgent     string
	httpClient    *http.Client
	defaultEngine string
}

// NewClient returns a new OpenAI GPT-3 API client. An apiKey is required to use the client
func NewClient(apiKey string, options ...ClientOption) Client {
	httpClient := &http.Client{
		Timeout: time.Duration(defaultTimeoutSeconds * time.Second),
	}

	c := &client{
		userAgent:     defaultUserAgent,
		apiKey:        apiKey,
		baseURL:       defaultBaseURL,
		httpClient:    httpClient,
		defaultEngine: DefaultEngine,
	}
	for _, o := range options {
		o(c)
	}
	return c
}

func (c *client) Engines(ctx context.Context) (*EnginesResponse, error) {
	req, err := c.newRequest(ctx, "GET", "/engines", nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	output := new(EnginesResponse)
	if err := getResponseObject(resp, output); err != nil {
		return nil, err
	}
	return output, nil
}

func (c *client) Engine(ctx context.Context, engine string) (*EngineObject, error) {
	req, err := c.newRequest(ctx, "GET", fmt.Sprintf("/engines/%s", engine), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	output := new(EngineObject)
	if err := getResponseObject(resp, output); err != nil {
		return nil, err
	}
	return output, nil
}

func (c *client) Completion(ctx context.Context, request CompletionRequest) (*CompletionResponse, error) {
	return c.CompletionWithEngine(ctx, c.defaultEngine, request)
}

func (c *client) CompletionWithEngine(ctx context.Context, engine string, request CompletionRequest) (*CompletionResponse, error) {
	req, err := c.newRequest(ctx, "POST", fmt.Sprintf("/engines/%s/completions", engine), request)
	if err != nil {
		return nil, err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	output := new(CompletionResponse)
	if err := getResponseObject(resp, output); err != nil {
		return nil, err
	}
	return output, nil
}

func (c *client) Search(ctx context.Context, request SearchRequest) (*SearchResponse, error) {
	return c.SearchWithEngine(ctx, c.defaultEngine, request)
}

func (c *client) SearchWithEngine(ctx context.Context, engine string, request SearchRequest) (*SearchResponse, error) {
	req, err := c.newRequest(ctx, "POST", fmt.Sprintf("/engines/%s/search", engine), request)
	if err != nil {
		return nil, err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	output := new(SearchResponse)
	if err := getResponseObject(resp, output); err != nil {
		return nil, err
	}
	return output, nil
}

func getResponseObject(rsp *http.Response, v interface{}) error {
	defer rsp.Body.Close()
	return json.NewDecoder(rsp.Body).Decode(v)
}

func jsonBodyReader(body interface{}) (io.Reader, error) {
	if body == nil {
		return bytes.NewBuffer(nil), nil
	}
	raw, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed encoding json: %w", err)
	}
	return bytes.NewBuffer(raw), nil
}

func (c *client) newRequest(ctx context.Context, method, path string, payload interface{}) (*http.Request, error) {
	bodyReader, err := jsonBodyReader(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, method, path, bodyReader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))
	return req, nil
}
