package agent

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kagent-dev/kagent/go/cli/internal/agent/frameworks/adk/python"
	"github.com/kagent-dev/kagent/go/cli/internal/agent/frameworks/common"
	crewai "github.com/kagent-dev/kagent/go/cli/internal/agent/frameworks/crewai/python"
	langgraph "github.com/kagent-dev/kagent/go/cli/internal/agent/frameworks/langgraph/python"
)

// InitConfig holds configuration for agent initialization
type InitConfig struct {
	AgentName     string
	Framework     string
	Instruction   string
	ModelProvider string
	ModelName     string
	Description   string
	OutputDir     string
	Verbose       bool
	KagentVersion string
}

// InitAgent initializes a new agent project
func InitAgent(config InitConfig) error {
	// Validate agent name
	if config.AgentName == "" {
		return fmt.Errorf("agent name is required")
	}

	// Create generator registry
	registry := common.NewGeneratorRegistry()

	// Register all available generators
	if err := registry.Register(python.NewADKGenerator()); err != nil {
		return fmt.Errorf("failed to register ADK generator: %w", err)
	}
	if err := registry.Register(crewai.NewCrewAIGenerator()); err != nil {
		return fmt.Errorf("failed to register CrewAI generator: %w", err)
	}
	if err := registry.Register(langgraph.NewLangGraphGenerator()); err != nil {
		return fmt.Errorf("failed to register LangGraph generator: %w", err)
	}

	// Validate and normalize framework
	framework := strings.ToLower(config.Framework)
	if framework == "" {
		framework = "adk" // Default to ADK
	}

	// Get the appropriate generator
	generator, err := registry.Get(framework)
	if err != nil {
		return fmt.Errorf("failed to get generator: %w", err)
	}

	// Determine output directory
	outputDir := config.OutputDir
	if outputDir == "" {
		outputDir = config.AgentName
	}

	// Create output directory
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Get absolute path
	absOutputDir, err := filepath.Abs(outputDir)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	// Set default model provider and name if not specified
	if config.ModelProvider == "" {
		config.ModelProvider = "Gemini"
	}
	if config.ModelName == "" {
		config.ModelName = "Gemini-2.0-flash"
	}

	if config.Verbose {
		fmt.Printf("Initializing %s agent with %s framework\n", config.AgentName, framework)
		fmt.Printf("Output directory: %s\n", absOutputDir)
		fmt.Printf("Model: %s/%s\n", config.ModelProvider, config.ModelName)
	}

	// Generate the project
	err = generator.Generate(
		absOutputDir,
		config.AgentName,
		config.Instruction,
		config.ModelProvider,
		config.ModelName,
		config.Description,
		config.Verbose,
		config.KagentVersion,
	)

	if err != nil {
		return fmt.Errorf("failed to generate project: %w", err)
	}

	return nil
}

// ListFrameworks returns all supported frameworks
func ListFrameworks() []string {
	return common.SupportedFrameworks
}
