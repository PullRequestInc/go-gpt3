package gpt3_test

import (
	"net/http"

	fakes "github.com/PullRequestInc/go-gpt3/go-gpt3fakes"
)

func fakeHttpClient() (*fakes.FakeRoundTripper, *http.Client) {
	rt := &fakes.FakeRoundTripper{}
	return rt, &http.Client{
		Transport: rt,
	}
}
