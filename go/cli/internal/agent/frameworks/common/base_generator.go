package common

import (
        "io/fs"
	"fmt"
	"strings"
	"text/template"

	"github.com/kagent-dev/kagent/go/cli/internal/common/generator"
)

// AgentConfig holds the configuration for agent project generation
type AgentConfig struct {
	Name          string
	Directory     string
	Verbose       bool
	Instruction   string
	ModelProvider string
	ModelName     string
	Framework     string
	Language      string
	KagentVersion string
	McpServers    []McpServerType
	EnvVars       []string
}

// Implement ProjectConfig interface for AgentConfig
func (c AgentConfig) GetDirectory() string {
	return c.Directory
}

func (c AgentConfig) IsVerbose() bool {
	return c.Verbose
}

func (c AgentConfig) ShouldInitGit() bool {
	return true
}

func (c AgentConfig) ShouldSkipPath(path string) bool {
	// Skip mcp_server directory - these templates are processed separately
	return path == "mcp_server"
}

// BaseGenerator provides common functionality for all project generators.
// This now wraps the shared generator.BaseGenerator.
type BaseGenerator struct {
	*generator.BaseGenerator
}

// renderTemplate renders a template string with the provided data
func (g *BaseGenerator) renderTemplate(tmplContent string, data interface{}) (string, error) {
	tmpl, err := template.New("template").Funcs(template.FuncMap{
		"ToPascalCase": ToPascalCase,
		"ToUpper":      ToUpper,
                "upper":        ToUpper,
	}).Parse(tmplContent)

	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}
	var result strings.Builder
	if err := tmpl.Execute(&result, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return result.String(), nil
}

// GenerateProject generates a new project using the provided templates.
// This delegates to the shared generator implementation.
func (g *BaseGenerator) GenerateProject(config AgentConfig) error {
	return g.BaseGenerator.GenerateProject(config)
}

// RenderTemplate renders a template string with the provided data.
// This delegates to the shared generator implementation.
func (g *BaseGenerator) RenderTemplate(tmplContent string, data interface{}) (string, error) {
	return g.BaseGenerator.RenderTemplate(tmplContent, data)
}

// NewBaseGenerator creates a new BaseGenerator with the given template files
func NewBaseGenerator(templateFiles fs.FS) *BaseGenerator {
	return &BaseGenerator{
		BaseGenerator: generator.NewBaseGenerator(templateFiles, "templates"),
	}
}
