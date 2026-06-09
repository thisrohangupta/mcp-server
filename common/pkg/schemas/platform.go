package schemas

// ConnectorOutput represents a connector in structured output format.
// Used by: get_connector_details
type ConnectorOutput struct {
	Identifier        string               `json:"identifier"`
	Name              string               `json:"name"`
	Description       string               `json:"description,omitempty"`
	Type              string               `json:"type"`
	AccountIdentifier string               `json:"account_identifier,omitempty"`
	OrgIdentifier     string               `json:"org_identifier,omitempty"`
	ProjectIdentifier string               `json:"project_identifier,omitempty"`
	Tags              map[string]string    `json:"tags,omitempty"`
	Status            *ConnectorStatusInfo `json:"status,omitempty"`
	CreatedAt         string               `json:"created_at,omitempty"`
	LastModifiedAt    string               `json:"last_modified_at,omitempty"`
	HarnessManaged    bool                 `json:"harness_managed,omitempty"`
	Spec              map[string]any       `json:"spec,omitempty"`
}

// ConnectorStatusInfo contains connector connectivity status.
type ConnectorStatusInfo struct {
	Status          string `json:"status"`
	ErrorSummary    string `json:"error_summary,omitempty"`
	TestedAt        string `json:"tested_at,omitempty"`
	LastConnectedAt string `json:"last_connected_at,omitempty"`
}

// ConnectorListItem represents a connector in list responses.
// Used by: list_connectors
type ConnectorListItem struct {
	Identifier        string               `json:"identifier"`
	Name              string               `json:"name"`
	Description       string               `json:"description,omitempty"`
	Type              string               `json:"type"`
	OrgIdentifier     string               `json:"org_identifier,omitempty"`
	ProjectIdentifier string               `json:"project_identifier,omitempty"`
	Tags              map[string]string    `json:"tags,omitempty"`
	Status            *ConnectorStatusInfo `json:"status,omitempty"`
	CreatedAt         string               `json:"created_at,omitempty"`
	LastModifiedAt    string               `json:"last_modified_at,omitempty"`
}

// ConnectorListOutput is the structured output for list_connectors.
type ConnectorListOutput = PaginatedResult[ConnectorListItem]

// ConnectorCatalogueOutput represents the connector catalogue.
// Used by: list_connector_catalogue
type ConnectorCatalogueOutput struct {
	Categories []ConnectorCategory `json:"categories"`
}

// ConnectorCategory represents a category in the connector catalogue.
type ConnectorCategory struct {
	Category   string          `json:"category"`
	Connectors []ConnectorType `json:"connectors"`
}

// ConnectorType represents a connector type in the catalogue.
type ConnectorType struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description,omitempty"`
}

// SecretOutput represents a secret in structured output format.
// Used by: get_secret
type SecretOutput struct {
	Identifier              string            `json:"identifier"`
	Name                    string            `json:"name"`
	Description             string            `json:"description,omitempty"`
	Type                    string            `json:"type"`
	OrgIdentifier           string            `json:"org_identifier,omitempty"`
	ProjectIdentifier       string            `json:"project_identifier,omitempty"`
	Tags                    map[string]string `json:"tags,omitempty"`
	SecretManagerIdentifier string            `json:"secret_manager_identifier,omitempty"`
	CreatedAt               string            `json:"created_at,omitempty"`
	UpdatedAt               string            `json:"updated_at,omitempty"`
}

// SecretListItem represents a secret in list responses.
// Used by: list_secrets
type SecretListItem struct {
	Identifier              string            `json:"identifier"`
	Name                    string            `json:"name"`
	Description             string            `json:"description,omitempty"`
	Type                    string            `json:"type"`
	OrgIdentifier           string            `json:"org_identifier,omitempty"`
	ProjectIdentifier       string            `json:"project_identifier,omitempty"`
	Tags                    map[string]string `json:"tags,omitempty"`
	SecretManagerIdentifier string            `json:"secret_manager_identifier,omitempty"`
}

// SecretListOutput is the structured output for list_secrets.
type SecretListOutput = PaginatedResult[SecretListItem]

// ServiceOutput represents a service in structured output format.
// Used by: get_service
type ServiceOutput struct {
	Identifier        string            `json:"identifier"`
	Name              string            `json:"name"`
	Description       string            `json:"description,omitempty"`
	OrgIdentifier     string            `json:"org_identifier,omitempty"`
	ProjectIdentifier string            `json:"project_identifier,omitempty"`
	Tags              map[string]string `json:"tags,omitempty"`
	Yaml              string            `json:"yaml,omitempty"`
	CreatedAt         string            `json:"created_at,omitempty"`
	LastModifiedAt    string            `json:"last_modified_at,omitempty"`
}

// ServiceListItem represents a service in list responses.
// Used by: list_services
type ServiceListItem struct {
	Identifier        string            `json:"identifier"`
	Name              string            `json:"name"`
	Description       string            `json:"description,omitempty"`
	OrgIdentifier     string            `json:"org_identifier,omitempty"`
	ProjectIdentifier string            `json:"project_identifier,omitempty"`
	Tags              map[string]string `json:"tags,omitempty"`
}

// ServiceListOutput is the structured output for list_services.
type ServiceListOutput = PaginatedResult[ServiceListItem]

