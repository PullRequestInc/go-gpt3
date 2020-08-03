package gpt3_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"testing"

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
			"Get https://api.openai.com/v1/engines: request error",
		},
		{
			"Engine",
			func() (interface{}, error) {
				return client.Engine(ctx, gpt3.DefaultEngine)
			},
			"Get https://api.openai.com/v1/engines/davinci: request error",
		},
		{
			"Completion",
			func() (interface{}, error) {
				return client.Completion(ctx, gpt3.CompletionRequest{})
			},
			"Post https://api.openai.com/v1/engines/davinci/completions: request error",
		}, {
			"CompletionStream",
			func() (interface{}, error) {
				var rsp *gpt3.CompletionResponse
				onData := func(data *gpt3.CompletionResponse) {
					rsp = data
				}
				return rsp, client.CompletionStream(ctx, gpt3.CompletionRequest{}, onData)
			},
			"Post https://api.openai.com/v1/engines/davinci/completions: request error",
		}, {
			"CompletionWithEngine",
			func() (interface{}, error) {
				return client.CompletionWithEngine(ctx, gpt3.AdaEngine, gpt3.CompletionRequest{})
			},
			"Post https://api.openai.com/v1/engines/ada/completions: request error",
		}, {
			"CompletionStreamWithEngine",
			func() (interface{}, error) {
				var rsp *gpt3.CompletionResponse
				onData := func(data *gpt3.CompletionResponse) {
					rsp = data
				}
				return rsp, client.CompletionStreamWithEngine(ctx, gpt3.AdaEngine, gpt3.CompletionRequest{}, onData)
			},
			"Post https://api.openai.com/v1/engines/ada/completions: request error",
		}, {
			"Search",
			func() (interface{}, error) {
				return client.Search(ctx, gpt3.SearchRequest{})
			},
			"Post https://api.openai.com/v1/engines/davinci/search: request error",
		}, {
			"SearchWithEngine",
			func() (interface{}, error) {
				return client.SearchWithEngine(ctx, gpt3.AdaEngine, gpt3.SearchRequest{})
			},
			"Post https://api.openai.com/v1/engines/ada/search: request error",
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

func TestResponseBadStatusCode(t *testing.T) {
	ctx := context.Background()
	rt, httpClient := fakeHttpClient()
	client := gpt3.NewClient("test-key", gpt3.WithHTTPClient(httpClient))

	testCases := map[string]func() (interface{}, error){
		"Engines": func() (interface{}, error) {
			return client.Engines(ctx)
		},
		"Engine": func() (interface{}, error) {
			return client.Engine(ctx, gpt3.DefaultEngine)
		},
		"Completion": func() (interface{}, error) {
			return client.Completion(ctx, gpt3.CompletionRequest{})
		},
		"CompletionStream": func() (interface{}, error) {
			var rsp *gpt3.CompletionResponse
			onData := func(data *gpt3.CompletionResponse) {
				rsp = data
			}
			return rsp, client.CompletionStream(ctx, gpt3.CompletionRequest{}, onData)
		},
		"CompletionWithEngine": func() (interface{}, error) {
			return client.CompletionWithEngine(ctx, gpt3.DefaultEngine, gpt3.CompletionRequest{})
		},
		"CompletionStreamWithEngine": func() (interface{}, error) {
			var rsp *gpt3.CompletionResponse
			onData := func(data *gpt3.CompletionResponse) {
				rsp = data
			}
			return rsp, client.CompletionStreamWithEngine(ctx, gpt3.DefaultEngine, gpt3.CompletionRequest{}, onData)
		},
		"Search": func() (interface{}, error) {
			return client.Search(ctx, gpt3.SearchRequest{})
		},
		"SearchWithEngine": func() (interface{}, error) {
			return client.SearchWithEngine(ctx, gpt3.DefaultEngine, gpt3.SearchRequest{})
		},
	}

	for name, apiCall := range testCases {
		t.Run(name, func(t *testing.T) {
			for _, code := range []int{301, 302, 400, 401, 404, 422, 500} {
				mockResponse := &http.Response{
					Status:     http.StatusText(200),
					StatusCode: code,
					Body:       ioutil.NopCloser(bytes.NewBufferString("unknown error")),
				}

				rt.RoundTripReturns(mockResponse, nil)
				rsp, err := apiCall()
				// assert.Nil(t, err)
				assert.Errorf(t, err, "request error")
				assert.Nil(t, rsp)
			}
		})
	}
}

func TestEnginesJsonDecodeFailure(t *testing.T) {
	rt, httpClient := fakeHttpClient()
	client := gpt3.NewClient("test-key", gpt3.WithHTTPClient(httpClient))

	mockResponse := &http.Response{
		Status:     http.StatusText(200),
		StatusCode: 200,
		Body:       ioutil.NopCloser(bytes.NewBufferString("invalid json")),
	}

	rt.RoundTripReturns(mockResponse, nil)

	ctx := context.Background()
	rsp, err := client.Engines(ctx)
	assert.Error(t, err, "invalid json response: invalid character 'i' looking for beginning of value")
	assert.Nil(t, rsp)
}

func TestEnginesSuccess(t *testing.T) {
	rt, httpClient := fakeHttpClient()
	client := gpt3.NewClient("test-key", gpt3.WithHTTPClient(httpClient))

	engines := &gpt3.EnginesResponse{
		Data: []gpt3.EngineObject{
			gpt3.EngineObject{
				ID:     "123",
				Object: "list",
				Owner:  "owner",
				Ready:  true,
			},
		},
	}

	data, err := json.Marshal(engines)
	assert.NoError(t, err)

	mockResponse := &http.Response{
		Status:     http.StatusText(200),
		StatusCode: 200,
		Body:       ioutil.NopCloser(bytes.NewBuffer(data)),
	}

	rt.RoundTripReturns(mockResponse, nil)

	ctx := context.Background()
	rsp, err := client.Engines(ctx)
	assert.NoError(t, err)
	assert.Equal(t, engines, rsp)
}
