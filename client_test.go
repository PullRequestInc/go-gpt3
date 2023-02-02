package gpt3_test

import (
	"testing"

	"github.com/PullRequestInc/go-gpt3"
	"github.com/stretchr/testify/assert"
)

func TestInitNewClient(t *testing.T) {
	client, err := gpt3.NewClient("test-key")
	assert.Nil(t, err)
	assert.NotNil(t, client)
}
