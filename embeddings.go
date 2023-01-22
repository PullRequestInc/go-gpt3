package gpt3

import "context"

const (
	TextSimilarityAda001      = "text-similarity-ada-001"
	TextSimilarityBabbage001  = "text-similarity-babbage-001"
	TextSimilarityCurie001    = "text-similarity-curie-001"
	TextSimilarityDavinci001  = "text-similarity-davinci-001"
	TextSearchAdaDoc001       = "text-search-ada-doc-001"
	TextSearchAdaQuery001     = "text-search-ada-query-001"
	TextSearchBabbageDoc001   = "text-search-babbage-doc-001"
	TextSearchBabbageQuery001 = "text-search-babbage-query-001"
	TextSearchCurieDoc001     = "text-search-curie-doc-001"
	TextSearchCurieQuery001   = "text-search-curie-query-001"
	TextSearchDavinciDoc001   = "text-search-davinci-doc-001"
	TextSearchDavinciQuery001 = "text-search-davinci-query-001"
	CodeSearchAdaCode001      = "code-search-ada-code-001"
	CodeSearchAdaText001      = "code-search-ada-text-001"
	CodeSearchBabbageCode001  = "code-search-babbage-code-001"
	CodeSearchBabbageText001  = "code-search-babbage-text-001"
	TextEmbeddingAda002       = "text-embedding-ada-002"
)

// EmbeddingsRequest is a request for the Embeddings API
type EmbeddingsRequest struct {
	// Input text to get embeddings for, encoded as a string or array of tokens. To get embeddings
	// for multiple inputs in a single request, pass an array of strings or array of token arrays.
	// Each input must not exceed 2048 tokens in length.
	Input []string `json:"input"`
	// ID of the model to use
	Model string `json:"model"`
	// The request user is an optional parameter meant to be used to trace abusive requests
	// back to the originating user. OpenAI states:
	// "The [user] IDs should be a string that uniquely identifies each user. We recommend hashing
	// their username or email address, in order to avoid sending us any identifying information.
	// If you offer a preview of your product to non-logged in users, you can send a session ID
	// instead."
	User string `json:"user,omitempty"`
}

// EmbeddingsResponse is the response from a create embeddings request.
type EmbeddingsResponse struct {
	Object string             `json:"object"`
	Data   []EmbeddingsResult `json:"data"`
	Usage  EmbeddingsUsage    `json:"usage"`
}

// The inner result of a create embeddings request, containing the embeddings for a single input.
type EmbeddingsResult struct {
	// The type of object returned (e.g., "list", "object")
	Object string `json:"object"`
	// The embedding data for the input
	Embedding []float64 `json:"embedding"`
	Index     int       `json:"index"`
}

// The usage stats for an embeddings response
type EmbeddingsUsage struct {
	// The number of tokens used by the prompt
	PromptTokens int `json:"prompt_tokens"`
	// The total tokens used
	TotalTokens int `json:"total_tokens"`
}

// Embeddings creates text embeddings for a supplied slice of inputs with a provided model.
//
// See: https://beta.openai.com/docs/api-reference/embeddings
func (c *client) Embeddings(ctx context.Context, request EmbeddingsRequest) (*EmbeddingsResponse, error) {
	req, err := c.newRequest(ctx, "POST", "/embeddings", request)
	if err != nil {
		return nil, err
	}
	resp, err := c.performRequest(req)
	if err != nil {
		return nil, err
	}

	output := EmbeddingsResponse{}
	if err := getResponseObject(resp, &output); err != nil {
		return nil, err
	}
	return &output, nil
}
