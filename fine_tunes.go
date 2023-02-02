package gpt3

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
)

// CreateFineTuneRequest is a request for the FineTune API
type CreateFineTuneRequest struct {
	// The ID of an uploaded file that contains training data
	TrainingFile string `json:"training_file"`
	// The ID of an uploaded file that contains validation data
	ValidationFile string `json:"validation_file"`
	// The ID of the model to fine-tune
	Model string `json:"model"`
	// The number of epochs to train for
	NEpochs int `json:"n_epochs"`
	// The batch size to use for training
	BatchSize int `json:"batch_size"`
	// The learning rate to use for training
	LearningRate float32 `json:"learning_rate"`
	// The weight to use for loss on the prompt tokens
	PromptLossWeight float32 `json:"prompt_loss_weight"`
	// If set, we calculate classification-specific metrics using the validation set at the end of each epoch
	ComputeClassificationMetrics bool `json:"compute_classification_metrics"`
	// The number of classes in a classification task
	ClassificationNClasses int `json:"classification_n_classes"`
	// The positive class in binary classification
	ClassificationPositiveClass string `json:"classification_positive_class"`
	// If this is provided, we calculate F-beta scores at the specified beta values
	ClassificationBetas []float32 `json:"classification_betas"`
	// A string of up to 40 characters that will be added to your fine-tuned model name
	Suffix string `json:"suffix"`
}

// FineTuneRequest is a request for the FineTune API
type FineTuneRequest struct {
	// The ID of the fine-tune job
	FineTuneID string `json:"fine_tune_id"`
}

// FineTuneEventsRequest is a request for the FineTune API
type FineTuneEventsRequest struct {
	// The ID of the fine-tune job
	FineTuneID string `json:"fine_tune_id"`
	// Whether to stream events for the fine-tune job
	Stream bool `json:"stream"`
}

// DeleteFineTuneModelRequest is a request for the FineTune API
type DeleteFineTuneModelRequest struct {
	// The ID of the fine-tune model to delete
	Model string `json:"model"`
}

// FineTuneEvent is a single fine tune event
type FineTuneEvent struct {
	Object    string `json:"object"`
	CreatedAt int    `json:"created_at"`
	Level     string `json:"level"`
	Message   string `json:"message"`
}

// FineTuneHyperparams is the hyperparams for a fine tune request
type FineTuneHyperparams struct {
	BatchSize              int     `json:"batch_size"`
	LearningRateMultiplier float64 `json:"learning_rate_multiplier"`
	NEpochs                int     `json:"n_epochs"`
	PromptLossWeight       float64 `json:"prompt_loss_weight"`
}

// FineTuneObject is a single fine tune object
//
// See: https://beta.openai.com/docs/api-reference/fine-tunes/retrieve
type FineTuneObject struct {
	ID              string              `json:"id"`
	Object          string              `json:"object"`
	Model           string              `json:"model"`
	CreatedAt       int                 `json:"created_at"`
	Events          []FineTuneEvent     `json:"events"`
	FineTuneModel   string              `json:"fine_tune_model"`
	Hyperparams     FineTuneHyperparams `json:"hyperparams"`
	OrganizationID  string              `json:"organization_id"`
	ResultFiles     []FileObject        `json:"result_files"`
	Status          string              `json:"status"`
	ValidationFiles []FileObject        `json:"validation_files"`
	TrainingFiles   []FileObject        `json:"training_files"`
	UpdatedAt       int                 `json:"updated_at"`
}

// ListFineTunesResponse is the response from a list fine tunes request.
//
// See: https://beta.openai.com/docs/api-reference/fine-tunes/list
type ListFineTunesResponse struct {
	Data   []FineTuneObject `json:"data"`
	Object string           `json:"object"`
}

// FineTuneEventsResponse
//
// See: https://beta.openai.com/docs/api-reference/fine-tunes/events
type FineTuneEventsResponse struct {
	Data   []FineTuneEvent `json:"data"`
	Object string          `json:"object"`
}

// DeleteFineTuneModelResponse
//
// See: https://beta.openai.com/docs/api-reference/fine-tunes/delete-model
type DeleteFineTuneModelResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Deleted bool   `json:"deleted"`
}

// CreateFineTune creates a job that fine-tunes a model on a dataset.
//
// See: https://beta.openai.com/docs/api-reference/fine-tunes/create
func (c *client) CreateFineTune(ctx context.Context, request CreateFineTuneRequest) (*FineTuneObject, error) {
	req, err := c.newRequest(ctx, "POST", "/fine-tunes", request)
	if err != nil {
		return nil, err
	}
	resp, err := c.performRequest(req)
	if err != nil {
		return nil, err
	}

	output := FineTuneObject{}
	if err := getResponseObject(resp, &output); err != nil {
		return nil, err
	}
	return &output, nil
}

