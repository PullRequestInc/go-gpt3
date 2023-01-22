package gpt3

import "context"

// CompletionRequest is a request for the completions API
type CompletionRequest struct {
	// The model ID to use for completion
	Model string `json:"model"`
	// A list of string prompts to use.
	// TODO there are other prompt types here for using token integers that we could add support for.
	Prompt []string `json:"prompt"`
	// The suffix that comes after a completion of inserted text
	Suffix string `json:"suffix,omitempty"`
	// How many tokens to complete up to. Max of 512
	MaxTokens *int `json:"max_tokens,omitempty"`
	// Sampling temperature to use
	Temperature *float32 `json:"temperature,omitempty"`
	// Alternative to temperature for nucleus sampling
	TopP *float32 `json:"top_p,omitempty"`
	// How many choice to create for each prompt
	N *int `json:"n"`
	// Include the probabilities of most likely tokens
	LogProbs *int `json:"logprobs"`
	// Echo back the prompt in addition to the completion
	Echo bool `json:"echo"`
	// Up to 4 sequences where the API will stop generating tokens. Response will not contain the stop sequence.
	Stop []string `json:"stop,omitempty"`
	// PresencePenalty number between 0 and 1 that penalizes tokens that have already appeared in the text so far.
	PresencePenalty float32 `json:"presence_penalty"`
	// FrequencyPenalty number between 0 and 1 that penalizes tokens on existing frequency in the text so far.
	FrequencyPenalty float32 `json:"frequency_penalty"`
	// BestOf generates n completions server-side and returns the "best" one.
	BestOf *int `json:"best_of,omitempty"`
	// LogitBias is a list of token logit biases to apply before sampling.
	LogitBias []interface{} `json:"logit_bias,omitempty"`
	// User is the user ID to associate with this request
	User string `json:"user,omitempty"`

	// Whether to stream back results or not. Don't set this value in the request yourself
	// as it will be overridden depending on if you use CompletionStream or Completion methods.
	Stream bool `json:"stream,omitempty"`
}

// CompletionResponse is the full response from a request to the completions API
type CompletionResponse struct {
	ID      string                     `json:"id"`
	Object  string                     `json:"object"`
	Created int                        `json:"created"`
	Model   string                     `json:"model"`
	Choices []CompletionResponseChoice `json:"choices"`
	Usage   CompletionResponseUsage    `json:"usage"`
}

// CompletionResponseChoice is one of the choices returned in the response to the Completions API
type CompletionResponseChoice struct {
	Text         string        `json:"text"`
	Index        int           `json:"index"`
	LogProbs     LogprobResult `json:"logprobs"`
	FinishReason string        `json:"finish_reason"`
}

// LogprobResult represents logprob result of Choice
type LogprobResult struct {
	Tokens        []string             `json:"tokens"`
	TokenLogprobs []float32            `json:"token_logprobs"`
	TopLogprobs   []map[string]float32 `json:"top_logprobs"`
	TextOffset    []int                `json:"text_offset"`
}

// CompletionResponseUsage is the object that returns how many tokens the completion's request used
type CompletionResponseUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// Completion is a single completion.
func (c *client) Completion(ctx context.Context, request CompletionRequest) (*CompletionResponse, error) {
	request.Stream = false
	req, err := c.newRequest(ctx, "POST", "/completions", request)
	if err != nil {
		return nil, err
	}
	resp, err := c.performRequest(req)
	if err != nil {
		return nil, err
	}

	output := new(CompletionResponse)
	if err := getResponseObject(resp, output); err != nil {
		return nil, err
	}
	return output, nil
}
