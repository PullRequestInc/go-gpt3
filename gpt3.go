package gpt3

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// Engine Types
const (
	TextAda001Engine     = "text-ada-001"
	TextBabbage001Engine = "text-babbage-001"
	TextCurie001Engine   = "text-curie-001"
	TextDavinci001Engine = "text-davinci-001"
	TextDavinci002Engine = "text-davinci-002"
	TextDavinci003Engine = "text-davinci-003"
	AdaEngine            = "ada"
	BabbageEngine        = "babbage"
	CurieEngine          = "curie"
	DavinciEngine        = "davinci"
	DefaultEngine        = DavinciEngine
)

type EmbeddingEngine string

const (
	TextSimilarityAda001      = "text-similarity-ada-001"
	TextSimilarityBabbage001  = "text-similarity-babbage-001"
	TextSimilarityCurie001    = "text-similarity-curie-001"
	TextSimilarityDavinci001  = "text-similarity-davinci-001"
	TextSearchAdaDoc001       = "text-search-ada-doc-001"
	TextSearchAdaQuery001     = "text-search-ada-query-001"
	TextSearchBabbageDoc001   = "text-search-babbage-doc-001"
	TextSearchBabbageQuery001 = "text-search-babbage-query-001"
	TextSearchCurieDoc001     = "text-search-curie-doc-001"
	TextSearchCurieQuery001   = "text-search-curie-query-001"
	TextSearchDavinciDoc001   = "text-search-davinci-doc-001"
	TextSearchDavinciQuery001 = "text-search-davinci-query-001"
	CodeSearchAdaCode001      = "code-search-ada-code-001"
	CodeSearchAdaText001      = "code-search-ada-text-001"
	CodeSearchBabbageCode001  = "code-search-babbage-code-001"
	CodeSearchBabbageText001  = "code-search-babbage-text-001"
	TextEmbeddingAda002       = "text-embedding-ada-002"
)

const (
	defaultBaseURL        = "https://api.openai.com/v1"
	defaultUserAgent      = "go-gpt3"
	defaultTimeoutSeconds = 30
)