// ListFineTunes lists the fine-tuning jobs that belong to the user's organization.
//
// See: https://beta.openai.com/docs/api-reference/fine-tunes/list
func (c *client) ListFineTunes(ctx context.Context) (*ListFineTunesResponse, error) {
	req, err := c.newRequest(ctx, "GET", "/fine-tunes", nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.performRequest(req)
	if err != nil {
		return nil, err
	}

	output := ListFineTunesResponse{}
	if err := getResponseObject(resp, &output); err != nil {
		return nil, err
	}
	return &output, nil
}

// FineTune retrieves a fine-tuning job from the user's organization.
//
// See: https://beta.openai.com/docs/api-reference/fine-tunes/retrieve
func (c *client) FineTune(ctx context.Context, fineTuneID string) (*FineTuneObject, error) {
	req, err := c.newRequest(ctx, "GET", fmt.Sprintf("/fine-tunes/%s", fineTuneID), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.performRequest(req)
	if err != nil {
		return nil, err
	}

	output := FineTuneObject{}
	if err := getResponseObject(resp, &output); err != nil {
		return nil, err
	}
	return &output, nil
}

// CancelFineTune cancels a fine-tuning job from the user's organization.
//
// See: https://beta.openai.com/docs/api-reference/fine-tunes/cancel
func (c *client) CancelFineTune(ctx context.Context, fineTuneID string) (*FineTuneObject, error) {
	req, err := c.newRequest(ctx, "POST", fmt.Sprintf("/fine-tunes/%s/cancel", fineTuneID), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.performRequest(req)
	if err != nil {
		return nil, err
	}

	output := FineTuneObject{}
	if err := getResponseObject(resp, &output); err != nil {
		return nil, err
	}
	return &output, nil
}

// FineTuneEvents lists the events that belong to a fine-tuning job.
//
// See: https://beta.openai.com/docs/api-reference/fine-tunes/events
func (c *client) FineTuneEvents(ctx context.Context, request FineTuneEventsRequest) (*FineTuneEventsResponse, error) {
	request.Stream = false
	req, err := c.newRequest(ctx, "GET", fmt.Sprintf("/fine-tunes/%s/events", request.FineTuneID), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.performRequest(req)
	if err != nil {
		return nil, err
	}

	output := FineTuneEventsResponse{}
	if err := getResponseObject(resp, &output); err != nil {
		return nil, err
	}
	return &output, nil
}

// FineTuneStreamEvents streams the events that belong to a fine-tuning job.
//
// See: https://beta.openai.com/docs/api-reference/fine-tunes/events
func (c *client) FineTuneStreamEvents(ctx context.Context, request FineTuneEventsRequest, onData func(*FineTuneEvent)) error {
	request.Stream = true
	req, err := c.newRequest(ctx, "GET", fmt.Sprintf("/fine-tunes/%s/events", request.FineTuneID), nil)
	if err != nil {
		return err
	}
	resp, err := c.performRequest(req)
	if err != nil {
		return err
	}

	reader := bufio.NewReader(resp.Body)
	defer resp.Body.Close()

	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			return err
		}

		line = bytes.TrimSpace(line)
		if !bytes.HasPrefix(line, dataPrefix) {
			continue
		}
		line = bytes.TrimPrefix(line, dataPrefix)

		if bytes.HasPrefix(line, streamTerminationPrefix) {
			break
		}
		output := new(FineTuneEvent)
		if err := json.Unmarshal(line, output); err != nil {
			return fmt.Errorf("invalid json stream data: %v", err)
		}
		onData(output)
	}
	return nil
}

// DeleteFineTuneModel deletes a fine-tuned model from the user's organization.
//
// See: https://beta.openai.com/docs/api-reference/fine-tunes/delete-model
func (c *client) DeleteFineTuneModel(ctx context.Context, modelID string) (*DeleteFineTuneModelResponse, error) {
	req, err := c.newRequest(ctx, "DELETE", fmt.Sprintf("/models/%s", modelID), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.performRequest(req)
	if err != nil {
		return nil, err
	}

	output := DeleteFineTuneModelResponse{}
	if err := getResponseObject(resp, &output); err != nil {
		return nil, err
	}
	return &output, nil
}
