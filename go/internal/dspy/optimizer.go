package dspy

import (
	"encoding/json"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	v1alpha2 "github.com/kagent-dev/kagent/go/api/v1alpha2"
	"github.com/kagent-dev/kagent/go/internal/database"
	"gorm.io/gorm"
)

// Optimizer handles DSPy optimization jobs
type Optimizer struct {
	db       *gorm.DB
	compiler *Compiler
	serializer ConfigSerializer
}

// NewOptimizer creates a new optimizer
func NewOptimizer(db *gorm.DB, compiler *Compiler) *Optimizer {
	serializer := NewJSONConfigSerializer()
	return &Optimizer{
		db:         db,
		compiler:   compiler,
		serializer: serializer,
	}
}

// CreateOptimizationJob creates a new optimization job
func (o *Optimizer) CreateOptimizationJob(ctx context.Context, agentID string, config *v1alpha2.OptimizationConfig) (string, error) {
	jobID := uuid.New().String()
	now := time.Now()

	job := &database.OptimizationJob{
		ID:        jobID,
		AgentID:   agentID,
		Status:    string(OptimizationJobStatusPending),
		Optimizer: config.Optimizer,
		StartedAt: &now,
	}

	// Serialize config to JSON
	configJSON, err := o.serializer.Serialize(config)
	if err != nil {
		return "", fmt.Errorf("failed to serialize config: %w", err)
	}
	job.Config = configJSON

	if err := o.db.WithContext(ctx).Create(job).Error; err != nil {
		return "", fmt.Errorf("failed to create optimization job: %w", err)
	}

	return jobID, nil
}

// GetOptimizationJob retrieves an optimization job by ID
func (o *Optimizer) GetOptimizationJob(ctx context.Context, jobID string) (*database.OptimizationJob, error) {
	var job database.OptimizationJob
	if err := o.db.WithContext(ctx).Where("id = ?", jobID).First(&job).Error; err != nil {
		return nil, fmt.Errorf("failed to get optimization job: %w", err)
	}
	return &job, nil
}

// UpdateOptimizationJobStatus updates the status of an optimization job
func (o *Optimizer) UpdateOptimizationJobStatus(ctx context.Context, jobID string, status OptimizationJobStatus, errorMsg *string) error {
	updates := map[string]interface{}{
		"status": string(status),
	}

	if status == OptimizationJobStatusCompleted || status == OptimizationJobStatusFailed {
		now := time.Now()
		updates["completed_at"] = &now
	}

	if errorMsg != nil {
		updates["error_msg"] = *errorMsg
	}

	if err := o.db.WithContext(ctx).Model(&database.OptimizationJob{}).Where("id = ?", jobID).Updates(updates).Error; err != nil {
		return fmt.Errorf("failed to update optimization job status: %w", err)
	}

	return nil
}

// ListOptimizationJobs lists all optimization jobs for an agent
func (o *Optimizer) ListOptimizationJobs(ctx context.Context, agentID string) ([]database.OptimizationJob, error) {
	var jobs []database.OptimizationJob
	if err := o.db.WithContext(ctx).Where("agent_id = ?", agentID).Order("created_at DESC").Find(&jobs).Error; err != nil {
		return nil, fmt.Errorf("failed to list optimization jobs: %w", err)
	}
	return jobs, nil
}

// Helper function to serialize config
func serializeConfig(config interface{}) (string, error) {
	// Serialize the config to JSON
	data, err := json.Marshal(config)
	if err != nil {
		return "", fmt.Errorf("failed to serialize config: %w", err)
	}
	return string(data), nil
}
