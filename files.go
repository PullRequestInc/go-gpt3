package gpt3

import (
	"context"
	"fmt"
	"io/ioutil"
)

// UploadFileRequest is a request for the Files API
type UploadFileRequest struct {
	// The file name of the JSON Lines file to upload
	File string `json:"file"`
	// The purpose of the file. Use "fine-tune" for a file that will be used to fine-tune a model.
	Purpose string `json:"purpose"`
}

// FileObject is a single file object
type FileObject struct {
	ID        string `json:"id"`
	Object    string `json:"object"`
	Bytes     int    `json:"bytes"`
	CreatedAt int    `json:"created_at"`
	Filename  string `json:"filename"`
	Purpose   string `json:"purpose"`
}

// FilesResponse is the response from a list files request.
type FilesResponse struct {
	Data   []FileObject `json:"data"`
	Object string       `json:"object"`
}

// DeleteFileResponse is the response from a delete file request.
type DeleteFileResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Deleted bool   `json:"deleted"`
}

// Files lists the files that belong to the user's organization.
//
// See: https://beta.openai.com/docs/api-reference/files/list
func (c *client) Files(ctx context.Context) (*FilesResponse, error) {
	req, err := c.newRequest(ctx, "GET", "/files", nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.performRequest(req)
	if err != nil {
		return nil, err
	}

	output := FilesResponse{}
	if err := getResponseObject(resp, &output); err != nil {
		return nil, err
	}
	return &output, nil
}

// UploadFile uploads a file that contains document(s) to be used across various endpoints.
//
// See: https://beta.openai.com/docs/api-reference/files/upload
func (c *client) UploadFile(ctx context.Context, request UploadFileRequest) (*FileObject, error) {
	req, err := c.newRequest(ctx, "POST", "/files", request)
	if err != nil {
		return nil, err
	}
	resp, err := c.performRequest(req)
	if err != nil {
		return nil, err
	}

	output := FileObject{}
	if err := getResponseObject(resp, &output); err != nil {
		return nil, err
	}
	return &output, nil
}

// DeleteFile deletes a file that contains document(s) to be used across various endpoints.
//
// See: https://beta.openai.com/docs/api-reference/files/delete
func (c *client) DeleteFile(ctx context.Context, fileID string) (*DeleteFileResponse, error) {
	req, err := c.newRequest(ctx, "DELETE", fmt.Sprintf("/files/%s", fileID), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.performRequest(req)
	if err != nil {
		return nil, err
	}

	output := DeleteFileResponse{}
	if err := getResponseObject(resp, &output); err != nil {
		return nil, err
	}
	return &output, nil
}

// File retrieves a file that contains document(s) to be used across various endpoints.
//
// See: https://beta.openai.com/docs/api-reference/files/retrieve
func (c *client) File(ctx context.Context, fileID string) (*FileObject, error) {
	req, err := c.newRequest(ctx, "GET", fmt.Sprintf("/files/%s", fileID), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.performRequest(req)
	if err != nil {
		return nil, err
	}

	output := FileObject{}
	if err := getResponseObject(resp, &output); err != nil {
		return nil, err
	}
	return &output, nil
}

// FileContent retrieves the content of a file that contains document(s) to be used across various endpoints.
//
// See: https://beta.openai.com/docs/api-reference/files/retrieve-content
func (c *client) FileContent(ctx context.Context, fileID string) ([]byte, error) {
	req, err := c.newRequest(ctx, "GET", fmt.Sprintf("/files/%s/content", fileID), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.performRequest(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}
