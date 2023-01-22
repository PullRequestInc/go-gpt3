package gpt3_test

import (
	"net/http"
	"testing"

	"github.com/PullRequestInc/go-gpt3"
	fakes "github.com/PullRequestInc/go-gpt3/go-gpt3fakes"
	"github.com/stretchr/testify/assert"
)

func TestInitNewClient(t *testing.T) {
	client, err := gpt3.NewClient("test-key")
	assert.Nil(t, err)
	assert.NotNil(t, client)
}

func fakeHttpClient() (*fakes.FakeRoundTripper, *http.Client) {
	rt := &fakes.FakeRoundTripper{}
	return rt, &http.Client{
		Transport: rt,
	}
}
