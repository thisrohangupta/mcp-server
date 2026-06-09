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

// GetRepositoryTool creates a tool for getting a specific repository
func GetRepositoryTool(config *config.McpServerConfig, client *client.RepositoryService) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool("get_repository",
			mcp.WithDescription("Get details of a specific repository in Harness."),
			mcp.WithString("repo_identifier",
				mcp.Required(),
				mcp.Description("The identifier of the repository"),
			),
			common.WithScope(config, false),
			mcp.WithOutputSchema[schemas.RepositoryOutput](),
			mcp.WithReadOnlyHintAnnotation(true),
			mcp.WithIdempotentHintAnnotation(true),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			repoIdentifier, err := RequiredParam[string](request, "repo_identifier")
			if err != nil {
				return schemas.NewMissingParamError("repo_identifier"), nil
			}

			scope, err := common.FetchScope(ctx, config, request, false)
			if err != nil {
				return schemas.NewScopeError(err.Error()), nil
			}

			data, err := client.Get(ctx, scope, repoIdentifier)
			if err != nil {
				return nil, fmt.Errorf("failed to get repository: %w", err)
			}

			return schemas.NewStructuredResult(data)
		}
}

// ListRepositoriesTool creates a tool for listing repositories
func ListRepositoriesTool(config *config.McpServerConfig, client *client.RepositoryService) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool("list_repositories",
			mcp.WithDescription("List repositories in Harness."),
			mcp.WithString("query",
				mcp.Description("Optional search term to filter repositories"),
			),
			mcp.WithString("sort",
				mcp.Description("Optional field to sort by (e.g., identifier)"),
			),
			mcp.WithString("order",
				mcp.Description("Optional sort order"),
				mcp.Enum("asc", "desc"),
			),
			mcp.WithNumber("page",
				mcp.DefaultNumber(1),
				mcp.Description("Page number for pagination"),
			),
			mcp.WithNumber("limit",
				mcp.DefaultNumber(5),
				mcp.Max(20),
				mcp.Description("Number of items per page"),
			),
			common.WithScope(config, false),
			mcp.WithOutputSchema[schemas.RepositoryListOutput](),
			mcp.WithReadOnlyHintAnnotation(true),
			mcp.WithIdempotentHintAnnotation(true),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			scope, err := common.FetchScope(ctx, config, request, false)
			if err != nil {
				return schemas.NewScopeError(err.Error()), nil
			}

			opts := &dto.RepositoryOptions{}

			// Handle pagination
			page, err := OptionalParam[float64](request, "page")
			if err != nil {
				return schemas.NewInvalidParamError("page", err.Error()), nil
			}
			if page > 0 {
				opts.Page = int(page)
			}

			limit, err := OptionalParam[float64](request, "limit")
			if err != nil {
				return schemas.NewInvalidParamError("limit", err.Error()), nil
			}
			if limit > 0 {
				opts.Limit = int(limit)
			}

			// Handle other optional parameters
			query, err := OptionalParam[string](request, "query")
			if err != nil {
				return schemas.NewInvalidParamError("query", err.Error()), nil
			}
			if query != "" {
				opts.Query = query
			}

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

			data, err := client.List(ctx, scope, opts)
			if err != nil {
				return nil, fmt.Errorf("failed to list repositories: %w", err)
			}

			return schemas.NewStructuredResult(data)
		}
}
