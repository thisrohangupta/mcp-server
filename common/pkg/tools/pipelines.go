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

func GetPipelineTool(config *config.McpServerConfig, client *client.PipelineService) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool("get_pipeline",
			mcp.WithDescription("Get details of a specific pipeline in a Harness repository. Use list_pipelines (if available) first to find the correct pipeline_id if you're unsure of the exact ID."),
			mcp.WithString("pipeline_id",
				mcp.Required(),
				mcp.Description("The ID of the pipeline"),
			),
			common.WithScope(config, true),
			// MCP v2: Add output schema for typed responses
			mcp.WithOutputSchema[schemas.PipelineOutput](),
			// MCP v2: Add tool annotations for client hints
			mcp.WithReadOnlyHintAnnotation(true),
			mcp.WithIdempotentHintAnnotation(true),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			pipelineID, err := RequiredParam[string](request, "pipeline_id")
			if err != nil {
				return schemas.NewMissingParamError("pipeline_id"), nil
			}
			scope, err := common.FetchScope(ctx, config, request, true)
			if err != nil {
				return schemas.NewScopeError(err.Error()), nil
			}

			data, err := client.Get(ctx, scope, pipelineID)
			if err != nil {
				return nil, fmt.Errorf("failed to get pipeline: %w", err)
			}

			// MCP v2: Return structured output with text fallback
			output := schemas.PipelineOutput{
				Identifier:                    pipelineID,
				YamlPipeline:                  data.Data.YamlPipeline,
				ResolvedTemplatesPipelineYaml: data.Data.ResolvedTemplatesPipelineYaml,
				StoreType:                     data.Data.StoreType,
				ConnectorRef:                  data.Data.ConnectorRef,
				Modules:                       data.Data.Modules,
			}

			return schemas.NewStructuredResult(output)
		}
}

func ListPipelinesTool(config *config.McpServerConfig, client *client.PipelineService) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool("list_pipelines",
			mcp.WithDescription("List pipelines in a Harness repository."),
			mcp.WithString("search_term",
				mcp.Description("Optional search term to filter pipelines"),
			),
			common.WithScope(config, true),
			WithPagination(),
			// MCP v2: Add output schema for typed responses
			mcp.WithOutputSchema[schemas.PipelineListOutput](),
			// MCP v2: Add tool annotations
			mcp.WithReadOnlyHintAnnotation(true),
			mcp.WithIdempotentHintAnnotation(true),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			scope, err := common.FetchScope(ctx, config, request, true)
			if err != nil {
				return schemas.NewScopeError(err.Error()), nil
			}

			page, size, err := FetchPagination(request)
			if err != nil {
				return schemas.NewInvalidParamError("page/size", err.Error()), nil
			}

			searchTerm, err := OptionalParam[string](request, "search_term")
			if err != nil {
				return schemas.NewInvalidParamError("search_term", err.Error()), nil
			}

			opts := &dto.PipelineListOptions{
				SearchTerm: searchTerm,
				PaginationOptions: dto.PaginationOptions{
					Page: page,
					Size: size,
				},
			}

			data, err := client.List(ctx, scope, opts)
			if err != nil {
				return nil, fmt.Errorf("failed to list pipelines: %w", err)
			}

			// Transform to structured output format
			items := make([]schemas.PipelineListItem, 0, len(data.Data.Content))
			for _, p := range data.Data.Content {
				p.ExecutionSummaryInfo.FormatTimestamps()

				item := schemas.PipelineListItem{
					Identifier:    p.Identifier,
					Name:          p.Name,
					Description:   p.Description,
					Tags:          p.Tags,
					NumOfStages:   p.NumOfStages,
					StageNames:    p.StageNames,
					Modules:       p.Modules,
					StoreType:     p.StoreType,
					CreatedAt:     dto.FormatUnixMillisToRFC3339(p.CreatedAt),
					LastUpdatedAt: dto.FormatUnixMillisToRFC3339(p.LastUpdatedAt),
				}

				if p.ExecutionSummaryInfo.LastExecutionId != "" {
					item.LastExecution = &schemas.ExecutionSummary{
						ExecutionID: p.ExecutionSummaryInfo.LastExecutionId,
						Status:      p.ExecutionSummaryInfo.LastExecutionStatus,
						Timestamp:   p.ExecutionSummaryInfo.LastExecutionTsTime,
					}
				}

				items = append(items, item)
			}

			output := schemas.NewPaginatedResult(items, page, size, data.Data.TotalElements)
			return schemas.NewStructuredResult(output)
		}
}

