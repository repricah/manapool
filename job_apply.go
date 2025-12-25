package manapool

import (
	"bytes"
	"context"
	"fmt"
	"mime/multipart"
)

const defaultApplicationFilename = "application.zip"

// SubmitJobApplication submits a job application.
func (c *Client) SubmitJobApplication(ctx context.Context, req JobApplicationRequest) (*JobApplicationResponse, error) {
	if req.FirstName == "" || req.LastName == "" || req.Email == "" || len(req.Application) == 0 {
		return nil, NewValidationError("application", "first name, last name, email, and application are required")
	}
	filename := req.ApplicationFilename
	if filename == "" {
		filename = defaultApplicationFilename
	}

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	if err := writer.WriteField("first_name", req.FirstName); err != nil {
		return nil, NewValidationError("first_name", "failed to encode first_name field: "+err.Error())
	}
	if err := writer.WriteField("last_name", req.LastName); err != nil {
		return nil, NewValidationError("last_name", "failed to encode last_name field: "+err.Error())
	}
	if err := writer.WriteField("email", req.Email); err != nil {
		return nil, NewValidationError("email", "failed to encode email field: "+err.Error())
	}
	if req.LinkedInURL != "" {
		if err := writer.WriteField("linkedin_url", req.LinkedInURL); err != nil {
			return nil, NewValidationError("linkedin_url", "failed to encode linkedin_url field: "+err.Error())
		}
	}
	if req.GitHubURL != "" {
		if err := writer.WriteField("github_url", req.GitHubURL); err != nil {
			return nil, NewValidationError("github_url", "failed to encode github_url field: "+err.Error())
		}
	}
	fileWriter, err := writer.CreateFormFile("application", filename)
	if err != nil {
		return nil, NewValidationError("application", "failed to create application form file: "+err.Error())
	}
	if _, err := fileWriter.Write(req.Application); err != nil {
		return nil, NewValidationError("application", "failed to encode application file data: "+err.Error())
	}
	if err := writer.Close(); err != nil {
		return nil, NewValidationError("application", "failed to finalize multipart body: "+err.Error())
	}

	resp, err := c.doRequestWithBody(ctx, "POST", "/job-apply", nil, &body, writer.FormDataContentType())
	if err != nil {
		return nil, fmt.Errorf("failed to submit job application: %w", err)
	}

	var response JobApplicationResponse
	if err := c.decodeResponse(resp, &response); err != nil {
		return nil, fmt.Errorf("failed to decode job application response: %w", err)
	}

	return &response, nil
}
