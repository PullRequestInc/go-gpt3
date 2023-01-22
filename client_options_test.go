package gpt3_test

import (
	"net/http"
	"testing"

	"github.com/PullRequestInc/go-gpt3"
	"github.com/stretchr/testify/assert"
)

func TestClientWithOrg(t *testing.T) {
	client, err := gpt3.NewClient("test-key", gpt3.WithOrg("test-org"))
	assert.Nil(t, err)
	assert.NotNil(t, client)
}

func TestClientWithDefaultModel(t *testing.T) {
	client, err := gpt3.NewClient("test-key", gpt3.WithDefaultModel("test-model"))
	assert.Nil(t, err)
	assert.NotNil(t, client)
}

func TestClientWithUserAgent(t *testing.T) {
	client, err := gpt3.NewClient("test-key", gpt3.WithUserAgent("test-agent"))
	assert.Nil(t, err)
	assert.NotNil(t, client)
}

func TestClientWithBaseURL(t *testing.T) {
	client, err := gpt3.NewClient("test-key", gpt3.WithBaseURL("test-url"))
	assert.Nil(t, err)
	assert.NotNil(t, client)
}

func TestClientWithHTTPClient(t *testing.T) {
	httpClient := &http.Client{}
	client, err := gpt3.NewClient("test-key", gpt3.WithHTTPClient(httpClient))
	assert.Nil(t, err)
	assert.NotNil(t, client)
}

func TestClientWithTimeout(t *testing.T) {
	client, err := gpt3.NewClient("test-key", gpt3.WithTimeout(10))
	assert.Nil(t, err)
	assert.NotNil(t, client)
}
