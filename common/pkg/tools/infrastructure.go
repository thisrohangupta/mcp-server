package tools

import (
	"context"
	"fmt"

	config "github.com/harness/mcp-server/common"
	"github.com/harness/mcp-server/common/client"
	"github.com/harness/mcp-server/common/client/dto"
	"github.com/harness/mcp-server/common/pkg/common"
	"github.com/harness/mcp-server/common/pkg/schemas"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// ListInfrastructuresTool creates a tool for listing infrastructures
// https://apidocs.harness.io/tag/Infrastructures#operation/getInfrastructureList
func ListInfrastructuresTool(config *config.McpServerConfig, client *client.InfrastructureClient) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool("list_infrastructures",
			mcp.WithDescription("List infrastructure definitions in Harness."),
			mcp.WithString("deploymentType",
				mcp.Description("Optional filter for deployment type (e.g., Kubernetes, ECS)"),
			),
			mcp.WithString("environmentIdentifier",
				mcp.Required(),
				mcp.Description("Filter for environment"),
			),
			mcp.WithString("sort",
				mcp.Description("Optional field to sort by (e.g., name)"),
			),
			mcp.WithString("order",
				mcp.Description("Optional sort order (asc or desc)"),
				mcp.Enum("asc", "desc"),
			),
			mcp.WithNumber("page",
				mcp.DefaultNumber(0),
				mcp.Description("Page number for pagination (0-based)"),
			),
			mcp.WithNumber("limit",
				mcp.DefaultNumber(5),
				mcp.Max(20),
				mcp.Description("Number of infrastructures per page"),
			),
			common.WithScope(config, false),
			mcp.WithOutputSchema[schemas.InfrastructureListOutput](),
			mcp.WithReadOnlyHintAnnotation(true),
			mcp.WithIdempotentHintAnnotation(true),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			scope, err := common.FetchScope(ctx, config, request, false)
			if err != nil {
				return schemas.NewScopeError(err.Error()), nil
			}

			opts := &dto.InfrastructureOptions{}

			// Handle pagination
			page, err := OptionalParam[float64](request, "page")
			if err != nil {
				return schemas.NewInvalidParamError("page", err.Error()), nil
			}
			if page >= 0 {
				opts.Page = int(page)
			}

			limit, err := OptionalParam[float64](request, "limit")
			if err != nil {
				return schemas.NewInvalidParamError("limit", err.Error()), nil
			}
			if limit > 0 {
				opts.Limit = int(limit)
			}

			// Handle filters
			deploymentType, err := OptionalParam[string](request, "deploymentType")
			if err != nil {
				return schemas.NewInvalidParamError("deploymentType", err.Error()), nil
			}
			if deploymentType != "" {
				opts.DeploymentType = deploymentType
			}

			environmentIdentifier, err := RequiredParam[string](request, "environmentIdentifier")
			if err != nil {
				return schemas.NewMissingParamError("environmentIdentifier"), nil
			}
			opts.EnvironmentIdentifier = environmentIdentifier

			// Handle sorting
			sort, err := OptionalParam[string](request, "sort")
			if err != nil {
				return schemas.NewInvalidParamError("sort", err.Error()), nil
			}
			if sort != "" {
				opts.Sort = sort
			}

			order, err := OptionalParam[string](request, "order")
			if err != nil {
				return schemas.NewInvalidParamError("order", err.Error()), nil
			}
			if order != "" {
				opts.Order = order
			}

			infrastructures, totalCount, err := client.List(ctx, scope, opts)
			if err != nil {
				return nil, fmt.Errorf("failed to list infrastructures: %w", err)
			}

			// Transform to structured output
			items := make([]schemas.InfrastructureListItem, 0, len(infrastructures))
			for _, infra := range infrastructures {
				items = append(items, schemas.InfrastructureListItem{
					Identifier:            infra.ID,
					Name:                  infra.Name,
					Description:           infra.Description,
					EnvironmentIdentifier: infra.EnvironmentRef,
					DeploymentType:        infra.DeploymentType,
					Type:                  infra.Type,
					Tags:                  infra.Tags,
				})
			}

			output := schemas.NewPaginatedResult(items, opts.Page, opts.Limit, totalCount)
			return schemas.NewStructuredResult(output)
		}
}

