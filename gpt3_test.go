package gpt3_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/PullRequestInc/go-gpt3"
	fakes "github.com/PullRequestInc/go-gpt3/go-gpt3fakes"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 net/http.RoundTripper

func TestInitNewClient(t *testing.T) {
	client := gpt3.NewClient("test-key")
	assert.NotNil(t, client)
}

func fakeHttpClient() (*fakes.FakeRoundTripper, *http.Client) {
	rt := &fakes.FakeRoundTripper{}
	return rt, &http.Client{
		Transport: rt,
	}
}

func TestRequestCreationFails(t *testing.T) {
	ctx := context.Background()
	rt, httpClient := fakeHttpClient()
	client := gpt3.NewClient("test-key", gpt3.WithHTTPClient(httpClient))
	rt.RoundTripReturns(nil, errors.New("request error"))

	type testCase struct {
		name        string
		apiCall     func() (interface{}, error)
		errorString string
	}

	testCases := []testCase{
		{
			"Engines",
			func() (interface{}, error) {
				return client.Engines(ctx)
			},
			"Get \"https://api.openai.com/v1/engines\": request error",
		},
		{
			"Engine",
			func() (interface{}, error) {
				return client.Engine(ctx, gpt3.DefaultEngine)
			},
			"Get \"https://api.openai.com/v1/engines/davinci\": request error",
		},
		{
			"ChatCompletion",
			func() (interface{}, error) {
				return client.ChatCompletion(ctx, gpt3.ChatCompletionRequest{})
			},
			"Post \"https://api.openai.com/v1/chat/completions\": request error",
		},
		{
			"Completion",
			func() (interface{}, error) {
				return client.Completion(ctx, gpt3.CompletionRequest{})
			},
			"Post \"https://api.openai.com/v1/engines/davinci/completions\": request error",
		},
		{
			"CompletionStream",
			func() (interface{}, error) {
				var rsp *gpt3.CompletionResponse
				onData := func(data *gpt3.CompletionResponse) {
					rsp = data
				}
				return rsp, client.CompletionStream(ctx, gpt3.CompletionRequest{}, onData)
			},
			"Post \"https://api.openai.com/v1/engines/davinci/completions\": request error",
		},
		{
			"CompletionWithEngine",
			func() (interface{}, error) {
				return client.CompletionWithEngine(ctx, gpt3.AdaEngine, gpt3.CompletionRequest{})
			},
			"Post \"https://api.openai.com/v1/engines/ada/completions\": request error",
		},
		{
			"CompletionStreamWithEngine",
			func() (interface{}, error) {
				var rsp *gpt3.CompletionResponse
				onData := func(data *gpt3.CompletionResponse) {
					rsp = data
				}
				return rsp, client.CompletionStreamWithEngine(ctx, gpt3.AdaEngine, gpt3.CompletionRequest{}, onData)
			},
			"Post \"https://api.openai.com/v1/engines/ada/completions\": request error",
		},
		{
			"Edits",
			func() (interface{}, error) {
				return client.Edits(ctx, gpt3.EditsRequest{})
			},
			"Post \"https://api.openai.com/v1/edits\": request error",
		},
		{
			"Search",
			func() (interface{}, error) {
				return client.Search(ctx, gpt3.SearchRequest{})
			},
			"Post \"https://api.openai.com/v1/engines/davinci/search\": request error",
		},
		{
			"SearchWithEngine",
			func() (interface{}, error) {
				return client.SearchWithEngine(ctx, gpt3.AdaEngine, gpt3.SearchRequest{})
			},
			"Post \"https://api.openai.com/v1/engines/ada/search\": request error",
		},
		{
			"Embeddings",
			func() (interface{}, error) {
				return client.Embeddings(ctx, gpt3.EmbeddingsRequest{})
			},
			"Post \"https://api.openai.com/v1/embeddings\": request error",
		},
		{
			"Moderation",
			func() (interface{}, error) {
				return client.Moderation(ctx, gpt3.ModerationRequest{})
			},
			"Post \"https://api.openai.com/v1/moderations\": request error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rsp, err := tc.apiCall()
			assert.EqualError(t, err, tc.errorString)
			assert.Nil(t, rsp)
		})
	}
}

type errReader int

func (errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("read error")
}

