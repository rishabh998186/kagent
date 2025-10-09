package dspy

// DSPyField represents a unified field type for DSPy signatures
type DSPyField struct {
	Name        string
	Type        string
	Description string
	Prefix      string
}

// ToInternalField converts to internal API representation with pointers
func (f DSPyField) ToInternalField() SignatureField {
	field := SignatureField{
		Name: f.Name,
		Type: f.Type,
	}
	
	if f.Description != "" {
		field.Description = &f.Description
	}
	
	if f.Prefix != "" {
		field.Prefix = &f.Prefix
	}
	
	return field
}

// FromAPIField creates DSPyField from API v1alpha2 type
func FromAPIField(apiField interface{}) DSPyField {
	// Type assertion and conversion logic
	return DSPyField{}
}
