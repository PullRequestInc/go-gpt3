package gpt3_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/PullRequestInc/go-gpt3"
	fakes "github.com/PullRequestInc/go-gpt3/go-gpt3fakes"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 net/http.RoundTripper

func fakeHttpClient() (*fakes.FakeRoundTripper, *http.Client) {
	rt := &fakes.FakeRoundTripper{}
	return rt, &http.Client{
		Transport: rt,
	}
}

func TestRequestCreationFails(t *testing.T) {
	ctx := context.Background()
	rt, httpClient := fakeHttpClient()
	client, err := gpt3.NewClient("test-key", gpt3.WithHTTPClient(httpClient))
	assert.Nil(t, err)
	rt.RoundTripReturns(nil, errors.New("request error"))

	type testCase struct {
		name        string
		apiCall     func() (interface{}, error)
		errorString string
	}

	testCases := []testCase{
		{
			"Models",
			func() (interface{}, error) {
				return client.Models(ctx)
			},
			"Get \"https://api.openai.com/v1/models\": request error",
		},
		{
			"Model",
			func() (interface{}, error) {
				return client.Model(ctx, gpt3.DefaultModel)
			},
			"Get \"https://api.openai.com/v1/models/davinci\": request error",
		},
		{
			"Completion",
			func() (interface{}, error) {
				return client.Completion(ctx, gpt3.CompletionRequest{})
			},
			"Post \"https://api.openai.com/v1/completions\": request error",
		}, {
			"CompletionStream",
			func() (interface{}, error) {
				var rsp *gpt3.CompletionResponse
				onData := func(data *gpt3.CompletionResponse) {
					rsp = data
				}
				return rsp, client.CompletionStream(ctx, gpt3.CompletionRequest{}, onData)
			},
			"Post \"https://api.openai.com/v1/completions\": request error",
		}, {
			"Edits",
			func() (interface{}, error) {
				return client.Edits(ctx, gpt3.EditsRequest{})
			},
			"Post \"https://api.openai.com/v1/edits\": request error",
		}, {
			"Embeddings",
			func() (interface{}, error) {
				return client.Embeddings(ctx, gpt3.EmbeddingsRequest{})
			},
			"Post \"https://api.openai.com/v1/embeddings\": request error",
		}, {
			"Files",
			func() (interface{}, error) {
				return client.Files(ctx)
			},
			"Get \"https://api.openai.com/v1/files\": request error",
		}, {
			"UploadFile",
			func() (interface{}, error) {
				return client.UploadFile(ctx, gpt3.UploadFileRequest{})
			},
			"Post \"https://api.openai.com/v1/files\": request error",
		}, {
			"DeleteFile",
			func() (interface{}, error) {
				return client.DeleteFile(ctx, "file-id")
			},
			"Delete \"https://api.openai.com/v1/files/file-id\": request error",
		}, {
			"File",
			func() (interface{}, error) {
				return client.File(ctx, "file-id")
			},
			"Get \"https://api.openai.com/v1/files/file-id\": request error",
		}, {
			"FileContent",
			func() (interface{}, error) {
				return client.FileContent(ctx, "file-id")
			},
			"Get \"https://api.openai.com/v1/files/file-id/content\": request error",
		}, {
			"FineTunes",
			func() (interface{}, error) {
				return client.FineTunes(ctx)
			},
			"Get \"https://api.openai.com/v1/fine-tunes\": request error",
		}, {
			"FineTune",
			func() (interface{}, error) {
				return client.FineTune(ctx, "fine-tune-id")
			},
			"Get \"https://api.openai.com/v1/fine-tunes/fine-tune-id\": request error",
		}, {
			"CancelFineTune",
			func() (interface{}, error) {
				return client.CancelFineTune(ctx, "fine-tune-id")
			},
			"Post \"https://api.openai.com/v1/fine-tunes/fine-tune-id/cancel\": request error",
		}, {
			"FineTuneEvents",
			func() (interface{}, error) {
				return client.FineTuneEvents(ctx, gpt3.FineTuneEventsRequest{
					FineTuneID: "fine-tune-id",
				})
			},
			"Get \"https://api.openai.com/v1/fine-tunes/fine-tune-id/events\": request error",
		},
		{
			"FineTuneStreamEvents",
			func() (interface{}, error) {
				var rsp *gpt3.FineTuneEvent
				onData := func(data *gpt3.FineTuneEvent) {
					rsp = data
				}
				return rsp, client.FineTuneStreamEvents(ctx, gpt3.FineTuneEventsRequest{
					FineTuneID: "fine-tune-id",
				}, onData)
			},
			"Get \"https://api.openai.com/v1/fine-tunes/fine-tune-id/events\": request error",
		}, {
			"DeleteFineTuneModel",
			func() (interface{}, error) {
				return client.DeleteFineTuneModel(ctx, "model-id")
			},
			"Delete \"https://api.openai.com/v1/models/model-id\": request error",
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
	client, err := gpt3.NewClient("test-key", gpt3.WithHTTPClient(httpClient))
	assert.Nil(t, err)

	type testCase struct {
		name           string
		apiCall        func() (interface{}, error)
		responseObject interface{}
	}

	testCases := []testCase{
		{
			"Models",
			func() (interface{}, error) {
				return client.Models(ctx)
			},
			&gpt3.ModelsResponse{
				Data: []gpt3.ModelObject{
					{
						ID:          "123",
						Object:      "list",
						OwnedBy:     "organization-owner",
						Permissions: []string{},
					},
				},
			},
		},
		{
			"Model",
			func() (interface{}, error) {
				return client.Model(ctx, gpt3.DefaultModel)
			},
			&gpt3.ModelObject{
				ID:          "123",
				Object:      "list",
				OwnedBy:     "organization-owner",
				Permissions: []string{},
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
		}, {
			"CompletionStream",
			func() (interface{}, error) {
				var rsp *gpt3.CompletionResponse
				onData := func(data *gpt3.CompletionResponse) {
					rsp = data
				}
				return rsp, client.CompletionStream(ctx, gpt3.CompletionRequest{}, onData)
			},
			nil, // streaming responses are tested separately
		}, {
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
		}, {
			"Files",
			func() (interface{}, error) {
				return client.Files(ctx)
			},
			&gpt3.FilesResponse{
				Object: "list",
				Data: []gpt3.FileObject{
					{
						ID:        "123",
						Object:    "object",
						Bytes:     123,
						CreatedAt: 123456789,
						Filename:  "file.txt",
						Purpose:   "fine-tune",
					},
				},
			},
		}, {
			"File",
			func() (interface{}, error) {
				return client.File(ctx, "file-id")
			},
			&gpt3.FileObject{
				ID:        "123",
				Object:    "object",
				Bytes:     123,
				CreatedAt: 123456789,
				Filename:  "file.txt",
				Purpose:   "fine-tune",
			},
		}, {
			"UploadFile",
			func() (interface{}, error) {
				return client.UploadFile(ctx, gpt3.UploadFileRequest{
					File:    "file.jsonl",
					Purpose: "fine-tune",
				})

			},
			&gpt3.FileObject{
				ID:        "123",
				Object:    "object",
				Bytes:     123,
				CreatedAt: 123456789,
				Filename:  "file.txt",
				Purpose:   "fine-tune",
			},
		}, {
			"DeleteFile",
			func() (interface{}, error) {
				return client.DeleteFile(ctx, "file-id")
			},
			nil,
		}, {
			"FineTunes",
			func() (interface{}, error) {
				return client.FineTunes(ctx)
			},
			&gpt3.FineTunesResponse{
				Object: "list",
				Data: []gpt3.FineTuneObject{
					{
						ID:            "123",
						Object:        "object",
						Model:         "davinci-12",
						CreatedAt:     123456789,
						Events:        []gpt3.FineTuneEvent{},
						FineTuneModel: "davince:ft:123",
						Hyperparams: gpt3.FineTuneHyperparams{
							BatchSize:              1,
							LearningRateMultiplier: 1.0,
							NEpochs:                1,
							PromptLossWeight:       1.0,
						},
						OrganizationID: "org-id",
						ResultFiles: []gpt3.FileObject{
							{
								ID:        "123",
								Object:    "object",
								Bytes:     123,
								CreatedAt: 123456789,
								Filename:  "file.txt",
								Purpose:   "fine-tune",
							},
						},
						Status: "complete",
						ValidationFiles: []gpt3.FileObject{
							{
								ID:        "123",
								Object:    "object",
								Bytes:     123,
								CreatedAt: 123456789,
								Filename:  "file.txt",
								Purpose:   "fine-tune",
							},
						},
						TrainingFiles: []gpt3.FileObject{
							{
								ID:        "123",
								Object:    "object",
								Bytes:     123,
								CreatedAt: 123456789,
								Filename:  "file.txt",
								Purpose:   "fine-tune",
							},
						},
						UpdatedAt: 123456789,
					},
				},
			},
		}, {
			"FineTune",
			func() (interface{}, error) {
				return client.FineTune(ctx, "fine-tune-id")
			},
			&gpt3.FineTuneObject{
				ID:              "123",
				Object:          "object",
				Model:           "davinci-12",
				CreatedAt:       123456789,
				Events:          []gpt3.FineTuneEvent{},
				FineTuneModel:   "davince:ft:123",
				Hyperparams:     gpt3.FineTuneHyperparams{},
				OrganizationID:  "org-id",
				ResultFiles:     []gpt3.FileObject{},
				Status:          "complete",
				ValidationFiles: []gpt3.FileObject{},
				TrainingFiles:   []gpt3.FileObject{},
				UpdatedAt:       123456789,
			},
		}, {
			"CreateFineTune",
			func() (interface{}, error) {
				return client.CreateFineTune(ctx, gpt3.CreateFineTuneRequest{})
			},
			&gpt3.FineTuneObject{
				ID:              "123",
				Object:          "object",
				Model:           "davinci-12",
				CreatedAt:       123456789,
				Events:          []gpt3.FineTuneEvent{},
				FineTuneModel:   "davince:ft:123",
				Hyperparams:     gpt3.FineTuneHyperparams{},
				OrganizationID:  "org-id",
				ResultFiles:     []gpt3.FileObject{},
				Status:          "complete",
				ValidationFiles: []gpt3.FileObject{},
				TrainingFiles:   []gpt3.FileObject{},
				UpdatedAt:       123456789,
			},
		}, {
			"CancelFineTune",
			func() (interface{}, error) {
				return client.CancelFineTune(ctx, "fine-tune-id")
			},
			&gpt3.FineTuneObject{
				ID:              "123",
				Object:          "object",
				Model:           "davinci-12",
				CreatedAt:       123456789,
				Events:          []gpt3.FineTuneEvent{},
				FineTuneModel:   "davince:ft:123",
				Hyperparams:     gpt3.FineTuneHyperparams{},
				OrganizationID:  "org-id",
				ResultFiles:     []gpt3.FileObject{},
				Status:          "complete",
				ValidationFiles: []gpt3.FileObject{},
				TrainingFiles:   []gpt3.FileObject{},
				UpdatedAt:       123456789,
			},
		}, {
			"FineTuneEvents",
			func() (interface{}, error) {
				return client.FineTuneEvents(ctx, gpt3.FineTuneEventsRequest{
					FineTuneID: "fine-tune-id",
				})
			},
			&gpt3.FineTuneEventsResponse{
				Object: "list",
				Data: []gpt3.FineTuneEvent{
					{
						Object:    "object",
						CreatedAt: 123456789,
						Level:     "info",
						Message:   "message",
					},
				},
			},
		}, {
			"FineTuneStreamEvents",
			func() (interface{}, error) {
				var events []gpt3.FineTuneEvent
				onEvent := func(event *gpt3.FineTuneEvent) {
					events = append(events, *event)
				}
				return nil, client.FineTuneStreamEvents(ctx, gpt3.FineTuneEventsRequest{
					FineTuneID: "fine-tune-id",
				}, onEvent)
			},
			nil,
		}, {
			"DeleteFineTuneModel",
			func() (interface{}, error) {
				return client.DeleteFineTuneModel(ctx, "model-id")
			},
			&gpt3.DeleteFineTuneModelResponse{
				ID:      "model-id",
				Object:  "object",
				Deleted: true,
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

// TODO: add streaming response tests
// TODO: add file content tests
