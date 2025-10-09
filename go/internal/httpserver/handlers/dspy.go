package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/kagent-dev/kagent/go/api/v1alpha2"
	"github.com/kagent-dev/kagent/go/internal/dspy"
	"github.com/kagent-dev/kagent/go/internal/httpserver/errors"
	"github.com/kagent-dev/kagent/go/internal/utils"
	"github.com/kagent-dev/kagent/go/pkg/auth" // ← Change this from internal/httpserver/auth
	"github.com/kagent-dev/kagent/go/pkg/client/api"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	ctrllog "sigs.k8s.io/controller-runtime/pkg/log"
)

// DSPyHandler handles DSPy-related requests
type DSPyHandler struct {
	*Base
	compiler  *dspy.Compiler
	optimizer *dspy.Optimizer
}

// NewDSPyHandler creates a new DSPyHandler
func NewDSPyHandler(base *Base, compiler *dspy.Compiler, optimizer *dspy.Optimizer) *DSPyHandler {
	return &DSPyHandler{
		Base:      base,
		compiler:  compiler,
		optimizer: optimizer,
	}
}

// CompileRequest represents a request to compile a DSPy prompt
type CompileRequest struct {
	DSPyConfig v1alpha2.DSPyConfig `json:"dspy_config"`
}

// CompileResponse represents the response from DSPy compilation
type CompileResponse struct {
	CompiledPrompt string                 `json:"compiled_prompt"`
	SignatureDict  map[string]interface{} `json:"signature_dict"`
	ModuleType     string                 `json:"module_type"`
}

// OptimizeRequest represents a request to start optimization
type OptimizeRequest struct {
	Optimizer       string                 `json:"optimizer"`
	TrainingDataRef string                 `json:"training_data_ref,omitempty"`
	MetricName      string                 `json:"metric_name,omitempty"`
	Config          map[string]interface{} `json:"config,omitempty"`
}

// HandleCompilePrompt handles POST /api/agents/{namespace}/{name}/dspy/compile
func (h *DSPyHandler) HandleCompilePrompt(w ErrorResponseWriter, r *http.Request) {
	log := ctrllog.FromContext(r.Context()).WithName("dspy-handler").WithValues("operation", "compile")

	agentName, err := GetPathParam(r, "name")
	if err != nil {
		w.RespondWithError(errors.NewBadRequestError("Failed to get name from path", err))
		return
	}

	agentNamespace, err := GetPathParam(r, "namespace")
	if err != nil {
		w.RespondWithError(errors.NewBadRequestError("Failed to get namespace from path", err))
		return
	}

	log = log.WithValues("agentName", agentName, "agentNamespace", agentNamespace)

	if err := Check(h.Authorizer, r, auth.Resource{
		Type: "Agent",
		Name: types.NamespacedName{Namespace: agentNamespace, Name: agentName}.String(),
	}); err != nil {
		w.RespondWithError(err)
		return
	}

	// Get the agent
	agent := &v1alpha2.Agent{}
	if err := h.KubeClient.Get(r.Context(), client.ObjectKey{
		Namespace: agentNamespace,
		Name:      agentName,
	}, agent); err != nil {
		if k8serrors.IsNotFound(err) {
			w.RespondWithError(errors.NewNotFoundError("Agent not found", err))
			return
		}
		w.RespondWithError(errors.NewInternalServerError("Failed to get Agent", err))
		return
	}

	// Check if agent is declarative and has DSPy config
	if agent.Spec.Type != v1alpha2.AgentType_Declarative || agent.Spec.Declarative.DSPyConfig == nil {
		w.RespondWithError(errors.NewBadRequestError("Agent must be declarative with DSPy configuration", nil))
		return
	}

	// Compile the prompt
	log.V(1).Info("Compiling DSPy prompt")
	result, err := h.compiler.Compile(r.Context(), agent.Spec.Declarative.DSPyConfig)
	if err != nil {
		w.RespondWithError(errors.NewInternalServerError("Failed to compile DSPy prompt", err))
		return
	}

	response := CompileResponse{
		CompiledPrompt: result.CompiledPrompt,
		SignatureDict:  result.SignatureDict,
		ModuleType:     result.ModuleType,
	}

	log.Info("Successfully compiled DSPy prompt")
	data := api.NewResponse(response, "Successfully compiled DSPy prompt", false)
	RespondWithJSON(w, http.StatusOK, data)
}

