# DSPy Integration for Kagent

This document describes the DSPy prompt optimization integration in Kagent.

## Overview

DSPy (Declarative Self-improving Python) is integrated into Kagent to provide automatic prompt compilation and optimization capabilities for AI agents.

## Architecture

### Components

1. **Go Backend (`go/internal/dspy/`)**
   - `compiler.go`: HTTP client for DSPy compilation
   - `optimizer.go`: Manages optimization jobs and stores results in database

2. **Python Service (`python-dspy-server/`)**
   - FastAPI service exposing DSPy functionality
   - Endpoints: `/compile`, `/optimize`, `/health`

3. **HTTP API Routes**
   - `POST /api/agents/{namespace}/{name}/dspy/compile` - Compile a prompt
   - `POST /api/agents/{namespace}/{name}/dspy/optimize` - Start optimization
   - `GET /api/agents/{namespace}/{name}/dspy/optimize/{jobId}` - Get optimization job
   - `GET /api/agents/{namespace}/{name}/dspy/optimize` - List optimization jobs

## Setup

### Prerequisites

- Kubernetes cluster
- Helm 3+
- OpenAI API key (or other LLM provider)

### Deployment

1. **Create secret for OpenAI API key:**

kubectl create secret generic dspy-secrets
--from-literal=openai-api-key=YOUR_API_KEY

2. **Deploy via Helm:**

helm install kagent ./helm/kagent
--set dspy.image=your-registry/dspy-service
--set dspy.tag=latest

3. **Verify deployment:**

kubectl get pods -l app=dspy-service
kubectl logs -l app=dspy-service

## Usage Examples

### 1. Compile a Prompt


curl -X POST http://localhost:8001/api/agents/default/my-agent/dspy/compile
-H "Content-Type: application/json"
-d '{
"module": "ChainOfThought",
"inputs": [
{"name": "question", "type": "string", "description": "User question"}
],
"outputs": [
{"name": "answer", "type": "string", "description": "Agent answer"}
],
"instructions": "Answer the question thoughtfully"
}'



### 2. Optimize a Prompt


curl -X POST http://localhost:8001/api/agents/default/my-agent/dspy/optimize
-H "Content-Type: application/json"
-d '{
"module": "ChainOfThought",
"inputs": [{"name": "question", "type": "string"}],
"outputs": [{"name": "answer", "type": "string"}],
"optimizer": "BootstrapFewShot",
"training_data": [
{
"inputs": {"question": "What is 2+2?"},
"outputs": {"answer": "4"}
}
]
}'

## Database Schema

Optimization results are stored in the `dspy_optimization_jobs` table:


CREATE TABLE dspy_optimization_jobs (
id VARCHAR PRIMARY KEY,
agent_id VARCHAR NOT NULL,
module_type VARCHAR NOT NULL,
optimizer VARCHAR NOT NULL,
status VARCHAR NOT NULL,
optimized_prompt TEXT,
metrics JSONB,
created_at TIMESTAMP,
updated_at TIMESTAMP
);

## Development

### Building the Python Service


cd python-dspy-server
docker build -t dspy-service:latest .
docker run -p 8000:8000 -e OPENAI_API_KEY=your-key dspy-service:latest

### Testing


Test health endpoint
curl http://localhost:8000/health

Test compilation (mock mode without API key)
curl -X POST http://localhost:8000/compile
-H "Content-Type: application/json"
-d '{"module": "Predict", "inputs": [], "outputs": []

## Configuration

### Environment Variables (Python Service)

- `DSPY_LM_PROVIDER`: LLM provider (default: "openai")
- `OPENAI_API_KEY`: OpenAI API key (required for real optimization)

### Helm Values


dspy:
image: your-registry/dspy-service
tag: latest
llmProvider: openai

## Troubleshooting

### Service not starting
- Check if OpenAI API key secret exists: `kubectl get secret dspy-secrets`
- View logs: `kubectl logs -l app=dspy-service`

### Optimization fails
- Service runs in MOCK mode without API key
- Check Python service logs for detailed errors
- Verify training data format matches DSPy requirements

## References

- [DSPy Documentation](https://dspy-docs.vercel.app/)
- [Kagent Documentation](../README.md)

## Design Decisions

### Type Definitions

The codebase maintains separate `SignatureField` types in two locations:

1. **Internal Types** (`go/internal/dspy/types.go`): 
   - Used for HTTP communication with DSPy service
   - Pointer fields for optional values

2. **API Types** (`go/api/v1alpha2/agent_types.go`):
   - Used for Kubernetes CRD definitions
   - String fields for simpler serialization

**Rationale**: This separation follows clean architecture principles, keeping API contracts separate from internal implementation details.

### Database Models

Two separate models exist for optimization data:

1. **OptimizationJob**: Tracks the optimization process and status
2. **PromptArtifact**: Stores the final optimized prompts

**Rationale**: Separating job tracking from artifact storage allows for better querying and lifecycle management.


### Why Duplicate Type Definitions?

The `SignatureField` type appears in two locations with different implementations:

**Internal Type** (`go/internal/dspy/types.go`):

type SignatureField struct {
Name string json:"name"
Type string json:"type"
Description *string json:"description,omitempty" // Pointer!
Prefix *string json:"prefix,omitempty" // Pointer!
}

**API Type** (`go/api/v1alpha2/agent_types.go`):

type SignatureField struct {
Name string json:"name"
Type string json:"type"
Description string json:"description,omitempty" // Plain string
Prefix string json:"prefix,omitempty" // Plain string
}

**Key Differences**:
1. **Pointer vs Value**: Internal uses pointers for true optionality in JSON
2. **Kubernetes Markers**: API type has CRD validation markers
3. **Purpose**: API is user-facing contract, internal is implementation detail

**Pattern**: This follows the **Anti-Corruption Layer** pattern from Domain-Driven Design, preventing external API changes from affecting internal implementation.

**Verdict**: ✅ Keep both - this is intentional good architecture!