func FetchExecutionURLTool(config *config.McpServerConfig, client *client.PipelineService) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool("fetch_execution_url",
			mcp.WithDescription("Fetch the execution URL for a pipeline execution in Harness."),
			mcp.WithString("pipeline_id",
				mcp.Required(),
				mcp.Description("The ID of the pipeline"),
			),
			mcp.WithString("plan_execution_id",
				mcp.Required(),
				mcp.Description("The ID of the plan execution"),
			),
			common.WithScope(config, true),
			// MCP v2: Add output schema for typed responses
			mcp.WithOutputSchema[schemas.ExecutionURLOutput](),
			// MCP v2: Add tool annotations
			mcp.WithReadOnlyHintAnnotation(true),
			mcp.WithIdempotentHintAnnotation(true),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			pipelineID, err := RequiredParam[string](request, "pipeline_id")
			if err != nil {
				return schemas.NewMissingParamError("pipeline_id"), nil
			}

			planExecutionID, err := RequiredParam[string](request, "plan_execution_id")
			if err != nil {
				return schemas.NewMissingParamError("plan_execution_id"), nil
			}

			scope, err := common.FetchScope(ctx, config, request, true)
			if err != nil {
				return schemas.NewScopeError(err.Error()), nil
			}

			url, err := client.FetchExecutionURL(ctx, scope, pipelineID, planExecutionID)
			if err != nil {
				return nil, fmt.Errorf("failed to fetch execution URL: %w", err)
			}

			output := schemas.URLResult{URL: url}
			return schemas.NewStructuredResult(output)
		}
}

