package dspy

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	v1alpha2 "github.com/kagent-dev/kagent/go/api/v1alpha2"
)

// Compiler handles DSPy prompt compilation
type Compiler struct {
	serviceURL string
	httpClient *http.Client
}

// NewCompiler creates a new DSPy compiler client
func NewCompiler(serviceURL string) *Compiler {
	return &Compiler{
		serviceURL: serviceURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Compile compiles a DSPy signature into a prompt
func (c *Compiler) Compile(ctx context.Context, config *v1alpha2.DSPyConfig) (*CompileResponse, error) {
	req := c.buildCompileRequest(config)

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(
		ctx,
		"POST",
		c.serviceURL+"/compile",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("compilation failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var result CompileResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// buildCompileRequest converts v1alpha2.DSPyConfig to CompileRequest
func (c *Compiler) buildCompileRequest(config *v1alpha2.DSPyConfig) *CompileRequest {
	req := &CompileRequest{
		Module:       config.Module,
		Instructions: config.Signature.Instructions,
		Inputs:       make([]SignatureField, len(config.Signature.Inputs)),
		Outputs:      make([]SignatureField, len(config.Signature.Outputs)),
	}

	// Convert input fields
	for i, field := range config.Signature.Inputs {
		req.Inputs[i] = SignatureField{
			Name:        field.Name,
			Type:        field.Type,
			Description: func() *string {
				if field.Description != "" {
					return &field.Description
				}
				return nil
			}(),
			Prefix: func() *string {
				if field.Prefix != "" {
					return &field.Prefix
				}
				return nil
			}(),
		}
	}

	// Convert output fields
	for i, field := range config.Signature.Outputs {
		req.Outputs[i] = SignatureField{
			Name:        field.Name,
			Type:        field.Type,
			Description: func() *string {
				if field.Description != "" {
					return &field.Description
				}
				return nil
			}(),
			Prefix: func() *string {
				if field.Prefix != "" {
					return &field.Prefix
				}
				return nil
			}(),
		}
	}

	return req
}

// HealthCheck checks if the DSPy service is healthy
func (c *Compiler) HealthCheck(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, "GET", c.serviceURL+"/health", nil)
	if err != nil {
		return fmt.Errorf("failed to create health check request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("health check failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("health check returned status %d", resp.StatusCode)
	}

	return nil
}
