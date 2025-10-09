/*
Copyright 2025.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha2

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"trpc.group/trpc-go/trpc-a2a-go/server"
)

// AgentType represents the agent type
// +kubebuilder:validation:Enum=Declarative;BYO
type AgentType string

const (
	AgentType_Declarative AgentType = "Declarative"
	AgentType_BYO         AgentType = "BYO"
)

// AgentSpec defines the desired state of Agent.
// +kubebuilder:validation:XValidation:message="type must be specified",rule="has(self.type)"
// +kubebuilder:validation:XValidation:message="type must be either Declarative or BYO",rule="self.type == 'Declarative' || self.type == 'BYO'"
// +kubebuilder:validation:XValidation:message="declarative must be specified if type is Declarative, or byo must be specified if type is BYO",rule="(self.type == 'Declarative' && has(self.declarative)) || (self.type == 'BYO' && has(self.byo))"
type AgentSpec struct {
	// +kubebuilder:validation:Enum=Declarative;BYO
	// +kubebuilder:default=Declarative
	Type AgentType `json:"type"`

	// +optional
	BYO *BYOAgentSpec `json:"byo,omitempty"`
	// +optional
	Declarative *DeclarativeAgentSpec `json:"declarative,omitempty"`

	// +optional
	Description string `json:"description,omitempty"`
}

// +kubebuilder:validation:XValidation:rule="!has(self.systemMessage) || !has(self.systemMessageFrom)",message="systemMessage and systemMessageFrom are mutually exclusive"
// +kubebuilder:validation:XValidation:rule="!has(self.dspyConfig) || (!has(self.systemMessage) && !has(self.systemMessageFrom))",message="dspyConfig cannot be used with systemMessage or systemMessageFrom"
type DeclarativeAgentSpec struct {
	// SystemMessage is a string specifying the system message for the agent
	// +optional
	SystemMessage string `json:"systemMessage,omitempty"`
	// SystemMessageFrom is a reference to a ConfigMap or Secret containing the system message.
	// +optional
	SystemMessageFrom *ValueSource `json:"systemMessageFrom,omitempty"`
	// DSPyConfig enables DSPy-based prompt compilation and optimization
	// When specified, the agent will use DSPy framework to compile and optimize prompts
	// instead of using systemMessage or systemMessageFrom directly.
	// +optional
	DSPyConfig *DSPyConfig `json:"dspyConfig,omitempty"`

	// The name of the model config to use.
	// If not specified, the default value is "default-model-config".
	// Must be in the same namespace as the Agent.
	// +optional
	ModelConfig string `json:"modelConfig,omitempty"`
	// Whether to stream the response from the model.
	// If not specified, the default value is true.
	// +optional
	Stream *bool `json:"stream,omitempty"`
	// +kubebuilder:validation:MaxItems=20
	Tools []*Tool `json:"tools,omitempty"`
	// A2AConfig instantiates an A2A server for this agent,
	// served on the HTTP port of the kagent kubernetes
	// controller (default 8083).
	// The A2A server URL will be served at
	// <kagent-controller-ip>:8083/api/a2a/<agent-namespace>/<agent-name>
	// Read more about the A2A protocol here: https://github.com/google/A2A
	// +optional
	A2AConfig *A2AConfig `json:"a2aConfig,omitempty"`

	// +optional
	Deployment *DeclarativeDeploymentSpec `json:"deployment,omitempty"`
}

type DeclarativeDeploymentSpec struct {
	// +optional
	ImageRegistry string `json:"imageRegistry,omitempty"`

	SharedDeploymentSpec `json:",inline"`
}

type BYOAgentSpec struct {
	// Trust relationship to the agent.
	// +optional
	Deployment *ByoDeploymentSpec `json:"deployment,omitempty"`
}

type ByoDeploymentSpec struct {
	// +kubebuilder:validation:MinLength=1
	Image string `json:"image,omitempty"`
	// +optional
	Cmd *string `json:"cmd,omitempty"`
	// +optional
	Args []string `json:"args,omitempty"`

	SharedDeploymentSpec `json:",inline"`
}

type SharedDeploymentSpec struct {
	// If not specified, the default value is 1.
	// +optional
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:default=1
	Replicas *int32 `json:"replicas,omitempty"`
	// +optional
	ImagePullSecrets []corev1.LocalObjectReference `json:"imagePullSecrets,omitempty"`
	// +optional
	Volumes []corev1.Volume `json:"volumes,omitempty"`
	// +optional
	VolumeMounts []corev1.VolumeMount `json:"volumeMounts,omitempty"`
	// +optional
	Labels map[string]string `json:"labels,omitempty"`
	// +optional
	Annotations map[string]string `json:"annotations,omitempty"`
	// +optional
	Env []corev1.EnvVar `json:"env,omitempty"`
	// +optional
	ImagePullPolicy corev1.PullPolicy `json:"imagePullPolicy,omitempty"`
	// +optional
	Resources *corev1.ResourceRequirements `json:"resources,omitempty"`
}

// ToolProviderType represents the tool provider type
// +kubebuilder:validation:Enum=McpServer;Agent
type ToolProviderType string

const (
	ToolProviderType_McpServer ToolProviderType = "McpServer"
	ToolProviderType_Agent     ToolProviderType = "Agent"
)

// +kubebuilder:validation:XValidation:message="type.mcpServer must be nil if the type is not McpServer",rule="!(has(self.mcpServer) && self.type != 'McpServer')"
// +kubebuilder:validation:XValidation:message="type.mcpServer must be specified for McpServer filter.type",rule="!(!has(self.mcpServer) && self.type == 'McpServer')"
// +kubebuilder:validation:XValidation:message="type.agent must be nil if the type is not Agent",rule="!(has(self.agent) && self.type != 'Agent')"
// +kubebuilder:validation:XValidation:message="type.agent must be specified for Agent filter.type",rule="!(!has(self.agent) && self.type == 'Agent')"
type Tool struct {
	// +kubebuilder:validation:Enum=McpServer;Agent
	Type ToolProviderType `json:"type,omitempty"`
	// +optional
	McpServer *McpServerTool `json:"mcpServer,omitempty"`
	// +optional
	Agent *TypedLocalReference `json:"agent,omitempty"`

	// HeadersFrom specifies a list of configuration values to be added as
	// headers to requests sent to the Tool from this agent. The value of
	// each header is resolved from either a Secret or ConfigMap in the same
	// namespace as the Agent. Headers specified here will override any
	// headers of the same name/key specified on the tool.
	// +optional
	HeadersFrom []ValueRef `json:"headersFrom,omitempty"`
}

func (s *Tool) ResolveHeaders(ctx context.Context, client client.Client, namespace string) (map[string]string, error) {
	result := map[string]string{}

	for _, h := range s.HeadersFrom {
		k, v, err := h.Resolve(ctx, client, namespace)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve header: %v", err)
		}

		result[k] = v
	}

	return result, nil
}

type McpServerTool struct {
	// The reference to the ToolServer that provides the tool.
	// Can either be a reference to the name of a ToolServer in the same namespace as the referencing Agent, or a reference to the name of an ToolServer in a different namespace in the form <namespace>/<name>
	// +optional
	TypedLocalReference `json:",inline"`

	// The names of the tools to be provided by the ToolServer
	// For a list of all the tools provided by the server,
	// the client can query the status of the ToolServer object after it has been created
	ToolNames []string `json:"toolNames,omitempty"`
}

type TypedLocalReference struct {
	// +optional
	Kind string `json:"kind"`
	// +optional
	ApiGroup string `json:"apiGroup"`
	Name     string `json:"name"`
}

func (t *TypedLocalReference) GroupKind() schema.GroupKind {
	return schema.GroupKind{
		Group: t.ApiGroup,
		Kind:  t.Kind,
	}
}

type A2AConfig struct {
	// +kubebuilder:validation:MinItems=1
	Skills []AgentSkill `json:"skills,omitempty"`
}

type AgentSkill server.AgentSkill

const (
	AgentConditionTypeAccepted = "Accepted"
	AgentConditionTypeReady    = "Ready"
)

// AgentStatus defines the observed state of Agent.
type AgentStatus struct {
	ObservedGeneration int64              `json:"observedGeneration"`
	Conditions         []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Type",type="string",JSONPath=".spec.type",description="The type of the agent."
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status",description="Whether or not the agent is ready to serve requests."
// +kubebuilder:printcolumn:name="Accepted",type="string",JSONPath=".status.conditions[?(@.type=='Accepted')].status",description="Whether or not the agent has been accepted by the system."
// +kubebuilder:storageversion

// Agent is the Schema for the agents API.
type Agent struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AgentSpec   `json:"spec,omitempty"`
	Status AgentStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// AgentList contains a list of Agent.
type AgentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Agent `json:"items"`
}

// DSPyConfig defines the configuration for DSPy-based prompt compilation
type DSPyConfig struct {
	// Enabled indicates whether DSPy compilation is active for this agent
	// +kubebuilder:default=true
	Enabled bool `json:"enabled"`

	// Signature defines the input/output schema for the DSPy module
	Signature DSPySignature `json:"signature"`

	// Module specifies which DSPy module to use
	// +kubebuilder:validation:Enum=Predict;ChainOfThought;ReAct
	// +kubebuilder:default=ChainOfThought
	Module string `json:"module"`

	// CompiledPromptRef is a reference to a ConfigMap or Secret containing
	// the compiled prompt artifact. This is populated after successful compilation.
	// +optional
	CompiledPromptRef *ValueSource `json:"compiledPromptRef,omitempty"`

	// OptimizationConfig defines settings for prompt optimization
	// +optional
	OptimizationConfig *OptimizationConfig `json:"optimizationConfig,omitempty"`
}

// DSPySignature defines the input and output schema for a DSPy module
type DSPySignature struct {
	// Instructions provide context for what the module should do
	// +optional
	Instructions string `json:"instructions,omitempty"`

	// Inputs define the input fields for the signature
	// +kubebuilder:validation:MinItems=1
	Inputs []SignatureField `json:"inputs"`

	// Outputs define the output fields for the signature
	// +kubebuilder:validation:MinItems=1
	Outputs []SignatureField `json:"outputs"`
}

// SignatureField represents a single input or output field in a DSPy signature
type SignatureField struct {
	// Name of the field
	// +kubebuilder:validation:MinLength=1
	Name string `json:"name"`

	// Type of the field
	// +kubebuilder:validation:Enum=string;int;bool;list;float
	// +kubebuilder:default=string
	Type string `json:"type"`

	// Description provides context about what this field represents
	// +optional
	Description string `json:"description,omitempty"`

	// Prefix is an optional prefix for the field in the generated prompt
	// +optional
	Prefix string `json:"prefix,omitempty"`
}

// OptimizationConfig defines configuration for DSPy prompt optimization
type OptimizationConfig struct {
	// Enabled indicates whether optimization should be performed
	// +kubebuilder:default=false
	Enabled bool `json:"enabled"`

	// Optimizer specifies which DSPy optimizer algorithm to use
	// +kubebuilder:validation:Enum=MIPRO;MIPROv2;BootstrapFewShot;BootstrapFewShotWithRandomSearch;COPRO
	// +kubebuilder:default=MIPROv2
	Optimizer string `json:"optimizer"`

	// TrainingDataRef references a ConfigMap or Secret containing training examples
	// The data should be in JSON format with an array of input/output examples
	// +optional
	TrainingDataRef *ValueSource `json:"trainingDataRef,omitempty"`

	// MetricName specifies the evaluation metric to optimize for
	// +optional
	MetricName string `json:"metricName,omitempty"`

	// MaxBootstrappedDemos limits the number of demonstrations to bootstrap
	// Only applicable for bootstrap-based optimizers
	// +optional
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=20
	MaxBootstrappedDemos *int `json:"maxBootstrappedDemos,omitempty"`

	// MaxLabeledDemos limits the number of labeled demonstrations to use
	// +optional
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=50
	MaxLabeledDemos *int `json:"maxLabeledDemos,omitempty"`
}

// DSPyModuleType represents the supported DSPy module types
type DSPyModuleType string

const (
	DSPyModuleType_Predict        DSPyModuleType = "Predict"
	DSPyModuleType_ChainOfThought DSPyModuleType = "ChainOfThought"
	DSPyModuleType_ReAct          DSPyModuleType = "ReAct"
)

// DSPyOptimizerType represents the supported DSPy optimizer algorithms
type DSPyOptimizerType string

const (
	DSPyOptimizerType_MIPRO                            DSPyOptimizerType = "MIPRO"
	DSPyOptimizerType_MIPROv2                          DSPyOptimizerType = "MIPROv2"
	DSPyOptimizerType_BootstrapFewShot                 DSPyOptimizerType = "BootstrapFewShot"
	DSPyOptimizerType_BootstrapFewShotWithRandomSearch DSPyOptimizerType = "BootstrapFewShotWithRandomSearch"
	DSPyOptimizerType_COPRO                            DSPyOptimizerType = "COPRO"
)

func init() {
	SchemeBuilder.Register(&Agent{}, &AgentList{})
}