func GetExecutionTool(config *config.McpServerConfig, client *client.PipelineService) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool("get_execution",
			mcp.WithDescription("Get details of a specific pipeline execution in Harness."),
			mcp.WithString("plan_execution_id",
				mcp.Required(),
				mcp.Description("The ID of the plan execution"),
			),
			mcp.WithString("stage_node_id",
				mcp.Description("Optional ID of the stage node to filter the execution details"),
			),
			mcp.WithString("child_stage_node_id",
				mcp.Description("Optional ID of the child stage node to filter the execution details"),
			),
			common.WithScope(config, true),
			// MCP v2: Add output schema for typed responses
			mcp.WithOutputSchema[schemas.ExecutionOutput](),
			// MCP v2: Add tool annotations
			mcp.WithReadOnlyHintAnnotation(true),
			mcp.WithIdempotentHintAnnotation(true),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			planExecutionID, err := RequiredParam[string](request, "plan_execution_id")
			if err != nil {
				return schemas.NewMissingParamError("plan_execution_id"), nil
			}

			// Get optional stage node ID
			stageNodeID, _ := OptionalParam[string](request, "stage_node_id")

			// Get optional child stage node ID
			childStageNodeID, _ := OptionalParam[string](request, "child_stage_node_id")

			scope, err := common.FetchScope(ctx, config, request, true)
			if err != nil {
				return schemas.NewScopeError(err.Error()), nil
			}

			// Pass both stageNodeID and childStageNodeID to the client
			data, err := client.GetExecutionWithLogKeys(ctx, scope, planExecutionID, stageNodeID, childStageNodeID)
			if err != nil {
				return nil, fmt.Errorf("failed to get execution details: %w", err)
			}

			// Format timestamps for the execution
			data.Data.Execution.FormatTimestamps()

			// Transform to structured output
			exec := data.Data.Execution
			output := schemas.ExecutionOutput{
				PlanExecutionID:   exec.PlanExecutionId,
				PipelineID:        exec.PipelineIdentifier,
				Name:              exec.Name,
				Status:            exec.Status,
				OrgIdentifier:     exec.OrgIdentifier,
				ProjectIdentifier: exec.ProjectIdentifier,
				StartTime:         exec.StartTsTime,
				EndTime:           exec.EndTsTime,
				CreatedAt:         exec.CreatedAtTime,
				RunSequence:       exec.RunSequence,
				StagesSummary: &schemas.StagesSummary{
					SuccessfulCount: exec.SuccessfulStagesCount,
					FailedCount:     exec.FailedStagesCount,
					RunningCount:    exec.RunningStagesCount,
					StagesExecuted:  exec.StagesExecuted,
				},
			}

			if exec.ExecutionTriggerInfo != nil {
				output.TriggerInfo = &schemas.TriggerInfo{
					TriggerType: exec.ExecutionTriggerInfo.TriggerType,
					IsRerun:     exec.ExecutionTriggerInfo.IsRerun,
				}
				if exec.ExecutionTriggerInfo.TriggeredBy != nil {
					output.TriggerInfo.TriggeredBy = exec.ExecutionTriggerInfo.TriggeredBy.Identifier
					if exec.ExecutionTriggerInfo.TriggeredBy.ExtraInfo != nil {
						output.TriggerInfo.Email = exec.ExecutionTriggerInfo.TriggeredBy.ExtraInfo.Email
					}
				}
			}

			if exec.FailureInfo.Message != "" {
				output.FailureInfo = &schemas.FailureInfo{
					Message:      exec.FailureInfo.Message,
					FailureTypes: exec.FailureInfo.FailureTypeList,
				}
			}

			// Include execution graph if available
			if data.Data.ExecutionGraph.RootNodeId != "" {
				output.ExecutionGraph = &schemas.ExecutionGraph{
					RootNodeID: data.Data.ExecutionGraph.RootNodeId,
					Nodes:      make(map[string]schemas.ExecutionNode),
				}
				for id, node := range data.Data.ExecutionGraph.NodeMap {
					output.ExecutionGraph.Nodes[id] = schemas.ExecutionNode{
						UUID:       node.Uuid,
						Name:       node.Name,
						Identifier: node.Identifier,
						StepType:   node.StepType,
						Status:     node.Status,
						LogBaseKey: node.LogBaseKey,
					}
				}
			}

			return schemas.NewStructuredResult(output)
		}
}

