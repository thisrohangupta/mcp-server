package schemas

// FMEWorkspaceOutput represents an FME workspace.
type FMEWorkspaceOutput struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// FMEWorkspaceListOutput is the structured output for list_fme_workspaces.
type FMEWorkspaceListOutput = PaginatedResult[FMEWorkspaceOutput]

// FMEEnvironmentOutput represents an FME environment.
type FMEEnvironmentOutput struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	WorkspaceID string `json:"workspace_id,omitempty"`
}

// FMEEnvironmentListOutput is the structured output for list_fme_environments.
type FMEEnvironmentListOutput = PaginatedResult[FMEEnvironmentOutput]

// FMEFeatureFlagOutput represents an FME feature flag.
type FMEFeatureFlagOutput struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	TrafficType string `json:"traffic_type,omitempty"`
	CreatedAt   string `json:"created_at,omitempty"`
}

// FMEFeatureFlagListOutput is the structured output for list_fme_feature_flags.
type FMEFeatureFlagListOutput = PaginatedResult[FMEFeatureFlagOutput]

// FMEFeatureFlagDefinitionOutput represents an FME feature flag definition.
type FMEFeatureFlagDefinitionOutput struct {
	Name              string            `json:"name"`
	TrafficType       string            `json:"traffic_type,omitempty"`
	Treatments        []FMETreatment    `json:"treatments,omitempty"`
	DefaultTreatment  string            `json:"default_treatment,omitempty"`
	DefaultRule       *FMERule          `json:"default_rule,omitempty"`
	Rules             []FMERule         `json:"rules,omitempty"`
	TrafficAllocation []FMETrafficSplit `json:"traffic_allocation,omitempty"`
}

// FMETreatment represents a treatment in an FME feature flag.
type FMETreatment struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// FMERule represents a targeting rule in an FME feature flag.
type FMERule struct {
	Buckets []FMEBucket `json:"buckets,omitempty"`
}

// FMEBucket represents a treatment bucket in an FME rule.
type FMEBucket struct {
	Treatment string `json:"treatment"`
	Size      int    `json:"size"`
}

// FMETrafficSplit represents traffic allocation for a treatment.
type FMETrafficSplit struct {
	Treatment  string `json:"treatment"`
	Percentage int    `json:"percentage"`
}