// MoveInfrastructureConfigsTool creates a tool for moving configurations between infrastructures
// https://apidocs.harness.io/tag/Infrastructures#operation/moveInfraConfigs
func MoveInfrastructureConfigsTool(config *config.McpServerConfig, client *client.InfrastructureClient) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool("move_infrastructure_configs",
			mcp.WithDescription("Move infrastructure YAML from inline to remote or vice versa in Harness."),
			mcp.WithString("infra_identifier",
				mcp.Required(),
				mcp.Description("Infrastructure identifier to move"),
			),
			mcp.WithString("environment_identifier",
				mcp.Required(),
				mcp.Description("Environment identifier for the infrastructure"),
			),
			mcp.WithString("move_config_type",
				mcp.Required(),
				mcp.Description("Specifies the direction of the move operation"),
				mcp.Enum("INLINE_TO_REMOTE", "REMOTE_TO_INLINE"),
			),
			mcp.WithString("org_identifier",
				mcp.Description("Organization identifier"),
			),
			mcp.WithString("project_identifier",
				mcp.Description("Project identifier"),
			),
			mcp.WithString("connector_ref",
				mcp.Description("Identifier of connector needed for operations on the entity"),
			),
			mcp.WithString("repo_name",
				mcp.Description("Name of the repository"),
			),
			mcp.WithString("branch",
				mcp.Description("Name of the branch"),
			),
			mcp.WithString("file_path",
				mcp.Description("File path of the entity"),
			),
			mcp.WithString("commit_msg",
				mcp.Description("Commit message to use for the merge commit"),
			),
			mcp.WithBoolean("is_new_branch",
				mcp.Description("Checks the new branch"),
			),
			mcp.WithString("base_branch",
				mcp.Description("Name of the default branch"),
			),
			mcp.WithBoolean("is_harness_code_repo",
				mcp.Description("Is Harness code repo enabled"),
			),
			mcp.WithOutputSchema[schemas.MoveConfigOutput](),
			mcp.WithDestructiveHintAnnotation(true),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			// Get scope
			scope, err := common.FetchScope(ctx, config, request, false)
			if err != nil {
				return schemas.NewScopeError(err.Error()), nil
			}

			// Extract required parameters
			infraIdentifier, err := RequiredParam[string](request, "infra_identifier")
			if err != nil {
				return schemas.NewMissingParamError("infra_identifier"), nil
			}

			environmentIdentifier, err := RequiredParam[string](request, "environment_identifier")
			if err != nil {
				return schemas.NewMissingParamError("environment_identifier"), nil
			}

			moveConfigTypeStr, err := RequiredParam[string](request, "move_config_type")
			if err != nil {
				return schemas.NewMissingParamError("move_config_type"), nil
			}

			// Validate move config type
			if moveConfigTypeStr != string(dto.InlineToRemote) && moveConfigTypeStr != string(dto.RemoteToInline) {
				return schemas.NewInvalidParamError("move_config_type", "must be INLINE_TO_REMOTE or REMOTE_TO_INLINE"), nil
			}
			moveConfigType := dto.MoveConfigType(moveConfigTypeStr)

			// Extract optional parameters
			orgIdentifier, err := OptionalParam[string](request, "org_identifier")
			if err != nil {
				return schemas.NewInvalidParamError("org_identifier", err.Error()), nil
			}

			projectIdentifier, err := OptionalParam[string](request, "project_identifier")
			if err != nil {
				return schemas.NewInvalidParamError("project_identifier", err.Error()), nil
			}

			connectorRef, err := OptionalParam[string](request, "connector_ref")
			if err != nil {
				return schemas.NewInvalidParamError("connector_ref", err.Error()), nil
			}

			repoName, err := OptionalParam[string](request, "repo_name")
			if err != nil {
				return schemas.NewInvalidParamError("repo_name", err.Error()), nil
			}

			branch, err := OptionalParam[string](request, "branch")
			if err != nil {
				return schemas.NewInvalidParamError("branch", err.Error()), nil
			}

			filePath, err := OptionalParam[string](request, "file_path")
			if err != nil {
				return schemas.NewInvalidParamError("file_path", err.Error()), nil
			}

			commitMsg, err := OptionalParam[string](request, "commit_msg")
			if err != nil {
				return schemas.NewInvalidParamError("commit_msg", err.Error()), nil
			}

			isNewBranch, err := OptionalParam[bool](request, "is_new_branch")
			if err != nil {
				return schemas.NewInvalidParamError("is_new_branch", err.Error()), nil
			}

			baseBranch, err := OptionalParam[string](request, "base_branch")
			if err != nil {
				return schemas.NewInvalidParamError("base_branch", err.Error()), nil
			}

			isHarnessCodeRepo, err := OptionalParam[bool](request, "is_harness_code_repo")
			if err != nil {
				return schemas.NewInvalidParamError("is_harness_code_repo", err.Error()), nil
			}

			// Create move request
			moveRequest := &dto.MoveInfraConfigsRequest{
				InfraIdentifier:       infraIdentifier,
				EnvironmentIdentifier: environmentIdentifier,
				AccountIdentifier:     scope.AccountID,
				OrgIdentifier:         orgIdentifier,
				ProjectIdentifier:     projectIdentifier,
				ConnectorRef:          connectorRef,
				RepoName:              repoName,
				Branch:                branch,
				FilePath:              filePath,
				CommitMsg:             commitMsg,
				MoveConfigType:        moveConfigType,
			}

			// Set boolean pointers if values were provided
			if isNewBranchProvided, ok := request.GetArguments()["is_new_branch"]; ok && isNewBranchProvided != nil {
				moveRequest.IsNewBranch = &isNewBranch
			}

			if isHarnessCodeRepoProvided, ok := request.GetArguments()["is_harness_code_repo"]; ok && isHarnessCodeRepoProvided != nil {
				moveRequest.IsHarnessCodeRepo = &isHarnessCodeRepo
			}

			// Add the base branch if provided
			if baseBranch != "" {
				moveRequest.BaseBranch = baseBranch
			}

			// Execute the move operation
			response, err := client.MoveConfigs(ctx, scope, moveRequest)
			if err != nil {
				return nil, fmt.Errorf("failed to move infrastructure configurations: %w", err)
			}

			output := schemas.MoveConfigOutput{
				Success: response.Data.Success,
			}
			if response.Data.Success {
				output.Message = "Infrastructure configuration moved successfully"
			}

			return schemas.NewStructuredResult(output)
		}
}