func ListExecutionsTool(config *config.McpServerConfig, client *client.PipelineService) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool("list_executions",
			mcp.WithDescription("List pipeline executions in a Harness repository."),
			mcp.WithString("search_term",
				mcp.Description("Optional search term to filter executions"),
			),
			mcp.WithString("pipeline_identifier",
				mcp.Description("Optional pipeline identifier to filter executions"),
			),
			// MCP v2: Use enum for constrained values
			mcp.WithString("status",
				mcp.Description("Optional status to filter executions"),
				mcp.Enum("Running", "Success", "Failed", "Aborted", "Expired", "ApprovalWaiting", "Paused"),
			),
			mcp.WithString("branch",
				mcp.Description("Optional branch to filter executions"),
			),
			mcp.WithBoolean("my_deployments",
				mcp.Description("Optional flag to show only my deployments"),
			),
			common.WithScope(config, true),
			WithPagination(),
			// MCP v2: Add output schema for typed responses
			mcp.WithOutputSchema[schemas.ExecutionListOutput](),
			// MCP v2: Add tool annotations
			mcp.WithReadOnlyHintAnnotation(true),
			mcp.WithIdempotentHintAnnotation(true),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			scope, err := common.FetchScope(ctx, config, request, true)
			if err != nil {
				return schemas.NewScopeError(err.Error()), nil
			}

			page, size, err := FetchPagination(request)
			if err != nil {
				return schemas.NewInvalidParamError("page/size", err.Error()), nil
			}

			searchTerm, err := OptionalParam[string](request, "search_term")
			if err != nil {
				return schemas.NewInvalidParamError("search_term", err.Error()), nil
			}

			pipelineIdentifier, err := OptionalParam[string](request, "pipeline_identifier")
			if err != nil {
				return schemas.NewInvalidParamError("pipeline_identifier", err.Error()), nil
			}

			status, err := OptionalParam[string](request, "status")
			if err != nil {
				return schemas.NewInvalidParamError("status", err.Error()), nil
			}

			branch, err := OptionalParam[string](request, "branch")
			if err != nil {
				return schemas.NewInvalidParamError("branch", err.Error()), nil
			}

			myDeployments, err := OptionalParam[bool](request, "my_deployments")
			if err != nil {
				return schemas.NewInvalidParamError("my_deployments", err.Error()), nil
			}

			opts := &dto.PipelineExecutionOptions{
				SearchTerm:         searchTerm,
				PipelineIdentifier: pipelineIdentifier,
				Status:             status,
				Branch:             branch,
				MyDeployments:      myDeployments,
				PaginationOptions: dto.PaginationOptions{
					Page: page,
					Size: size,
				},
			}

			data, err := client.ListExecutions(ctx, scope, opts)
			if err != nil {
				return nil, fmt.Errorf("failed to list pipeline executions: %w", err)
			}

			// Transform to structured output format
			items := make([]schemas.ExecutionListItem, 0, len(data.Data.Content))
			for _, e := range data.Data.Content {
				e.FormatTimestamps()

				item := schemas.ExecutionListItem{
					PlanExecutionID: e.PlanExecutionId,
					PipelineID:      e.PipelineIdentifier,
					Name:            e.Name,
					Status:          e.Status,
					StartTime:       e.StartTsTime,
					EndTime:         e.EndTsTime,
					RunSequence:     e.RunSequence,
					StagesSummary: &schemas.StagesSummary{
						SuccessfulCount: e.SuccessfulStagesCount,
						FailedCount:     e.FailedStagesCount,
						RunningCount:    e.RunningStagesCount,
						StagesExecuted:  e.StagesExecuted,
					},
				}

				if e.ExecutionTriggerInfo != nil {
					item.TriggerInfo = &schemas.TriggerInfo{
						TriggerType: e.ExecutionTriggerInfo.TriggerType,
						IsRerun:     e.ExecutionTriggerInfo.IsRerun,
					}
					if e.ExecutionTriggerInfo.TriggeredBy != nil {
						item.TriggerInfo.TriggeredBy = e.ExecutionTriggerInfo.TriggeredBy.Identifier
					}
				}

				items = append(items, item)
			}

			output := schemas.NewPaginatedResult(items, page, size, data.Data.TotalElements)
			return schemas.NewStructuredResult(output)
		}
}

func ListInputSetsTool(config *config.McpServerConfig, client *client.PipelineService) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool("list_input_sets",
			mcp.WithDescription("List input sets for a pipeline."),
			mcp.WithString("pipeline_identifier",
				mcp.Required(),
				mcp.Description("Pipeline identifier to filter input sets."),
			),
			mcp.WithString("search_term",
				mcp.Description("Optional search term to filter out Input Sets based on name, identifier, tags."),
			),
			common.WithScope(config, true),
			WithPagination(),
			// MCP v2: Add output schema for typed responses
			mcp.WithOutputSchema[schemas.InputSetListOutput](),
			// MCP v2: Add tool annotations
			mcp.WithReadOnlyHintAnnotation(true),
			mcp.WithIdempotentHintAnnotation(true),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			scope, err := common.FetchScope(ctx, config, request, true)
			if err != nil {
				return schemas.NewScopeError(err.Error()), nil
			}

			page, size, err := FetchPagination(request)
			if err != nil {
				return schemas.NewInvalidParamError("page/size", err.Error()), nil
			}

			pipelineIdentifier, err := RequiredParam[string](request, "pipeline_identifier")
			if err != nil {
				return schemas.NewMissingParamError("pipeline_identifier"), nil
			}

			searchTerm, err := OptionalParam[string](request, "search_term")
			if err != nil {
				return schemas.NewInvalidParamError("search_term", err.Error()), nil
			}

			opts := &dto.InputSetListOptions{
				PipelineIdentifier: pipelineIdentifier,
				SearchTerm:         searchTerm,
				PaginationOptions: dto.PaginationOptions{
					Page: page,
					Size: size,
				},
			}

			data, err := client.ListInputSets(ctx, scope, opts)
			if err != nil {
				return nil, fmt.Errorf("failed to list input sets: %w", err)
			}

			// Transform to structured output format
			items := make([]schemas.InputSetListItem, 0, len(data.Data.Content))
			for _, is := range data.Data.Content {
				items = append(items, schemas.InputSetListItem{
					Identifier:         is.Identifier,
					Name:               is.Name,
					Description:        is.Description,
					PipelineIdentifier: is.PipelineIdentifier,
					InputSetType:       is.InputSetType,
					Tags:               is.Tags,
					IsOutdated:         is.IsOutdated,
				})
			}

			output := schemas.NewPaginatedResult(items, page, size, data.Data.TotalItems)
			return schemas.NewStructuredResult(output)
		}
}

