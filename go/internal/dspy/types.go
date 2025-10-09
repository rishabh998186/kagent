package dspy

import (
	"time"
)

// CompileRequest represents a request to compile a DSPy signature
type CompileRequest struct {
	Inputs       []SignatureField `json:"inputs"`
	Outputs      []SignatureField `json:"outputs"`
	Instructions string           `json:"instructions,omitempty"`
	Module       string           `json:"module"`
}

// SignatureField represents a field in a DSPy signature
type SignatureField struct {
	Name        string  `json:"name"`
	Type        string  `json:"type"`
	Description *string `json:"description,omitempty"`
	Prefix      *string `json:"prefix,omitempty"`
}

// CompileResponse represents the response from DSPy compilation
type CompileResponse struct {
	CompiledPrompt string                 `json:"compiled_prompt"`
	SignatureDict  map[string]interface{} `json:"signature_dict"`
	ModuleType     string                 `json:"module_type"`
}

// OptimizationJobStatus represents the status of an optimization job
type OptimizationJobStatus string

const (
	OptimizationJobStatusPending   OptimizationJobStatus = "pending"
	OptimizationJobStatusRunning   OptimizationJobStatus = "running"
	OptimizationJobStatusCompleted OptimizationJobStatus = "completed"
	OptimizationJobStatusFailed    OptimizationJobStatus = "failed"
)

// OptimizeRequest represents a request to optimize a DSPy module
type OptimizeRequest struct {
	Signature    CompileRequest           `json:"signature"`
	Optimizer    string                   `json:"optimizer"`
	TrainingData []map[string]interface{} `json:"training_data"`
	MetricName   string                   `json:"metric_name"`
	Config       map[string]interface{}   `json:"config,omitempty"`
}

// OptimizeResponse represents the response from DSPy optimization
type OptimizeResponse struct {
	JobID           string                 `json:"job_id"`
	Status          OptimizationJobStatus  `json:"status"`
	OptimizedPrompt string                 `json:"optimized_prompt,omitempty"`
	Metrics         map[string]interface{} `json:"metrics,omitempty"`
	Error           string                 `json:"error,omitempty"`
}

// OptimizationProgress represents the progress of an optimization job
type OptimizationProgress struct {
	JobID            string                 `json:"job_id"`
	Status           OptimizationJobStatus  `json:"status"`
	Progress         float64                `json:"progress"`
	CurrentIteration int                    `json:"current_iteration"`
	TotalIterations  int                    `json:"total_iterations"`
	Metrics          map[string]interface{} `json:"metrics,omitempty"`
	StartedAt        *time.Time             `json:"started_at,omitempty"`
	CompletedAt      *time.Time             `json:"completed_at,omitempty"`
	Error            string                 `json:"error,omitempty"`
}