func TestResponses(t *testing.T) {
	ctx := context.Background()
	rt, httpClient := fakeHttpClient()
	client := gpt3.NewClient("test-key", gpt3.WithHTTPClient(httpClient))

	type testCase struct {
		name           string
		apiCall        func() (interface{}, error)
		responseObject interface{}
	}

	testCases := []testCase{
		{
			"Engines",
			func() (interface{}, error) {
				return client.Engines(ctx)
			},
			&gpt3.EnginesResponse{
				Data: []gpt3.EngineObject{
					{
						ID:     "123",
						Object: "list",
						Owner:  "owner",
						Ready:  true,
					},
				},
			},
		},
		{
			"Engine",
			func() (interface{}, error) {
				return client.Engine(ctx, gpt3.DefaultEngine)
			},
			&gpt3.EngineObject{
				ID:     "123",
				Object: "list",
				Owner:  "owner",
				Ready:  true,
			},
		},
		{
			"ChatCompletion",
			func() (interface{}, error) {
				return client.ChatCompletion(ctx, gpt3.ChatCompletionRequest{})
			},
			&gpt3.ChatCompletionResponse{
				ID:      "chatcmpl-123",
				Object:  "messages",
				Created: 123456789,
				Model:   "gpt-3.5-turbo",
				Choices: []gpt3.ChatCompletionResponseChoice{
					{
						Index:        0,
						FinishReason: "stop",
						Message: gpt3.ChatCompletionResponseMessage{
							Role:    "assistant",
							Content: "output",
						},
					},
				},
			},
		},
		{
			"ChatCompletionWithFunctionCall",
			func() (interface{}, error) {
				return client.ChatCompletion(ctx, gpt3.ChatCompletionRequest{})
			},
			&gpt3.ChatCompletionResponse{
				ID:      "chatcmpl-123",
				Object:  "messages",
				Created: 123456789,
				Model:   "gpt-3.5-turbo-0613",
				Choices: []gpt3.ChatCompletionResponseChoice{
					{
						Index:        0,
						FinishReason: "function_call",
						Message: gpt3.ChatCompletionResponseMessage{
							Role:    "assistant",
							Content: "",
							FunctionCall: &gpt3.Function{
								Name:      "get_current_weather",
								Arguments: `"{"location": "Boston, MA"}"`,
							},
						},
					},
				},
			},
		},
		{
			"Completion",
			func() (interface{}, error) {
				return client.Completion(ctx, gpt3.CompletionRequest{})
			},
			&gpt3.CompletionResponse{
				ID:      "123",
				Object:  "list",
				Created: 123456789,
				Model:   "davinci-12",
				Choices: []gpt3.CompletionResponseChoice{
					{
						Text:         "output",
						FinishReason: "stop",
					},
				},
			},
		},
		{
			"CompletionStream",
			func() (interface{}, error) {
				var rsp *gpt3.CompletionResponse
				onData := func(data *gpt3.CompletionResponse) {
					rsp = data
				}
				return rsp, client.CompletionStream(ctx, gpt3.CompletionRequest{}, onData)
			},
			nil, // streaming responses are tested separately
		},
		{
			"CompletionWithEngine",
			func() (interface{}, error) {
				return client.CompletionWithEngine(ctx, gpt3.AdaEngine, gpt3.CompletionRequest{})
			},
			&gpt3.CompletionResponse{
				ID:      "123",
				Object:  "list",
				Created: 123456789,
				Model:   "davinci-12",
				Choices: []gpt3.CompletionResponseChoice{
					{
						Text:         "output",
						FinishReason: "stop",
					},
				},
			},
		},
		{
			"CompletionStreamWithEngine",
			func() (interface{}, error) {
				var rsp *gpt3.CompletionResponse
				onData := func(data *gpt3.CompletionResponse) {
					rsp = data
				}
				return rsp, client.CompletionStreamWithEngine(ctx, gpt3.AdaEngine, gpt3.CompletionRequest{}, onData)
			},
			nil, // streaming responses are tested separately
		},
		{
			"Search",
			func() (interface{}, error) {
				return client.Search(ctx, gpt3.SearchRequest{})
			},
			&gpt3.SearchResponse{
				Data: []gpt3.SearchData{
					{
						Document: 1,
						Object:   "search_result",
						Score:    40.312,
					},
				},
			},
		},
		{
			"SearchWithEngine",
			func() (interface{}, error) {
				return client.SearchWithEngine(ctx, gpt3.AdaEngine, gpt3.SearchRequest{})
			},
			&gpt3.SearchResponse{
				Data: []gpt3.SearchData{
					{
						Document: 1,
						Object:   "search_result",
						Score:    40.312,
					},
				},
			},
		},
		{
			"Embeddings",
			func() (interface{}, error) {
				return client.Embeddings(ctx, gpt3.EmbeddingsRequest{})
			},
			&gpt3.EmbeddingsResponse{
				Object: "list",
				Data: []gpt3.EmbeddingsResult{{
					Object:    "object",
					Embedding: []float64{0.1, 0.2, 0.3},
					Index:     0,
				}},
				Usage: gpt3.EmbeddingsUsage{
					PromptTokens: 1,
					TotalTokens:  2,
				},
			},
		},
		{
			"Moderation",
			func() (interface{}, error) {
				return client.Moderation(ctx, gpt3.ModerationRequest{})
			},
			&gpt3.ModerationResponse{
				ID:    "123",
				Model: "text-moderation-001",
				Results: []gpt3.ModerationResult{{
					Flagged: false,
					Categories: gpt3.ModerationCategoryResult{
						Hate:            false,
						HateThreatening: false,
						SelfHarm:        false,
						Sexual:          false,
						SexualMinors:    false,
						Violence:        false,
						ViolenceGraphic: false,
					},
					CategoryScores: gpt3.ModerationCategoryScores{
						Hate:            0.22714105248451233,
						HateThreatening: 0.22714105248451233,
						SelfHarm:        0.005232391878962517,
						Sexual:          0.01407341007143259,
						SexualMinors:    0.0038522258400917053,
						Violence:        0.009223177433013916,
						ViolenceGraphic: 0.036865197122097015,
					},
				}},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Run("bad status codes", func(t *testing.T) {
				for _, code := range []int{400, 401, 404, 422, 500} {
					// first mock with error with body failure
					mockResponse := &http.Response{
						StatusCode: code,
						Body:       ioutil.NopCloser(errReader(0)),
					}

					rt.RoundTripReturns(mockResponse, nil)
					rsp, err := tc.apiCall()
					assert.Nil(t, rsp)
					assert.EqualError(t, err, "failed to read from body: read error")

					// then mock with an unknown error string
					mockResponse = &http.Response{
						StatusCode: code,
						Body:       ioutil.NopCloser(bytes.NewBufferString("unknown error")),
					}

					rt.RoundTripReturns(mockResponse, nil)
					rsp, err = tc.apiCall()
					assert.Nil(t, rsp)
					assert.EqualError(t, err, fmt.Sprintf("[%d:Unexpected] unknown error", code))

					// then mock with an json APIErrorResponse
					apiErrorResponse := &gpt3.APIErrorResponse{
						Error: gpt3.APIError{
							Type:    "test_type",
							Message: "test message",
						},
					}

					data, err := json.Marshal(apiErrorResponse)
					assert.NoError(t, err)

					mockResponse = &http.Response{
						StatusCode: code,
						Body:       ioutil.NopCloser(bytes.NewBuffer(data)),
					}

					rt.RoundTripReturns(mockResponse, nil)
					rsp, err = tc.apiCall()
					assert.Nil(t, rsp)
					assert.EqualError(t, err, fmt.Sprintf("[%d:test_type] test message", code))
					apiErrorResponse.Error.StatusCode = code
					assert.Equal(t, apiErrorResponse.Error, err)
				}
			})
			t.Run("success code json decode failure", func(t *testing.T) {
				mockResponse := &http.Response{
					StatusCode: 200,
					Body:       ioutil.NopCloser(bytes.NewBufferString("invalid json")),
				}

				rt.RoundTripReturns(mockResponse, nil)

				rsp, err := tc.apiCall()
				assert.Error(t, err, "invalid json response: invalid character 'i' looking for beginning of value")
				assert.Nil(t, rsp)
			})
			// skip streaming/nil response objects here as those will be tested separately
			if tc.responseObject != nil {
				t.Run("successful response", func(t *testing.T) {
					data, err := json.Marshal(tc.responseObject)
					assert.NoError(t, err)

					mockResponse := &http.Response{
						StatusCode: 200,
						Body:       ioutil.NopCloser(bytes.NewBuffer(data)),
					}

					rt.RoundTripReturns(mockResponse, nil)

					rsp, err := tc.apiCall()
					assert.NoError(t, err)
					assert.Equal(t, tc.responseObject, rsp)
				})
			}
		})
	}
}

