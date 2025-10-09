package promptopt

import (
	"github.com/kagent-dev/kagent/go/internal/database"
	"github.com/kagent-dev/kagent/go/internal/dspy"
)

type OptimizerType string

const (
	OptimizerTypeGenetic OptimizerType = "genetic"
	OptimizerTypeDSPy    OptimizerType = "dspy"
)

func NewOptimizer(optimizerType OptimizerType, db database.Client) PromptOptimizer {
	switch optimizerType {
	case OptimizerTypeDSPy:
		compiler := dspy.NewCompiler("http://dspy-service:8000")
		optimizer := dspy.NewOptimizer(db.DB(), compiler) // Use db.DB() to get *gorm.DB
		return NewDSPyOptimizer(compiler, optimizer)
	default:
		return NewGeneticOptimizer()
	}
}
