package dspy

import (
	"context"

	"github.com/kagent-dev/kagent/go/internal/database"
)

// OptimizationRepository interface for database operations
type OptimizationRepository interface {
	Create(ctx context.Context, job *database.OptimizationJob) error
	Get(ctx context.Context, id string) (*database.OptimizationJob, error)
	List(ctx context.Context, agentID string) ([]*database.OptimizationJob, error)
	Update(ctx context.Context, job *database.OptimizationJob) error
}

// GormOptimizationRepository implements OptimizationRepository using GORM
type GormOptimizationRepository struct {
	db database.Client
}

// NewGormOptimizationRepository creates a new GORM-based repository
func NewGormOptimizationRepository(db database.Client) OptimizationRepository {
	return &GormOptimizationRepository{db: db}
}

// Create inserts a new optimization job
func (r *GormOptimizationRepository) Create(ctx context.Context, job *database.OptimizationJob) error {
	return r.db.DB().WithContext(ctx).Create(job).Error
}

// Get retrieves an optimization job by ID
func (r *GormOptimizationRepository) Get(ctx context.Context, id string) (*database.OptimizationJob, error) {
	var job database.OptimizationJob
	err := r.db.DB().WithContext(ctx).Where("id = ?", id).First(&job).Error
	if err != nil {
		return nil, err
	}
	return &job, nil
}

// List retrieves all optimization jobs for an agent
func (r *GormOptimizationRepository) List(ctx context.Context, agentID string) ([]*database.OptimizationJob, error) {
	var jobs []*database.OptimizationJob
	err := r.db.DB().WithContext(ctx).
		Where("agent_id = ?", agentID).
		Order("created_at DESC").
		Find(&jobs).Error
	if err != nil {
		return nil, err
	}
	return jobs, nil
}

// Update updates an existing optimization job
func (r *GormOptimizationRepository) Update(ctx context.Context, job *database.OptimizationJob) error {
	return r.db.DB().WithContext(ctx).Save(job).Error
}
