package gpt3

import "time"

// ClientOption are options that can be passed when creating a new client
type ClientOption func(*client) error

// WithDefaultEngine is a client option that allows you to override the default engine of the client
func WithDefaultEngine(engine string) ClientOption {
	return func(c *client) error {
		c.defaultEngine = engine
		return nil
	}
}

// WithUserAgent is a client option that allows you to override the default user agent of the client
func WithUserAgent(userAgent string) ClientOption {
	return func(c *client) error {
		c.userAgent = userAgent
		return nil
	}
}

// WithBaseURL is a client option that allows you to override the default base url of the client.
// The default base url is "https://api.openai.com/v1"
func WithBaseURL(baseURL string) ClientOption {
	return func(c *client) error {
		c.baseURL = baseURL
		return nil
	}
}

// WithTimeout is a client option that allows you to override the default timeout duration of requests
// for the client. The default is 30 seconds.
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *client) error {
		c.httpClient.Timeout = timeout
		return nil
	}
}
