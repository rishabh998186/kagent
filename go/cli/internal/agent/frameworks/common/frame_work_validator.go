package common

import (
    "fmt"
    "slices"
)

// SupportedFrameworks lists all supported agent frameworks
var SupportedFrameworks = []string{"langgraph", "crewai"}

// ValidateFramework checks if the given framework is supported
func ValidateFramework(framework string) error {
    if slices.Contains(SupportedFrameworks, framework) {
        return nil
    }
    return fmt.Errorf("unsupported framework: %s (supported: %v)", framework, SupportedFrameworks)
}