// A Client is an API client to communicate with the OpenAI gpt-3 APIs
type Client interface {
	// Deprecated: Engines lists the currently available engines, and provides basic information about each
	// option such as the owner and availability.
	Engines(ctx context.Context) (*EnginesResponse, error)

	// Deprecated: Engine retrieves an engine instance, providing basic information about the engine such
	// as the owner and availability.
	Engine(ctx context.Context, engine string) (*EngineObject, error)

	// Models lists the currently available models, and provides basic information about each
	// option such as the owner and availability.
	Models(ctx context.Context) (*ModelsResponse, error)

	// Model retrieves a model instance, providing basic information about the model such
	// as the owner and permissioning.
	Model(ctx context.Context, model string) (*ModelObject, error)

	// Completion creates a completion with the default engine. This is the main endpoint of the API
	// which auto-completes based on the given prompt.
	Completion(ctx context.Context, request CompletionRequest) (*CompletionResponse, error)

	// CompletionStream creates a completion with the default engine and streams the results through
	// multiple calls to onData.
	CompletionStream(ctx context.Context, request CompletionRequest, onData func(*CompletionResponse)) error

	// CompletionWithEngine is the same as Completion except allows overriding the default engine on the client.
	CompletionWithEngine(ctx context.Context, engine string, request CompletionRequest) (*CompletionResponse, error)

	// CompletionStreamWithEngine is the same as CompletionStream except allows overriding the default engine on the client.
	CompletionStreamWithEngine(ctx context.Context, engine string, request CompletionRequest, onData func(*CompletionResponse)) error

	// Given a prompt and an instruction, the model will return an edited version of the prompt.
	Edits(ctx context.Context, request EditsRequest) (*EditsResponse, error)

	// Search performs a semantic search over a list of documents with the default engine.
	Search(ctx context.Context, request SearchRequest) (*SearchResponse, error)

	// SearchWithEngine performs a semantic search over a list of documents with the specified engine.
	SearchWithEngine(ctx context.Context, engine string, request SearchRequest) (*SearchResponse, error)

	// Embeddings returns an embedding using the provided request.
	Embeddings(ctx context.Context, request EmbeddingsRequest) (*EmbeddingsResponse, error)

	// Files lists the files that belong to the user's organization.
	Files(ctx context.Context) (*FilesResponse, error)

	// UploadFile uploads a file that contains document(s) to be used across various endpoints.
	UploadFile(ctx context.Context, request UploadFileRequest) (*FileObject, error)

	// DeleteFile deletes a file from the user's organization.
	DeleteFile(ctx context.Context, fileID string) (*DeleteFileResponse, error)

	// File retrieves a file from the user's organization.
	File(ctx context.Context, fileID string) (*FileObject, error)

	// FileContent retrieves the content of a file from the user's organization
	// and returns it as raw bytes.
	FileContent(ctx context.Context, fileID string) ([]byte, error)

	// CreateFineTune creates a job that fine-tunes a model on a dataset.
	CreateFineTune(ctx context.Context, request CreateFineTuneRequest) (*FineTuneObject, error)

	// FineTunes lists the fine-tuning jobs that belong to the user's organization.
	FineTunes(ctx context.Context) (*FineTunesResponse, error)

	// FineTune retrieves a fine-tuning job from the user's organization.
	FineTune(ctx context.Context, request FineTuneRequest) (*FineTuneObject, error)

	// CancelFineTune cancels a fine-tuning job from the user's organization.
	CancelFineTune(ctx context.Context, request FineTuneRequest) (*FineTuneObject, error)

	// FineTuneEvents lists the events that belong to a fine-tuning job.
	FineTuneEvents(ctx context.Context, request FineTuneEventsRequest) (*FineTuneEventsResponse, error)

	// FineTuneStreamEvents streams the events that belong to a fine-tuning job.
	FineTuneStreamEvents(ctx context.Context, request FineTuneEventsRequest, onData func(*FineTuneEvent)) error

	// DeleteFineTuneModel deletes a fine-tuned model from the user's organization.
	DeleteFineTuneModel(ctx context.Context, request DeleteFineTuneModelRequest) (*DeleteFineTuneModelResponse, error)
}

type client struct {
	baseURL       string
	apiKey        string
	userAgent     string
	httpClient    *http.Client
	defaultEngine string
	idOrg         string
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
		idOrg:         "",
	}
	for _, o := range options {
		err := o(c)
		if err != nil {
			log.Fatalln(err)
		}
	}
	return c
}

