package promptopt

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

type GeneticOptimizer struct {
	maxIterations int
	populationSize int
}

func NewGeneticOptimizer() PromptOptimizer {
	return &GeneticOptimizer{
		maxIterations: 10,
		populationSize: 5,
	}
}

func (g *GeneticOptimizer) Compile(ctx context.Context, config *PromptConfig) (*CompiledPrompt, error) {
	// Validate config
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}
	
	return &CompiledPrompt{
		PromptID: fmt.Sprintf("genetic-%s-%d", config.AgentID, time.Now().Unix()),
		Content:  config.BasePrompt,
		Metadata: map[string]string{
			"type": "genetic",
			"compiled_at": time.Now().Format(time.RFC3339),
		},
	}, nil
}

func (g *GeneticOptimizer) Optimize(ctx context.Context, job *OptimizationJob) (*OptimizedPrompt, error) {
	// Validate config
	if err := job.Config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid job config: %w", err)
	}
	
	metrics := NewMetrics()
	
	// 1. Generate initial population of prompt variants
	variants := g.generateVariants(job.Config.BasePrompt, g.populationSize)
	
	bestPrompt := job.Config.BasePrompt
	bestScore := 0.0
	
	// 2. Evolutionary loop
	for i := 0; i < g.maxIterations; i++ {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		
		// Evaluate each variant
		scores := g.evaluateVariants(ctx, variants, job.Config.TestCases)
		
		// Find best performer
		for j, score := range scores {
			if score > bestScore {
				bestScore = score
				bestPrompt = variants[j]
				metrics.RecordImprovement(score)
			}
		}
		
		// Generate next generation
		variants = g.selectBest(variants, scores, g.populationSize/2)
		variants = append(variants, g.mutateVariants(variants)...)
		
		metrics.Iterations++
	}
	
	metrics.Finalize()
	
	return &OptimizedPrompt{
		PromptID: job.JobID,
		Content:  bestPrompt,
		PerformanceMetrics: map[string]float64{
			"accuracy": bestScore,
			"iterations": float64(metrics.Iterations),
			"duration_seconds": metrics.Duration.Seconds(),
		},
		Iterations: metrics.Iterations,
	}, nil
}

// generateVariants creates prompt variations
func (g *GeneticOptimizer) generateVariants(basePrompt string, count int) []string {
	variants := []string{basePrompt} // Keep original
	
	prefixes := []string{
		"Please carefully ",
		"Think step by step and ",
		"Let's approach this systematically: ",
		"Consider the following: ",
	}
	
	suffixes := []string{
		" Provide a detailed answer.",
		" Be specific and clear.",
		" Explain your reasoning.",
	}
	
	for i := 1; i < count; i++ {
		variant := basePrompt
		if i%2 == 0 && len(prefixes) > 0 {
			variant = prefixes[rand.Intn(len(prefixes))] + variant
		}
		if i%3 == 0 && len(suffixes) > 0 {
			variant = variant + suffixes[rand.Intn(len(suffixes))]
		}
		variants = append(variants, variant)
	}
	
	return variants
}

// evaluateVariants scores each variant (simplified)
func (g *GeneticOptimizer) evaluateVariants(ctx context.Context, variants []string, testCases []TestCase) []float64 {
	scores := make([]float64, len(variants))
	
	for i, variant := range variants {
		// Simple heuristic scoring based on prompt characteristics
		score := 0.5 // Base score
		
		// Longer prompts tend to be more detailed
		if len(variant) > 100 {
			score += 0.1
		}
		
		// Check for key instructional words
		instructionalWords := []string{"step", "detail", "explain", "specific", "clear"}
		for _, word := range instructionalWords {
			if strings.Contains(strings.ToLower(variant), word) {
				score += 0.05
			}
		}
		
		// Cap at 1.0
		if score > 1.0 {
			score = 1.0
		}
		
		scores[i] = score
	}
	
	return scores
}

// selectBest picks top performers
func (g *GeneticOptimizer) selectBest(variants []string, scores []float64, count int) []string {
	type pair struct {
		variant string
		score   float64
	}
	
	pairs := make([]pair, len(variants))
	for i := range variants {
		pairs[i] = pair{variants[i], scores[i]}
	}
	
	// Simple selection - take first 'count' (in production, sort by score)
	selected := make([]string, 0, count)
	for i := 0; i < count && i < len(pairs); i++ {
		selected = append(selected, pairs[i].variant)
	}
	
	return selected
}

// mutateVariants creates variations of existing prompts
func (g *GeneticOptimizer) mutateVariants(variants []string) []string {
	mutations := make([]string, 0, len(variants))
	
	for _, variant := range variants {
		// Simple mutation: add emphasis
		mutated := strings.ReplaceAll(variant, ".", "!")
		if mutated != variant {
			mutations = append(mutations, mutated)
		}
	}
	
	return mutations
}
