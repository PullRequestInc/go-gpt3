package gpt3

import "context"

// EditsRequest is a request for the edits API
type EditsRequest struct {
	// ID of the model to use. You can use the List models API to see all of your available models, or see our Model overview for descriptions of them.
	Model string `json:"model"`
	// The input text to use as a starting point for the edit.
	Input string `json:"input"`
	// The instruction that tells the model how to edit the prompt.
	Instruction string `json:"instruction"`
	// Sampling temperature to use
	Temperature *float32 `json:"temperature,omitempty"`
	// Alternative to temperature for nucleus sampling
	TopP *float32 `json:"top_p,omitempty"`
	// How many edits to generate for the input and instruction. Defaults to 1
	N *int `json:"n"`
}

// EditsResponse is the full response from a request to the edits API
type EditsResponse struct {
	Object  string                `json:"object"`
	Created int                   `json:"created"`
	Choices []EditsResponseChoice `json:"choices"`
	Usage   EditsResponseUsage    `json:"usage"`
}

// EditsResponseChoice is one of the choices returned in the response to the Edits API
type EditsResponseChoice struct {
	Text  string `json:"text"`
	Index int    `json:"index"`
}

// EditsResponseUsage is a structure used in the response from a request to the edits API
type EditsResponseUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

func (c *client) Edits(ctx context.Context, request EditsRequest) (*EditsResponse, error) {
	req, err := c.newRequest(ctx, "POST", "/edits", request)
	if err != nil {
		return nil, err
	}
	resp, err := c.performRequest(req)
	if err != nil {
		return nil, err
	}

	output := new(EditsResponse)
	if err := getResponseObject(resp, output); err != nil {
		return nil, err
	}
	return output, nil
}