func GetInputSetTool(config *config.McpServerConfig, client *client.PipelineService) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool("get_input_set",
			mcp.WithDescription("Get details of a specific input set for a pipeline in Harness."),
			mcp.WithString("pipeline_identifier",
				mcp.Required(),
				mcp.Description("The identifier of the pipeline."),
			),
			mcp.WithString("input_set_identifier",
				mcp.Required(),
				mcp.Description("The identifier of the input set."),
			),
			common.WithScope(config, true),
			// MCP v2: Add output schema for typed responses
			mcp.WithOutputSchema[schemas.InputSetOutput](),
			// MCP v2: Add tool annotations
			mcp.WithReadOnlyHintAnnotation(true),
			mcp.WithIdempotentHintAnnotation(true),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			pipelineIdentifier, err := RequiredParam[string](request, "pipeline_identifier")
			if err != nil {
				return schemas.NewMissingParamError("pipeline_identifier"), nil
			}

			inputSetIdentifier, err := RequiredParam[string](request, "input_set_identifier")
			if err != nil {
				return schemas.NewMissingParamError("input_set_identifier"), nil
			}

			scope, err := common.FetchScope(ctx, config, request, true)
			if err != nil {
				return schemas.NewScopeError(err.Error()), nil
			}

			data, err := client.GetInputSet(ctx, scope, pipelineIdentifier, inputSetIdentifier)
			if err != nil {
				return nil, fmt.Errorf("failed to get input set: %w", err)
			}

			output := schemas.InputSetOutput{
				Identifier:         data.Data.Identifier,
				Name:               data.Data.Name,
				Description:        data.Data.Description,
				PipelineIdentifier: data.Data.PipelineIdentifier,
				InputSetYaml:       data.Data.InputSetYaml,
				Tags:               data.Data.Tags,
				IsOutdated:         data.Data.Outdated,
			}

			return schemas.NewStructuredResult(output)
		}
}

func GetPipelineSummaryTool(config *config.McpServerConfig, client *client.PipelineService) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool("get_pipeline_summary",
			mcp.WithDescription("Provides a concise summary of a pipeline's overall structure and execution info highlighting key aspects rather than detailed pipeline definition such as pipeline yaml, external references, etc."),
			mcp.WithString("pipeline_id",
				mcp.Required(),
				mcp.Description("Identifier of the pipeline."),
			),
			mcp.WithBoolean("get_metadata_only",
				mcp.Description("Whether to only fetch metadata without full pipeline details."),
			),
			common.WithScope(config, true),
			// MCP v2: Add output schema for typed responses
			mcp.WithOutputSchema[schemas.PipelineSummaryOutput](),
			// MCP v2: Add tool annotations
			mcp.WithReadOnlyHintAnnotation(true),
			mcp.WithIdempotentHintAnnotation(true),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			pipelineID, err := RequiredParam[string](request, "pipeline_id")
			if err != nil {
				return schemas.NewMissingParamError("pipeline_id"), nil
			}

			scope, err := common.FetchScope(ctx, config, request, true)
			if err != nil {
				return schemas.NewScopeError(err.Error()), nil
			}

			getMetadataOnly, err := OptionalParam[bool](request, "get_metadata_only")
			if err != nil {
				return schemas.NewInvalidParamError("get_metadata_only", err.Error()), nil
			}

			data, err := client.GetPipelineSummary(ctx, scope, pipelineID, getMetadataOnly)
			if err != nil {
				return nil, fmt.Errorf("failed to get pipeline summary: %w", err)
			}

			summary := data.Data
			output := schemas.PipelineSummaryOutput{
				Identifier:    summary.Identifier,
				Name:          summary.Name,
				Description:   summary.Description,
				Tags:          summary.Tags,
				NumOfStages:   summary.NumOfStages,
				StageNames:    summary.StageNames,
				Modules:       summary.Modules,
				Version:       summary.Version,
				CreatedAt:     dto.FormatUnixMillisToRFC3339(summary.CreatedAt),
				LastUpdatedAt: dto.FormatUnixMillisToRFC3339(summary.LastUpdatedAt),
			}

			if summary.ExecutionSummaryInfo != nil {
				summary.ExecutionSummaryInfo.FormatTimestamps()
				output.LastExecution = &schemas.ExecutionSummary{
					ExecutionID: summary.ExecutionSummaryInfo.LastExecutionId,
					Status:      summary.ExecutionSummaryInfo.LastExecutionStatus,
					Timestamp:   summary.ExecutionSummaryInfo.LastExecutionTsTime,
				}
			}

			return schemas.NewStructuredResult(output)
		}
}

