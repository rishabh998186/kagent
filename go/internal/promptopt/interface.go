package promptopt

import "context"

type PromptOptimizer interface {
	Compile(ctx context.Context, config *PromptConfig) (*CompiledPrompt, error)
	Optimize(ctx context.Context, job *OptimizationJob) (*OptimizedPrompt, error)
}

type PromptConfig struct {
	AgentID    string
	BasePrompt string
	TestCases  []TestCase
}

type TestCase struct {
	Input          map[string]interface{}
	ExpectedOutput map[string]interface{}
	Weight         float64
}

type CompiledPrompt struct {
	PromptID string
	Content  string
	Metadata map[string]string
}

type OptimizationJob struct {
	JobID   string
	AgentID string
	Config  *PromptConfig
	Status  string
}

type OptimizedPrompt struct {
	PromptID           string
	Content            string
	PerformanceMetrics map[string]float64
	Iterations         int
}
