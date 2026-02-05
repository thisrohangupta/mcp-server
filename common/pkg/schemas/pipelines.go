package schemas

// PipelineOutput represents a pipeline in structured output format.
// Used by: get_pipeline
type PipelineOutput struct {
	Identifier                    string            `json:"identifier"`
	Name                          string            `json:"name,omitempty"`
	Description                   string            `json:"description,omitempty"`
	Tags                          map[string]string `json:"tags,omitempty"`
	OrgIdentifier                 string            `json:"org_identifier,omitempty"`
	ProjectIdentifier             string            `json:"project_identifier,omitempty"`
	YamlPipeline                  string            `json:"yaml_pipeline,omitempty"`
	ResolvedTemplatesPipelineYaml string            `json:"resolved_templates_pipeline_yaml,omitempty"`
	StoreType                     string            `json:"store_type,omitempty"`
	ConnectorRef                  string            `json:"connector_ref,omitempty"`
	Modules                       []string          `json:"modules,omitempty"`
	GitDetails                    *GitDetails       `json:"git_details,omitempty"`
}

// GitDetails contains git repository details for a pipeline.
type GitDetails struct {
	Branch     string `json:"branch,omitempty"`
	FilePath   string `json:"file_path,omitempty"`
	RepoName   string `json:"repo_name,omitempty"`
	RootFolder string `json:"root_folder,omitempty"`
}

// PipelineListItem represents a pipeline in list responses.
// Used by: list_pipelines
type PipelineListItem struct {
	Identifier    string            `json:"identifier"`
	Name          string            `json:"name"`
	Description   string            `json:"description,omitempty"`
	Tags          map[string]string `json:"tags,omitempty"`
	NumOfStages   int               `json:"num_of_stages"`
	StageNames    []string          `json:"stage_names,omitempty"`
	Modules       []string          `json:"modules,omitempty"`
	StoreType     string            `json:"store_type,omitempty"`
	CreatedAt     string            `json:"created_at,omitempty"`
	LastUpdatedAt string            `json:"last_updated_at,omitempty"`
	LastExecution *ExecutionSummary `json:"last_execution,omitempty"`
}

// ExecutionSummary contains summary information about pipeline executions.
type ExecutionSummary struct {
	ExecutionID string `json:"execution_id,omitempty"`
	Status      string `json:"status,omitempty"`
	Timestamp   string `json:"timestamp,omitempty"`
}

// PipelineListOutput is the structured output for list_pipelines.
type PipelineListOutput = PaginatedResult[PipelineListItem]

// ExecutionOutput represents a pipeline execution in structured output format.
// Used by: get_execution
type ExecutionOutput struct {
	PlanExecutionID   string          `json:"plan_execution_id"`
	PipelineID        string          `json:"pipeline_id"`
	Name              string          `json:"name,omitempty"`
	Status            string          `json:"status"`
	OrgIdentifier     string          `json:"org_identifier,omitempty"`
	ProjectIdentifier string          `json:"project_identifier,omitempty"`
	StartTime         string          `json:"start_time,omitempty"`
	EndTime           string          `json:"end_time,omitempty"`
	CreatedAt         string          `json:"created_at,omitempty"`
	RunSequence       int32           `json:"run_sequence,omitempty"`
	TriggerInfo       *TriggerInfo    `json:"trigger_info,omitempty"`
	FailureInfo       *FailureInfo    `json:"failure_info,omitempty"`
	StagesSummary     *StagesSummary  `json:"stages_summary,omitempty"`
	ExecutionGraph    *ExecutionGraph `json:"execution_graph,omitempty"`
}

// TriggerInfo contains information about what triggered the execution.
type TriggerInfo struct {
	TriggerType string `json:"trigger_type,omitempty"`
	TriggeredBy string `json:"triggered_by,omitempty"`
	Email       string `json:"email,omitempty"`
	IsRerun     bool   `json:"is_rerun,omitempty"`
}

// FailureInfo contains information about execution failures.
type FailureInfo struct {
	Message      string   `json:"message,omitempty"`
	FailureTypes []string `json:"failure_types,omitempty"`
}

// StagesSummary provides a summary of stage execution counts.
type StagesSummary struct {
	SuccessfulCount int      `json:"successful_count"`
	FailedCount     int      `json:"failed_count"`
	RunningCount    int      `json:"running_count"`
	StagesExecuted  []string `json:"stages_executed,omitempty"`
}

// ExecutionGraph represents the execution flow graph.
type ExecutionGraph struct {
	RootNodeID string                   `json:"root_node_id,omitempty"`
	Nodes      map[string]ExecutionNode `json:"nodes,omitempty"`
}