type ActionData struct {
	Actions []struct {
		Text   string `json:"text"`
		Action string `json:"action"`
		Data   struct {
			PageName string            `json:"pageName"`
			Metadata map[string]string `json:"metadata"`
		} `json:"data"`
	} `json:"actions"`
}

func ListTriggersTool(config *config.McpServerConfig, client *client.PipelineService) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool("list_triggers",
			mcp.WithDescription("List triggers in a Harness pipeline."),
			mcp.WithString("target_identifier",
				mcp.Required(),
				mcp.Description("Identifier of the target pipeline."),
			),
			mcp.WithString("search_term",
				mcp.Description("Optional search term to filter triggers based on name, identifier, tags."),
			),
			common.WithScope(config, true),
			WithPagination(),
			// MCP v2: Add output schema for typed responses
			mcp.WithOutputSchema[schemas.TriggerListOutput](),
			// MCP v2: Add tool annotations
			mcp.WithReadOnlyHintAnnotation(true),
			mcp.WithIdempotentHintAnnotation(true),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			scope, err := common.FetchScope(ctx, config, request, true)
			if err != nil {
				return schemas.NewScopeError(err.Error()), nil
			}

			page, size, err := FetchPagination(request)
			if err != nil {
				return schemas.NewInvalidParamError("page/size", err.Error()), nil
			}

			searchTerm, err := OptionalParam[string](request, "search_term")
			if err != nil {
				return schemas.NewInvalidParamError("search_term", err.Error()), nil
			}

			targetIdentifier, err := RequiredParam[string](request, "target_identifier")
			if err != nil {
				return schemas.NewMissingParamError("target_identifier"), nil
			}

			opts := &dto.TriggerListOptions{
				SearchTerm:       searchTerm,
				TargetIdentifier: targetIdentifier,
				PaginationOptions: dto.PaginationOptions{
					Page: page,
					Size: size,
				},
			}

			data, err := client.ListTriggers(ctx, scope, opts)
			if err != nil {
				return nil, fmt.Errorf("failed to list triggers: %w", err)
			}

			// Transform to structured output format
			items := make([]schemas.TriggerListItem, 0, len(data.Data.Content))
			for _, t := range data.Data.Content {
				item := schemas.TriggerListItem{
					Identifier:         t.Identifier,
					Name:               t.Name,
					Description:        t.Description,
					Type:               t.Type,
					Enabled:            t.Enabled,
					PipelineIdentifier: t.PipelineIdentifier,
				}

				if t.TriggerStatus != nil {
					item.Status = &schemas.TriggerStatus{
						Status:         t.TriggerStatus.Status,
						DetailMessages: t.TriggerStatus.DetailMessages,
					}
				}

				items = append(items, item)
			}

			output := schemas.NewPaginatedResult(items, page, size, data.Data.TotalElements)
			return schemas.NewStructuredResult(output)
		}
}
