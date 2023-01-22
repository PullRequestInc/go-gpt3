package gpt3

import (
	"net/http"
	"time"
)

// ClientOption are options that can be passed when creating a new client
type ClientOption func(*client) error

// WithOrg is a client option that allows you to set the organization ID
func WithOrg(id string) ClientOption {
	return func(c *client) error {
		c.OrgID = id
		return nil
	}
}

// WithDefaultModel is a client option that allows you to override the default model of the client
func WithDefaultModel(model string) ClientOption {
	return func(c *client) error {
		c.DefaultModel = model
		return nil
	}
}

// WithUserAgent is a client option that allows you to override the default user agent of the client
func WithUserAgent(userAgent string) ClientOption {
	return func(c *client) error {
		c.UserAgent = userAgent
		return nil
	}
}

// WithBaseURL is a client option that allows you to override the default base url of the client.
// The default base url is "https://api.openai.com/v1"
func WithBaseURL(baseURL string) ClientOption {
	return func(c *client) error {
		c.BaseURL = baseURL
		return nil
	}
}

// WithHTTPClient allows you to override the internal http.Client used
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *client) error {
		c.HttpClient = httpClient
		return nil
	}
}

// WithTimeout is a client option that allows you to override the default timeout duration of requests
// for the client. The default is 30 seconds. If you are overriding the http client as well, just include
// the timeout there.
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *client) error {
		c.HttpClient.Timeout = timeout
		return nil
	}
}