// ExecutionNode represents a node in the execution graph.
type ExecutionNode struct {
	UUID       string `json:"uuid,omitempty"`
	Name       string `json:"name,omitempty"`
	Identifier string `json:"identifier,omitempty"`
	StepType   string `json:"step_type,omitempty"`
	Status     string `json:"status,omitempty"`
	LogBaseKey string `json:"log_base_key,omitempty"`
}

// ExecutionListItem represents an execution in list responses.
// Used by: list_executions
type ExecutionListItem struct {
	PlanExecutionID string         `json:"plan_execution_id"`
	PipelineID      string         `json:"pipeline_id"`
	Name            string         `json:"name,omitempty"`
	Status          string         `json:"status"`
	StartTime       string         `json:"start_time,omitempty"`
	EndTime         string         `json:"end_time,omitempty"`
	RunSequence     int32          `json:"run_sequence,omitempty"`
	TriggerInfo     *TriggerInfo   `json:"trigger_info,omitempty"`
	StagesSummary   *StagesSummary `json:"stages_summary,omitempty"`
}

// ExecutionListOutput is the structured output for list_executions.
type ExecutionListOutput = PaginatedResult[ExecutionListItem]

// ExecutionURLOutput is the structured output for fetch_execution_url.
type ExecutionURLOutput = URLResult

// InputSetOutput represents an input set in structured output format.
// Used by: get_input_set
type InputSetOutput struct {
	Identifier         string            `json:"identifier"`
	Name               string            `json:"name"`
	Description        string            `json:"description,omitempty"`
	PipelineIdentifier string            `json:"pipeline_identifier"`
	InputSetYaml       string            `json:"input_set_yaml,omitempty"`
	Tags               map[string]string `json:"tags,omitempty"`
	IsOutdated         bool              `json:"is_outdated,omitempty"`
	CreatedAt          string            `json:"created_at,omitempty"`
	LastUpdatedAt      string            `json:"last_updated_at,omitempty"`
}

// InputSetListItem represents an input set in list responses.
// Used by: list_input_sets
type InputSetListItem struct {
	Identifier         string            `json:"identifier"`
	Name               string            `json:"name"`
	Description        string            `json:"description,omitempty"`
	PipelineIdentifier string            `json:"pipeline_identifier"`
	InputSetType       string            `json:"input_set_type,omitempty"`
	Tags               map[string]string `json:"tags,omitempty"`
	IsOutdated         bool              `json:"is_outdated,omitempty"`
}

// InputSetListOutput is the structured output for list_input_sets.
type InputSetListOutput = PaginatedResult[InputSetListItem]

// TriggerOutput represents a trigger in structured output format.
type TriggerOutput struct {
	Identifier         string            `json:"identifier"`
	Name               string            `json:"name"`
	Description        string            `json:"description,omitempty"`
	Type               string            `json:"type"`
	Enabled            bool              `json:"enabled"`
	PipelineIdentifier string            `json:"pipeline_identifier"`
	Tags               map[string]string `json:"tags,omitempty"`
	WebhookURL         string            `json:"webhook_url,omitempty"`
	Status             *TriggerStatus    `json:"status,omitempty"`
}

// TriggerStatus contains status information for a trigger.
type TriggerStatus struct {
	Status         string   `json:"status,omitempty"`
	DetailMessages []string `json:"detail_messages,omitempty"`
}

// TriggerListItem represents a trigger in list responses.
// Used by: list_triggers
type TriggerListItem struct {
	Identifier         string         `json:"identifier"`
	Name               string         `json:"name"`
	Description        string         `json:"description,omitempty"`
	Type               string         `json:"type"`
	Enabled            bool           `json:"enabled"`
	PipelineIdentifier string         `json:"pipeline_identifier"`
	Status             *TriggerStatus `json:"status,omitempty"`
}

// TriggerListOutput is the structured output for list_triggers.
type TriggerListOutput = PaginatedResult[TriggerListItem]

// PipelineSummaryOutput is the structured output for get_pipeline_summary.
type PipelineSummaryOutput struct {
	Identifier    string            `json:"identifier"`
	Name          string            `json:"name"`
	Description   string            `json:"description,omitempty"`
	Tags          map[string]string `json:"tags,omitempty"`
	NumOfStages   int               `json:"num_of_stages"`
	StageNames    []string          `json:"stage_names,omitempty"`
	Modules       []string          `json:"modules,omitempty"`
	Version       int               `json:"version,omitempty"`
	CreatedAt     string            `json:"created_at,omitempty"`
	LastUpdatedAt string            `json:"last_updated_at,omitempty"`
	LastExecution *ExecutionSummary `json:"last_execution,omitempty"`
}
