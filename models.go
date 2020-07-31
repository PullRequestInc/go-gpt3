package gpt3

// EngineObject contained in an engine reponse
type EngineObject struct {
	ID     string `json:"id"`
	Object string `json:"object"`
	Owner  string `json:"owner"`
	Ready  bool   `json:"ready"`
}

// EnginesResponse is returned from the Engines API
type EnginesResponse struct {
	Data   []EngineObject `json:"data"`
	Object string         `json:"object"`
}

// CompletionRequest is a request for the completions API
type CompletionRequest struct {
	Prompt      string  `json:"prompt"`
	MaxTokens   int     `json:"max_tokens"`
	Temperature float32 `json:"temperature,omitempty"`
	TopP        float32 `json:"top_p,omitempty"`
	// N           float32 `json:"n"`
	Stream bool `json:"stream,omitempty"`
	// logprobs
	Stop string `json:"stop,omitempty"`
}

// CompletionResponseChoice is one of the choices returned in the response to the Completions API
type CompletionResponseChoice struct {
	Text  string `json:"text"`
	Index int    `json:"index"`
	// logprobs
	FinishReason string `json:"finish_reason"`
}

// CompletionResponse is the full response from a request to the completions API
type CompletionResponse struct {
	ID      string                     `json:"id"`
	Object  string                     `json:"object"`
	Created int                        `json:"created"`
	Model   string                     `json:"model"`
	Choices []CompletionResponseChoice `json:"choices"`
}

// SearchRequest is a request for the document search API
type SearchRequest struct {
	EngineID  string   `json:"engine_id"`
	Documents []string `json:"documents"`
	Query     string   `json:"query"`
}

// SearchData is a single search result from the document search API
type SearchData struct {
	Document int     `json:"document"`
	Object   string  `json:"search_resulet"`
	Score    float64 `json:"score"`
}

// SearchResponse is the full response from a request to the document search API
type SearchResponse struct {
	Data   []SearchData `json:"data"`
	Object string       `json:"object"`
}
