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

// GetServiceTool creates a tool for getting details of a specific service
// https://apidocs.harness.io/tag/Services#operation/getServiceV2
func GetServiceTool(config *config.McpServerConfig, client *client.ServiceClient) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool("get_service",
			mcp.WithDescription("Get details of a specific service in Harness."),
			mcp.WithString("service_identifier",
				mcp.Required(),
				mcp.Description("The identifier of the service"),
			),
			common.WithScope(config, false),
			mcp.WithOutputSchema[schemas.ServiceOutput](),
			mcp.WithReadOnlyHintAnnotation(true),
			mcp.WithIdempotentHintAnnotation(true),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			serviceIdentifier, err := RequiredParam[string](request, "service_identifier")
			if err != nil {
				return schemas.NewMissingParamError("service_identifier"), nil
			}

			scope, err := common.FetchScope(ctx, config, request, false)
			if err != nil {
				return schemas.NewScopeError(err.Error()), nil
			}

			data, err := client.Get(ctx, scope, serviceIdentifier)
			if err != nil {
				return nil, fmt.Errorf("failed to get service: %w", err)
			}

			return schemas.NewStructuredResult(data)
		}
}

// ListServicesTool creates a tool for listing services
// https://apidocs.harness.io/tag/Services#operation/getServiceList
func ListServicesTool(config *config.McpServerConfig, client *client.ServiceClient) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool("list_services",
			mcp.WithDescription("List services in Harness."),
			mcp.WithString("sort",
				mcp.Description("Optional field to sort by (e.g., name)"),
			),
			mcp.WithString("order",
				mcp.Description("Optional sort order"),
				mcp.Enum("asc", "desc"),
			),
			mcp.WithNumber("page",
				mcp.DefaultNumber(0),
				mcp.Description("Page number for pagination (0-based)"),
			),
			mcp.WithNumber("limit",
				mcp.DefaultNumber(5),
				mcp.Max(20),
				mcp.Description("Number of services per page"),
			),
			common.WithScope(config, false),
			mcp.WithOutputSchema[schemas.ServiceListOutput](),
			mcp.WithReadOnlyHintAnnotation(true),
			mcp.WithIdempotentHintAnnotation(true),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			scope, err := common.FetchScope(ctx, config, request, false)
			if err != nil {
				return schemas.NewScopeError(err.Error()), nil
			}

			opts := &dto.ServiceOptions{}

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

			// Handle other optional parameters
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

			services, totalCount, err := client.List(ctx, scope, opts)
			if err != nil {
				return nil, fmt.Errorf("failed to list services: %w", err)
			}

			// Transform to structured output
			items := make([]schemas.ServiceListItem, 0, len(services))
			for _, s := range services {
				items = append(items, schemas.ServiceListItem{
					Identifier:  s.ID,
					Name:        s.Name,
					Description: s.Description,
					Tags:        s.Tags,
				})
			}

			output := schemas.NewPaginatedResult(items, opts.Page, opts.Limit, totalCount)
			return schemas.NewStructuredResult(output)
		}
}
