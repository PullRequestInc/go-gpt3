package gpt3

import "fmt"

// APIError represents an error that occurred on an API
type APIError struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	Type       string `json:"type"`
}

func (e APIError) Error() string {
	return fmt.Sprintf("[%d:%s] %s", e.StatusCode, e.Type, e.Message)
}

// APIErrorResponse is the full error response that has been returned by an API.
type APIErrorResponse struct {
	Error APIError `json:"error"`
}
