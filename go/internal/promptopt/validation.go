package promptopt

import (
	"errors"
	"fmt"
)

// Validate checks if PromptConfig is valid
func (c *PromptConfig) Validate() error {
	if c.AgentID == "" {
		return errors.New("AgentID is required")
	}
	if c.BasePrompt == "" {
		return errors.New("BasePrompt is required")
	}
	if len(c.TestCases) == 0 {
		return fmt.Errorf("at least one test case is required")
	}
	return nil
}
