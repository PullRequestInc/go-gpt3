package gpt3

import "context"

// These are the legacy methods that are deprecated and will be removed in a future version.

func (c *client) CompletionWithEngine(ctx context.Context, engine string, request CompletionRequest) (*CompletionResponse, error) {
	// CompletionWithEngine is deprecated. Use Completion instead.
	request.Model = engine
	return c.Completion(ctx, request)
}

func (c *client) CompletionStreamWithEngine(ctx context.Context, engine string, request CompletionRequest, onData func(*CompletionResponse)) error {
	// CompletionStreamWithEngine is deprecated. Use CompletionStream instead.
	request.Model = engine
	return c.CompletionStream(ctx, request, onData)
}
