package dspy

import (
	"encoding/json"
	"fmt"
)

// ConfigSerializer interface for serializing configurations
type ConfigSerializer interface {
	Serialize(config interface{}) (string, error)
}

// JSONConfigSerializer implements ConfigSerializer using JSON
type JSONConfigSerializer struct{}

// NewJSONConfigSerializer creates a new JSON serializer
func NewJSONConfigSerializer() ConfigSerializer {
	return &JSONConfigSerializer{}
}

// Serialize converts config to JSON string
func (s *JSONConfigSerializer) Serialize(config interface{}) (string, error) {
	if config == nil {
		return "{}", nil
	}
	
	data, err := json.Marshal(config)
	if err != nil {
		return "", fmt.Errorf("failed to serialize config: %w", err)
	}
	
	return string(data), nil
}