func TestRateLimitHeaders(t *testing.T) {
	/*
		These values are taken directly from the documentation at https://platform.openai.com/docs/guides/rate-limits/overview

		| FIELD | SAMPLE VALUE | DESCRIPTION |
		| ---- | ---- | ---- |
		| x-ratelimit-limit-requests | 60 | The maximum number of requests that are permitted before exhausting the rate limit. |
		| x-ratelimit-limit-tokens | 150000 | The maximum number of tokens that are permitted before exhausting the rate limit. |
		| x-ratelimit-remaining-requests | 59 | The remaining number of requests that are permitted before exhausting the rate limit. |
		| x-ratelimit-remaining-tokens | 149984 | The remaining number of tokens that are permitted before exhausting the rate limit. |
		| x-ratelimit-reset-requests | 1s | The time until the rate limit (based on requests) resets to its initial state. |
		| x-ratelimit-reset-tokens | 6m0s | The time until the rate limit (based on tokens) resets to its initial state. |
	*/

	header := make(http.Header)
	header.Add("x-ratelimit-limit-requests", "60")
	header.Add("x-ratelimit-limit-tokens", "150000")
	header.Add("x-ratelimit-remaining-requests", "59")
	header.Add("x-ratelimit-remaining-tokens", "149984")
	header.Add("x-ratelimit-reset-requests", "1s")
	header.Add("x-ratelimit-reset-tokens", "6m0s")

	rateLimitHeaders := gpt3.NewRateLimitHeadersFromResponse(&http.Response{Header: header})
	assert.Equal(t, 60, rateLimitHeaders.LimitRequests)
	assert.Equal(t, 150000, rateLimitHeaders.LimitTokens)
	assert.Equal(t, 59, rateLimitHeaders.RemainingRequests)
	assert.Equal(t, 149984, rateLimitHeaders.RemainingTokens)
	assert.Equal(t, 1*time.Second, rateLimitHeaders.ResetRequests)
	assert.Equal(t, 6*time.Minute, rateLimitHeaders.ResetTokens)
}

// TODO: add streaming response tests