// HandleStartOptimization handles POST /api/agents/{namespace}/{name}/dspy/optimize
func (h *DSPyHandler) HandleStartOptimization(w ErrorResponseWriter, r *http.Request) {
	log := ctrllog.FromContext(r.Context()).WithName("dspy-handler").WithValues("operation", "optimize")

	agentName, err := GetPathParam(r, "name")
	if err != nil {
		w.RespondWithError(errors.NewBadRequestError("Failed to get name from path", err))
		return
	}

	agentNamespace, err := GetPathParam(r, "namespace")
	if err != nil {
		w.RespondWithError(errors.NewBadRequestError("Failed to get namespace from path", err))
		return
	}

	log = log.WithValues("agentName", agentName, "agentNamespace", agentNamespace)

	if err := Check(h.Authorizer, r, auth.Resource{
		Type: "Agent",
		Name: types.NamespacedName{Namespace: agentNamespace, Name: agentName}.String(),
	}); err != nil {
		w.RespondWithError(err)
		return
	}

	// Get the agent
	agent := &v1alpha2.Agent{}
	if err := h.KubeClient.Get(r.Context(), client.ObjectKey{
		Namespace: agentNamespace,
		Name:      agentName,
	}, agent); err != nil {
		if k8serrors.IsNotFound(err) {
			w.RespondWithError(errors.NewNotFoundError("Agent not found", err))
			return
		}
		w.RespondWithError(errors.NewInternalServerError("Failed to get Agent", err))
		return
	}

	// Check if agent has optimization config
	if agent.Spec.Declarative.DSPyConfig == nil || agent.Spec.Declarative.DSPyConfig.OptimizationConfig == nil {
		w.RespondWithError(errors.NewBadRequestError("Agent must have DSPy optimization configuration", nil))
		return
	}

	agentID := utils.ConvertToPythonIdentifier(utils.GetObjectRef(agent))

	// Create optimization job
	log.V(1).Info("Creating optimization job")
	jobID, err := h.optimizer.CreateOptimizationJob(
		r.Context(),
		agentID,
		agent.Spec.Declarative.DSPyConfig.OptimizationConfig,
	)
	if err != nil {
		w.RespondWithError(errors.NewInternalServerError("Failed to create optimization job", err))
		return
	}

	response := map[string]interface{}{
		"job_id": jobID,
		"status": "pending",
	}

	log.Info("Successfully created optimization job", "jobID", jobID)
	data := api.NewResponse(response, "Successfully started optimization", false)
	RespondWithJSON(w, http.StatusCreated, data)
}

// HandleGetOptimizationJob handles GET /api/agents/{namespace}/{name}/dspy/optimize/{jobId}
func (h *DSPyHandler) HandleGetOptimizationJob(w ErrorResponseWriter, r *http.Request) {
	log := ctrllog.FromContext(r.Context()).WithName("dspy-handler").WithValues("operation", "get-job")

	agentName, err := GetPathParam(r, "name")
	if err != nil {
		w.RespondWithError(errors.NewBadRequestError("Failed to get name from path", err))
		return
	}

	agentNamespace, err := GetPathParam(r, "namespace")
	if err != nil {
		w.RespondWithError(errors.NewBadRequestError("Failed to get namespace from path", err))
		return
	}

	jobID, err := GetPathParam(r, "jobId")
	if err != nil {
		w.RespondWithError(errors.NewBadRequestError("Failed to get jobId from path", err))
		return
	}

	log = log.WithValues("agentName", agentName, "agentNamespace", agentNamespace, "jobID", jobID)

	if err := Check(h.Authorizer, r, auth.Resource{
		Type: "Agent",
		Name: types.NamespacedName{Namespace: agentNamespace, Name: agentName}.String(),
	}); err != nil {
		w.RespondWithError(err)
		return
	}

	// Get optimization job
	job, err := h.optimizer.GetOptimizationJob(r.Context(), jobID)
	if err != nil {
		w.RespondWithError(errors.NewNotFoundError("Optimization job not found", err))
		return
	}

	// Parse metrics if available
	var metrics map[string]interface{}
	if job.Metrics != "" {
		if err := json.Unmarshal([]byte(job.Metrics), &metrics); err != nil {
			log.Error(err, "Failed to parse metrics")
		}
	}

	response := map[string]interface{}{
		"job_id":       job.ID,
		"status":       job.Status,
		"optimizer":    job.Optimizer,
		"started_at":   job.StartedAt,
		"completed_at": job.CompletedAt,
		"metrics":      metrics,
		"error":        job.ErrorMsg,
	}

	log.Info("Successfully retrieved optimization job")
	data := api.NewResponse(response, "Successfully retrieved optimization job", false)
	RespondWithJSON(w, http.StatusOK, data)
}

// HandleListOptimizationJobs handles GET /api/agents/{namespace}/{name}/dspy/optimize
func (h *DSPyHandler) HandleListOptimizationJobs(w ErrorResponseWriter, r *http.Request) {
	log := ctrllog.FromContext(r.Context()).WithName("dspy-handler").WithValues("operation", "list-jobs")

	agentName, err := GetPathParam(r, "name")
	if err != nil {
		w.RespondWithError(errors.NewBadRequestError("Failed to get name from path", err))
		return
	}

	agentNamespace, err := GetPathParam(r, "namespace")
	if err != nil {
		w.RespondWithError(errors.NewBadRequestError("Failed to get namespace from path", err))
		return
	}

	log = log.WithValues("agentName", agentName, "agentNamespace", agentNamespace)

	if err := Check(h.Authorizer, r, auth.Resource{
		Type: "Agent",
		Name: types.NamespacedName{Namespace: agentNamespace, Name: agentName}.String(),
	}); err != nil {
		w.RespondWithError(err)
		return
	}

	// Get agent to get ID
	agent := &v1alpha2.Agent{}
	if err := h.KubeClient.Get(r.Context(), client.ObjectKey{
		Namespace: agentNamespace,
		Name:      agentName,
	}, agent); err != nil {
		w.RespondWithError(errors.NewNotFoundError("Agent not found", err))
		return
	}

	agentID := utils.ConvertToPythonIdentifier(utils.GetObjectRef(agent))

	// List optimization jobs
	jobs, err := h.optimizer.ListOptimizationJobs(r.Context(), agentID)
	if err != nil {
		w.RespondWithError(errors.NewInternalServerError("Failed to list optimization jobs", err))
		return
	}

	log.Info("Successfully listed optimization jobs", "count", len(jobs))
	data := api.NewResponse(jobs, "Successfully listed optimization jobs", false)
	RespondWithJSON(w, http.StatusOK, data)
}
