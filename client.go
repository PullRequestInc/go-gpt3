package gpt3

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	DEFAULT_BASE_URL   = "https://api.openai.com/v1"
	DEFAULT_USER_AGENT = "gpt3-go"
	DEFAULT_TIMEOUT    = 30
)

var dataPrefix = []byte("data: ")
var doneSequence = []byte("[DONE]")

type Client interface {
	Models(ctx context.Context) (*ModelsResponse, error)
	Model(ctx context.Context, model string) (*ModelObject, error)
	Completion(ctx context.Context, request CompletionRequest) (*CompletionResponse, error)
	Edits(ctx context.Context, request EditsRequest) (*EditsResponse, error)
	Embeddings(ctx context.Context, request EmbeddingsRequest) (*EmbeddingsResponse, error)
	Files(ctx context.Context) (*FilesResponse, error)
	UploadFile(ctx context.Context, request UploadFileRequest) (*FileObject, error)
	DeleteFile(ctx context.Context, fileID string) (*DeleteFileResponse, error)
	File(ctx context.Context, fileID string) (*FileObject, error)
	FileContent(ctx context.Context, fileID string) ([]byte, error)
	CreateFineTune(ctx context.Context, request CreateFineTuneRequest) (*FineTuneObject, error)
	FineTunes(ctx context.Context) (*FineTunesResponse, error)
	FineTune(ctx context.Context, request FineTuneRequest) (*FineTuneObject, error)
	CancelFineTune(ctx context.Context, request FineTuneRequest) (*FineTuneObject, error)
	FineTuneEvents(ctx context.Context, request FineTuneEventsRequest) (*FineTuneEventsResponse, error)
	FineTuneStreamEvents(ctx context.Context, request FineTuneEventsRequest, onData func(*FineTuneEvent)) error
	DeleteFineTuneModel(ctx context.Context, request DeleteFineTuneModelRequest) (*DeleteFineTuneModelResponse, error)
}

type client struct {
	BaseURL      string
	APIKey       string
	OrgID        string
	UserAgent    string
	HttpClient   *http.Client
	DefaultModel string
}

func NewClient(apiKey string, options ...ClientOption) (Client, error) {
	c := &client{
		BaseURL:      DEFAULT_BASE_URL,
		APIKey:       apiKey,
		OrgID:        "",
		UserAgent:    DEFAULT_USER_AGENT,
		HttpClient:   &http.Client{Timeout: time.Duration(DEFAULT_TIMEOUT) * time.Second},
		DefaultModel: DavinciModel,
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
	url := c.BaseURL + path
	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.APIKey))
	req.Header.Set("User-Agent", c.UserAgent)
	if len(c.OrgID) > 0 {
		req.Header.Set("OpenAI-Organization", c.OrgID)
	}
	return req, nil
}

func (c *client) performRequest(req *http.Request) (*http.Response, error) {
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, checkForSuccess(resp)
}

func checkForSuccess(resp *http.Response) error {
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read from body: %w", err)
	}
	// If the content-type is not json, then we can't decode it, so just return the
	// response as is.
	// See: https://beta.openai.com/docs/api-reference/files/retrieve-content
	if resp.Header.Get("Content-Type") != "application/json" {
		return nil
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
	decoder := json.NewDecoder(rsp.Body)
	if err := decoder.Decode(v); err != nil {
		return fmt.Errorf("invalid json response: %w", err)
	}
	return nil
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
