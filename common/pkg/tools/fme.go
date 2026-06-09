package tools

import (
	"context"
	"fmt"

	config "github.com/harness/mcp-server/common"
	"github.com/harness/mcp-server/common/client"
	"github.com/harness/mcp-server/common/pkg/schemas"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// ListFMEWorkspacesTool creates a tool for listing FME workspaces
func ListFMEWorkspacesTool(config *config.McpServerConfig, fmeService *client.FMEService) (mcp.Tool, server.ToolHandlerFunc) {
	return mcp.NewTool("list_fme_workspaces",
			mcp.WithDescription("List Feature Management & Experimentation (FME) workspaces."),
			mcp.WithReadOnlyHintAnnotation(true),
			mcp.WithIdempotentHintAnnotation(true),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			workspaces, err := fmeService.ListWorkspaces(ctx)
			if err != nil {
				return nil, fmt.Errorf("failed to list FME workspaces: %w", err)
			}

			return schemas.NewStructuredResult(workspaces)
		}
}

// ListFMEEnvironmentsTool creates a tool for listing FME environments for a specific workspace
func ListFMEEnvironmentsTool(config *config.McpServerConfig, fmeService *client.FMEService) (mcp.Tool, server.ToolHandlerFunc) {
	return mcp.NewTool("list_fme_environments",
			mcp.WithDescription("List Feature Management & Experimentation (FME) environments for a specific workspace."),
			mcp.WithString("ws_id",
				mcp.Required(),
				mcp.Description("The workspace ID to list environments for"),
			),
			mcp.WithReadOnlyHintAnnotation(true),
			mcp.WithIdempotentHintAnnotation(true),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			wsID, err := RequiredParam[string](request, "ws_id")
			if err != nil {
				return schemas.NewMissingParamError("ws_id"), nil
			}

			environments, err := fmeService.ListEnvironments(ctx, wsID)
			if err != nil {
				return nil, fmt.Errorf("failed to list FME environments: %w", err)
			}

			return schemas.NewStructuredResult(environments)
		}
}

// ListFMEFeatureFlagsTool creates a tool for listing FME feature flags for a specific workspace
func ListFMEFeatureFlagsTool(config *config.McpServerConfig, fmeService *client.FMEService) (mcp.Tool, server.ToolHandlerFunc) {
	return mcp.NewTool("list_fme_feature_flags",
			mcp.WithDescription("List Feature Management & Experimentation (FME) feature flags for a specific workspace."),
			mcp.WithString("ws_id",
				mcp.Required(),
				mcp.Description("The workspace ID to list feature flags for"),
			),
			mcp.WithReadOnlyHintAnnotation(true),
			mcp.WithIdempotentHintAnnotation(true),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			wsID, err := RequiredParam[string](request, "ws_id")
			if err != nil {
				return schemas.NewMissingParamError("ws_id"), nil
			}

			featureFlags, err := fmeService.ListFeatureFlags(ctx, wsID)
			if err != nil {
				return nil, fmt.Errorf("failed to list FME feature flags: %w", err)
			}

			return schemas.NewStructuredResult(featureFlags)
		}
}

// GetFMEFeatureFlagDefinitionTool creates a tool for getting a specific FME feature flag definition
func GetFMEFeatureFlagDefinitionTool(config *config.McpServerConfig, fmeService *client.FMEService) (mcp.Tool, server.ToolHandlerFunc) {
	return mcp.NewTool("get_fme_feature_flag_definition",
			mcp.WithDescription("Get the definition of a specific Feature Management & Experimentation (FME) feature flag in an environment."),
			mcp.WithString("ws_id",
				mcp.Required(),
				mcp.Description("The workspace ID"),
			),
			mcp.WithString("feature_flag_name",
				mcp.Required(),
				mcp.Description("The name of the feature flag"),
			),
			mcp.WithString("environment_id_or_name",
				mcp.Required(),
				mcp.Description("The environment ID or name"),
			),
			mcp.WithReadOnlyHintAnnotation(true),
			mcp.WithIdempotentHintAnnotation(true),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			wsID, err := RequiredParam[string](request, "ws_id")
			if err != nil {
				return schemas.NewMissingParamError("ws_id"), nil
			}

			flagName, err := RequiredParam[string](request, "feature_flag_name")
			if err != nil {
				return schemas.NewMissingParamError("feature_flag_name"), nil
			}

			envIDOrName, err := RequiredParam[string](request, "environment_id_or_name")
			if err != nil {
				return schemas.NewMissingParamError("environment_id_or_name"), nil
			}

			definition, err := fmeService.GetFeatureFlagDefinition(ctx, wsID, flagName, envIDOrName)
			if err != nil {
				return nil, fmt.Errorf("failed to get FME feature flag definition: %w", err)
			}

			return schemas.NewStructuredResult(definition)
		}
}
