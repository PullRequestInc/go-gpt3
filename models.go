package gpt3

import (
	"context"
	"fmt"
)

// Model Types
const (
	TextAda001Model     = "text-ada-001"
	TextBabbage001Model = "text-babbage-001"
	TextCurie001Model   = "text-curie-001"
	TextDavinci001Model = "text-davinci-001"
	TextDavinci002Model = "text-davinci-002"
	TextDavinci003Model = "text-davinci-003"
	AdaModel            = "ada"
	BabbageModel        = "babbage"
	CurieModel          = "curie"
	DavinciModel        = "davinci"
	DefaultModel        = DavinciModel
)

// ModelObject
type ModelObject struct {
	ID          string   `json:"id"`
	Object      string   `json:"object"`
	OwnedBy     string   `json:"owned_by"`
	Permissions []string `json:"permissions"`
}

// ModelsResponse is returned from the Models API
type ModelsResponse struct {
	Data   []ModelObject `json:"data"`
	Object string        `json:"object"`
}

// Models lists the currently available models, and provides basic information about each
// option such as the owner and permissioning.
func (c *client) Models(ctx context.Context) (*ModelsResponse, error) {
	req, err := c.newRequest(ctx, "GET", "/models", nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.performRequest(req)
	if err != nil {
		return nil, err
	}

	output := new(ModelsResponse)
	if err := getResponseObject(resp, output); err != nil {
		return nil, err
	}
	return output, nil
}

// Model retrieves a single model, providing basic information about the model such
// as the owner and permissioning.
func (c *client) Model(ctx context.Context, model string) (*ModelObject, error) {
	req, err := c.newRequest(ctx, "GET", fmt.Sprintf("/models/%s", model), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.performRequest(req)
	if err != nil {
		return nil, err
	}

	output := new(ModelObject)
	if err := getResponseObject(resp, output); err != nil {
		return nil, err
	}
	return output, nil
}
