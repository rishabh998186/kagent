package promptopt

import "time"

// OptimizationMetrics tracks optimization performance
type OptimizationMetrics struct {
	StartTime    time.Time
	Duration     time.Duration
	Iterations   int
	BestScore    float64
	Improvements []float64
}

// NewMetrics creates a new metrics tracker
func NewMetrics() *OptimizationMetrics {
	return &OptimizationMetrics{
		StartTime:    time.Now(),
		Improvements: []float64{},
	}
}

// RecordImprovement adds a score improvement
func (m *OptimizationMetrics) RecordImprovement(score float64) {
	m.Improvements = append(m.Improvements, score)
	if score > m.BestScore {
		m.BestScore = score
	}
}

// Finalize completes the metrics
func (m *OptimizationMetrics) Finalize() {
	m.Duration = time.Since(m.StartTime)
}
