package gpt3_test

import (
	"net/http"
	"testing"

	"github.com/PullRequestInc/go-gpt3"
	"github.com/stretchr/testify/assert"
)

func TestClient(t *testing.T) {
	testCases := []struct {
		name    string
		options []gpt3.ClientOption
	}{
		{
			name: "test-key",
			options: []gpt3.ClientOption{
				gpt3.WithOrg("test-org"),
				gpt3.WithDefaultModel("test-model"),
				gpt3.WithUserAgent("test-agent"),
				gpt3.WithBaseURL("test-url"),
				gpt3.WithHTTPClient(&http.Client{}),
				gpt3.WithTimeout(10),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			client, err := gpt3.NewClient(tc.name, tc.options...)
			assert.Nil(t, err)
			assert.NotNil(t, client)
		})
	}
}
