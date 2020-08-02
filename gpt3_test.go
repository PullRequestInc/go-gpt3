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

////go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 net/http.RoundTripper

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

func TestEnginesRequestFails(t *testing.T) {
	rt, httpClient := fakeHttpClient()
	client := gpt3.NewClient("test-key", gpt3.WithHTTPClient(httpClient))

	rt.RoundTripReturns(nil, errors.New("request error"))

	ctx := context.Background()
	rsp, err := client.Engines(ctx)
	assert.Error(t, err, "request error")
	assert.Nil(t, rsp)
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