// Deprecated: Engines lists the currently available engines, and provides basic information about each
// option such as the owner and availability.
//
// Use Models instead
func (c *client) Engines(ctx context.Context) (*EnginesResponse, error) {
	req, err := c.newRequest(ctx, "GET", "/engines", nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.performRequest(req)
	if err != nil {
		return nil, err
	}

	output := new(EnginesResponse)
	if err := getResponseObject(resp, output); err != nil {
		return nil, err
	}
	return output, nil
}

// Deprecated: Engine retrieves a single engine, providing basic information about the engine such
// as the owner and availability.
//
// Use Model instead
func (c *client) Engine(ctx context.Context, engine string) (*EngineObject, error) {
	req, err := c.newRequest(ctx, "GET", fmt.Sprintf("/engines/%s", engine), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.performRequest(req)
	if err != nil {
		return nil, err
	}

	output := new(EngineObject)
	if err := getResponseObject(resp, output); err != nil {
		return nil, err
	}
	return output, nil
}

// Models lists the currently available models, and provides basic information about each
// option such as the owner and permissioning.
func (c *client) Models(ctx context.Context) (*ModelsResponse, error) {
	req, err := c.newRequest(ctx, "GET", "/models", nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.performRequest(req)
	if err != nil {
		return nil, err
	}

	output := new(ModelsResponse)
	if err := getResponseObject(resp, output); err != nil {
		return nil, err
	}
	return output, nil
}

// Model retrieves a single model, providing basic information about the model such
// as the owner and permissioning.
func (c *client) Model(ctx context.Context, model string) (*ModelObject, error) {
	req, err := c.newRequest(ctx, "GET", fmt.Sprintf("/models/%s", model), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.performRequest(req)
	if err != nil {
		return nil, err
	}

	output := new(ModelObject)
	if err := getResponseObject(resp, output); err != nil {
		return nil, err
	}
	return output, nil
}

func (c *client) Completion(ctx context.Context, request CompletionRequest) (*CompletionResponse, error) {
	return c.CompletionWithEngine(ctx, c.defaultEngine, request)
}

func (c *client) CompletionWithEngine(ctx context.Context, engine string, request CompletionRequest) (*CompletionResponse, error) {
	request.Stream = false
	req, err := c.newRequest(ctx, "POST", fmt.Sprintf("/engines/%s/completions", engine), request)
	if err != nil {
		return nil, err
	}
	resp, err := c.performRequest(req)
	if err != nil {
		return nil, err
	}

	output := new(CompletionResponse)
	if err := getResponseObject(resp, output); err != nil {
		return nil, err
	}
	return output, nil
}

func (c *client) CompletionStream(ctx context.Context, request CompletionRequest, onData func(*CompletionResponse)) error {
	return c.CompletionStreamWithEngine(ctx, c.defaultEngine, request, onData)
}

var dataPrefix = []byte("data: ")
var doneSequence = []byte("[DONE]")

func (c *client) CompletionStreamWithEngine(
	ctx context.Context,
	engine string,
	request CompletionRequest,
	onData func(*CompletionResponse),
) error {
	request.Stream = true
	req, err := c.newRequest(ctx, "POST", fmt.Sprintf("/engines/%s/completions", engine), request)
	if err != nil {
		return err
	}
	resp, err := c.performRequest(req)
	if err != nil {
		return err
	}

	reader := bufio.NewReader(resp.Body)
	defer resp.Body.Close()

	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			return err
		}
		// make sure there isn't any extra whitespace before or after
		line = bytes.TrimSpace(line)
		// the completion API only returns data events
		if !bytes.HasPrefix(line, dataPrefix) {
			continue
		}
		line = bytes.TrimPrefix(line, dataPrefix)

		// the stream is completed when terminated by [DONE]
		if bytes.HasPrefix(line, doneSequence) {
			break
		}
		output := new(CompletionResponse)
		if err := json.Unmarshal(line, output); err != nil {
			return fmt.Errorf("invalid json stream data: %v", err)
		}
		onData(output)
	}

	return nil
}

func (c *client) Edits(ctx context.Context, request EditsRequest) (*EditsResponse, error) {
	req, err := c.newRequest(ctx, "POST", "/edits", request)
	if err != nil {
		return nil, err
	}
	resp, err := c.performRequest(req)
	if err != nil {
		return nil, err
	}

	output := new(EditsResponse)
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
	resp, err := c.performRequest(req)
	if err != nil {
		return nil, err
	}

	output := new(SearchResponse)
	if err := getResponseObject(resp, output); err != nil {
		return nil, err
	}
	return output, nil
}

// Embeddings creates text embeddings for a supplied slice of inputs with a provided model.
//
// See: https://beta.openai.com/docs/api-reference/embeddings
func (c *client) Embeddings(ctx context.Context, request EmbeddingsRequest) (*EmbeddingsResponse, error) {
	req, err := c.newRequest(ctx, "POST", "/embeddings", request)
	if err != nil {
		return nil, err
	}
	resp, err := c.performRequest(req)
	if err != nil {
		return nil, err
	}

	output := EmbeddingsResponse{}
	if err := getResponseObject(resp, &output); err != nil {
		return nil, err
	}
	return &output, nil
}

// Files lists the files that belong to the user's organization.
//
// See: https://beta.openai.com/docs/api-reference/files/list
func (c *client) Files(ctx context.Context) (*FilesResponse, error) {
	req, err := c.newRequest(ctx, "GET", "/files", nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.performRequest(req)
	if err != nil {
		return nil, err
	}

	output := FilesResponse{}
	if err := getResponseObject(resp, &output); err != nil {
		return nil, err
	}
	return &output, nil
}

// UploadFile uploads a file that contains document(s) to be used across various endpoints.
//
// See: https://beta.openai.com/docs/api-reference/files/upload
func (c *client) UploadFile(ctx context.Context, request UploadFileRequest) (*FileObject, error) {
	req, err := c.newRequest(ctx, "POST", "/files", request)
	if err != nil {
		return nil, err
	}
	resp, err := c.performRequest(req)
	if err != nil {
		return nil, err
	}

	output := FileObject{}
	if err := getResponseObject(resp, &output); err != nil {
		return nil, err
	}
	return &output, nil
}

// DeleteFile deletes a file that contains document(s) to be used across various endpoints.
//
// See: https://beta.openai.com/docs/api-reference/files/delete
func (c *client) DeleteFile(ctx context.Context, fileID string) (*DeleteFileResponse, error) {
	req, err := c.newRequest(ctx, "DELETE", fmt.Sprintf("/files/%s", fileID), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.performRequest(req)
	if err != nil {
		return nil, err
	}

	output := DeleteFileResponse{}
	if err := getResponseObject(resp, &output); err != nil {
		return nil, err
	}
	return &output, nil
}

// File retrieves a file that contains document(s) to be used across various endpoints.
//
// See: https://beta.openai.com/docs/api-reference/files/retrieve
func (c *client) File(ctx context.Context, fileID string) (*FileObject, error) {
	req, err := c.newRequest(ctx, "GET", fmt.Sprintf("/files/%s", fileID), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.performRequest(req)
	if err != nil {
		return nil, err
	}

	output := FileObject{}
	if err := getResponseObject(resp, &output); err != nil {
		return nil, err
	}
	return &output, nil
}

// FileContent retrieves the content of a file that contains document(s) to be used across various endpoints.
//
// See: https://beta.openai.com/docs/api-reference/files/retrieve-content
func (c *client) FileContent(ctx context.Context, fileID string) ([]byte, error) {
	req, err := c.newRequest(ctx, "GET", fmt.Sprintf("/files/%s/content", fileID), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.performRequest(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

// CreateFineTune creates a job that fine-tunes a model on a dataset.
//
// See: https://beta.openai.com/docs/api-reference/fine-tunes/create
func (c *client) CreateFineTune(ctx context.Context, request CreateFineTuneRequest) (*FineTuneObject, error) {
	req, err := c.newRequest(ctx, "POST", "/fine-tunes", request)
	if err != nil {
		return nil, err
	}
	resp, err := c.performRequest(req)
	if err != nil {
		return nil, err
	}

	output := FineTuneObject{}
	if err := getResponseObject(resp, &output); err != nil {
		return nil, err
	}
	return &output, nil
}

// FineTunes lists the fine-tuning jobs that belong to the user's organization.
//
// See: https://beta.openai.com/docs/api-reference/fine-tunes/list
func (c *client) FineTunes(ctx context.Context) (*FineTunesResponse, error) {
	req, err := c.newRequest(ctx, "GET", "/fine-tunes", nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.performRequest(req)
	if err != nil {
		return nil, err
	}

	output := FineTunesResponse{}
	if err := getResponseObject(resp, &output); err != nil {
		return nil, err
	}
	return &output, nil
}

// FineTune retrieves a fine-tuning job from the user's organization.
//
// See: https://beta.openai.com/docs/api-reference/fine-tunes/retrieve
func (c *client) FineTune(ctx context.Context, request FineTuneRequest) (*FineTuneObject, error) {
	req, err := c.newRequest(ctx, "GET", fmt.Sprintf("/fine-tunes/%s", request.FineTuneID), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.performRequest(req)
	if err != nil {
		return nil, err
	}

	output := FineTuneObject{}
	if err := getResponseObject(resp, &output); err != nil {
		return nil, err
	}
	return &output, nil
}

// CancelFineTune cancels a fine-tuning job from the user's organization.
//
// See: https://beta.openai.com/docs/api-reference/fine-tunes/cancel
func (c *client) CancelFineTune(ctx context.Context, request FineTuneRequest) (*FineTuneObject, error) {
	req, err := c.newRequest(ctx, "POST", fmt.Sprintf("/fine-tunes/%s/cancel", request.FineTuneID), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.performRequest(req)
	if err != nil {
		return nil, err
	}

	output := FineTuneObject{}
	if err := getResponseObject(resp, &output); err != nil {
		return nil, err
	}
	return &output, nil
}

// FineTuneEvents lists the events that belong to a fine-tuning job.
//
// See: https://beta.openai.com/docs/api-reference/fine-tunes/events
func (c *client) FineTuneEvents(ctx context.Context, request FineTuneEventsRequest) (*FineTuneEventsResponse, error) {
	request.Stream = false
	req, err := c.newRequest(ctx, "GET", fmt.Sprintf("/fine-tunes/%s/events", request.FineTuneID), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.performRequest(req)
	if err != nil {
		return nil, err
	}

	output := FineTuneEventsResponse{}
	if err := getResponseObject(resp, &output); err != nil {
		return nil, err
	}
	return &output, nil
}

// FineTuneStreamEvents streams the events that belong to a fine-tuning job.
//
// See: https://beta.openai.com/docs/api-reference/fine-tunes/events
func (c *client) FineTuneStreamEvents(ctx context.Context, request FineTuneEventsRequest, onData func(*FineTuneEvent)) error {
	request.Stream = true
	req, err := c.newRequest(ctx, "GET", fmt.Sprintf("/fine-tunes/%s/events", request.FineTuneID), nil)
	if err != nil {
		return err
	}
	resp, err := c.performRequest(req)
	if err != nil {
		return err
	}

	reader := bufio.NewReader(resp.Body)
	defer resp.Body.Close()

	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			return err
		}

		line = bytes.TrimSpace(line)
		if !bytes.HasPrefix(line, dataPrefix) {
			continue
		}
		line = bytes.TrimPrefix(line, dataPrefix)

		if bytes.HasPrefix(line, doneSequence) {
			break
		}
		output := new(FineTuneEvent)
		if err := json.Unmarshal(line, output); err != nil {
			return fmt.Errorf("invalid json stream data: %v", err)
		}
		onData(output)
	}
	return nil
}

// DeleteFineTuneModel deletes a fine-tuned model from the user's organization.
//
// See: https://beta.openai.com/docs/api-reference/fine-tunes/delete-model
func (c *client) DeleteFineTuneModel(ctx context.Context, request DeleteFineTuneModelRequest) (*DeleteFineTuneModelResponse, error) {
	req, err := c.newRequest(ctx, "DELETE", fmt.Sprintf("/fine-tunes/%s/model", request.Model), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.performRequest(req)
	if err != nil {
		return nil, err
	}

	output := DeleteFineTuneModelResponse{}
	if err := getResponseObject(resp, &output); err != nil {
		return nil, err
	}
	return &output, nil
}

func (c *client) performRequest(req *http.Request) (*http.Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, checkForSuccess(resp)
}

// checkForSuccess checks the response to see if it was successful. If it was
// not successful, it returns the error that was returned by the API.
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
	if len(c.idOrg) > 0 {
		req.Header.Set("OpenAI-Organization", c.idOrg)
	}
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))
	return req, nil
}