// EnvironmentOutput represents an environment in structured output format.
// Used by: get_environment
type EnvironmentOutput struct {
	Identifier        string            `json:"identifier"`
	Name              string            `json:"name"`
	Description       string            `json:"description,omitempty"`
	Type              string            `json:"type,omitempty"`
	OrgIdentifier     string            `json:"org_identifier,omitempty"`
	ProjectIdentifier string            `json:"project_identifier,omitempty"`
	Tags              map[string]string `json:"tags,omitempty"`
	Color             string            `json:"color,omitempty"`
	Yaml              string            `json:"yaml,omitempty"`
	CreatedAt         string            `json:"created_at,omitempty"`
	LastModifiedAt    string            `json:"last_modified_at,omitempty"`
}

// EnvironmentListItem represents an environment in list responses.
// Used by: list_environments
type EnvironmentListItem struct {
	Identifier        string            `json:"identifier"`
	Name              string            `json:"name"`
	Description       string            `json:"description,omitempty"`
	Type              string            `json:"type,omitempty"`
	OrgIdentifier     string            `json:"org_identifier,omitempty"`
	ProjectIdentifier string            `json:"project_identifier,omitempty"`
	Tags              map[string]string `json:"tags,omitempty"`
}

// EnvironmentListOutput is the structured output for list_environments.
type EnvironmentListOutput = PaginatedResult[EnvironmentListItem]

// MoveConfigOutput represents the result of a move config operation.
// Used by: move_environment_configs
type MoveConfigOutput struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

// AuditEventOutput represents an audit event in structured output format.
type AuditEventOutput struct {
	AuditID            string           `json:"audit_id"`
	ResourceType       string           `json:"resource_type"`
	ResourceIdentifier string           `json:"resource_identifier,omitempty"`
	Action             string           `json:"action"`
	Module             string           `json:"module,omitempty"`
	Timestamp          string           `json:"timestamp"`
	AuthenticationInfo *AuthInfo        `json:"authentication_info,omitempty"`
	HttpRequestInfo    *HttpRequestInfo `json:"http_request_info,omitempty"`
}

// AuthInfo contains authentication details for audit events.
type AuthInfo struct {
	Principal     string `json:"principal,omitempty"`
	PrincipalType string `json:"principal_type,omitempty"`
	Email         string `json:"email,omitempty"`
}

// HttpRequestInfo contains HTTP request details for audit events.
type HttpRequestInfo struct {
	RequestMethod string `json:"request_method,omitempty"`
	ClientIP      string `json:"client_ip,omitempty"`
}

// AuditEventListOutput is the structured output for list_audit_events.
type AuditEventListOutput = PaginatedResult[AuditEventOutput]

// TemplateOutput represents a template in structured output format.
type TemplateOutput struct {
	Identifier         string            `json:"identifier"`
	Name               string            `json:"name"`
	Description        string            `json:"description,omitempty"`
	VersionLabel       string            `json:"version_label,omitempty"`
	TemplateEntityType string            `json:"template_entity_type,omitempty"`
	ChildType          string            `json:"child_type,omitempty"`
	OrgIdentifier      string            `json:"org_identifier,omitempty"`
	ProjectIdentifier  string            `json:"project_identifier,omitempty"`
	Tags               map[string]string `json:"tags,omitempty"`
	Yaml               string            `json:"yaml,omitempty"`
	StableTemplate     bool              `json:"stable_template,omitempty"`
	CreatedAt          string            `json:"created_at,omitempty"`
	LastUpdatedAt      string            `json:"last_updated_at,omitempty"`
}

// TemplateListItem represents a template in list responses.
type TemplateListItem struct {
	Identifier         string            `json:"identifier"`
	Name               string            `json:"name"`
	Description        string            `json:"description,omitempty"`
	VersionLabel       string            `json:"version_label,omitempty"`
	TemplateEntityType string            `json:"template_entity_type,omitempty"`
	OrgIdentifier      string            `json:"org_identifier,omitempty"`
	ProjectIdentifier  string            `json:"project_identifier,omitempty"`
	Tags               map[string]string `json:"tags,omitempty"`
	StableTemplate     bool              `json:"stable_template,omitempty"`
}

// TemplateListOutput is the structured output for list_templates.
type TemplateListOutput = PaginatedResult[TemplateListItem]

// InfrastructureOutput represents an infrastructure in structured output format.
type InfrastructureOutput struct {
	Identifier            string            `json:"identifier"`
	Name                  string            `json:"name"`
	Description           string            `json:"description,omitempty"`
	EnvironmentIdentifier string            `json:"environment_identifier,omitempty"`
	DeploymentType        string            `json:"deployment_type,omitempty"`
	Type                  string            `json:"type,omitempty"`
	OrgIdentifier         string            `json:"org_identifier,omitempty"`
	ProjectIdentifier     string            `json:"project_identifier,omitempty"`
	Tags                  map[string]string `json:"tags,omitempty"`
	Yaml                  string            `json:"yaml,omitempty"`
}

// InfrastructureListItem represents an infrastructure in list responses.
type InfrastructureListItem struct {
	Identifier            string            `json:"identifier"`
	Name                  string            `json:"name"`
	Description           string            `json:"description,omitempty"`
	EnvironmentIdentifier string            `json:"environment_identifier,omitempty"`
	DeploymentType        string            `json:"deployment_type,omitempty"`
	Type                  string            `json:"type,omitempty"`
	Tags                  map[string]string `json:"tags,omitempty"`
}

// InfrastructureListOutput is the structured output for list_infrastructures.
type InfrastructureListOutput = PaginatedResult[InfrastructureListItem]
