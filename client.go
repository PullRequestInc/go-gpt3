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

const (
	DEFAULT_BASE_URL   = "https://api.openai.com/v1"
	DEFAULT_USER_AGENT = "gpt3-go"
	DEFAULT_TIMEOUT    = 30
)

var dataPrefix = []byte("data: ")
var streamTerminationPrefix = []byte("[DONE]")

type Client interface {
	Models(ctx context.Context) (*ModelsResponse, error)
	Model(ctx context.Context, model string) (*ModelObject, error)
	Completion(ctx context.Context, request CompletionRequest) (*CompletionResponse, error)
	CompletionStream(ctx context.Context, request CompletionRequest, onData func(*CompletionResponse)) error
	Edits(ctx context.Context, request EditsRequest) (*EditsResponse, error)
	Embeddings(ctx context.Context, request EmbeddingsRequest) (*EmbeddingsResponse, error)
	Files(ctx context.Context) (*FilesResponse, error)
	UploadFile(ctx context.Context, request UploadFileRequest) (*FileObject, error)
	DeleteFile(ctx context.Context, fileID string) (*DeleteFileResponse, error)
	File(ctx context.Context, fileID string) (*FileObject, error)
	FileContent(ctx context.Context, fileID string) ([]byte, error)
	CreateFineTune(ctx context.Context, request CreateFineTuneRequest) (*FineTuneObject, error)
	ListFineTunes(ctx context.Context) (*ListFineTunesResponse, error)
	FineTune(ctx context.Context, fineTuneID string) (*FineTuneObject, error)
	CancelFineTune(ctx context.Context, fineTuneID string) (*FineTuneObject, error)
	FineTuneEvents(ctx context.Context, request FineTuneEventsRequest) (*FineTuneEventsResponse, error)
	FineTuneStreamEvents(ctx context.Context, request FineTuneEventsRequest, onData func(*FineTuneEvent)) error
	DeleteFineTuneModel(ctx context.Context, modelID string) (*DeleteFineTuneModelResponse, error)

	// Deprecated
	CompletionWithEngine(ctx context.Context, engine string, request CompletionRequest) (*CompletionResponse, error)
	CompletionStreamWithEngine(ctx context.Context, engine string, request CompletionRequest, onData func(*CompletionResponse)) error
}

type client struct {
	baseURL      string
	apiKey       string
	orgID        string
	userAgent    string
	httpClient   *http.Client
	defaultModel string
}

func NewClient(apiKey string, options ...ClientOption) (Client, error) {
	c := &client{
		baseURL:      DEFAULT_BASE_URL,
		apiKey:       apiKey,
		orgID:        "",
		userAgent:    DEFAULT_USER_AGENT,
		httpClient:   &http.Client{Timeout: time.Duration(DEFAULT_TIMEOUT) * time.Second},
		defaultModel: DavinciModel,
	}

	for _, option := range options {
		if err := option(c); err != nil {
			return nil, err
		}
	}

	return c, nil
}

func (c *client) newRequest(ctx context.Context, method, path string, payload interface{}) (*http.Request, error) {
	bodyReader, err := jsonBodyReader(payload)
	if err != nil {
		return nil, err
	}
	url := c.baseURL + path
	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))
	req.Header.Set("User-Agent", c.userAgent)
	if len(c.orgID) > 0 {
		req.Header.Set("OpenAI-Organization", c.orgID)
	}
	return req, nil
}

func (c *client) performRequest(req *http.Request) (*http.Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return resp, checkForSuccess(resp)
}

func checkForSuccess(resp *http.Response) error {
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read from body: %w", err)
	}
	var result APIErrorResponse
	if err := json.Unmarshal(data, &result); err != nil {
		// if we can't decode the json error then create an unexpected error
		apiError := APIError{
			StatusCode: resp.StatusCode,
			Type:       "Unexpected",
			Message:    string(data),
		}
		return apiError
	}
	result.Error.StatusCode = resp.StatusCode
	return result.Error
}

func getResponseObject(rsp *http.Response, v interface{}) error {
	defer rsp.Body.Close()
	if err := json.NewDecoder(rsp.Body).Decode(v); err != nil {
		return fmt.Errorf("invalid json response: %w", err)
	}
	return nil
}

func jsonBodyReader(body interface{}) (io.Reader, error) {
	if body == nil {
		// the body is allowed to be nil so we return an empty buffer
		return bytes.NewBuffer(nil), nil
	}
	raw, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed encoding json: %w", err)
	}
	return bytes.NewBuffer(raw), nil
}
