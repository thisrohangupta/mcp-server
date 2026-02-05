# MCP v2 Structured Outputs Migration Guide

This guide explains how to migrate existing tools to use MCP v2 structured outputs.

## Overview

MCP v2 introduces structured outputs with typed schemas, enabling:
- Client-side validation
- Better UI rendering
- Downstream processing without parsing
- Type safety across the protocol

## Migration Steps

### Step 1: Define Output Schema Type

Create a Go struct that represents your tool's output in `common/pkg/schemas/`:

```go
// In common/pkg/schemas/yourmodule.go

// YourToolOutput is the structured output for your_tool.
type YourToolOutput struct {
    ID          string `json:"id"`
    Name        string `json:"name"`
    Description string `json:"description,omitempty"`
    // ... other fields
}

// For list operations, use the generic PaginatedResult
type YourItemListOutput = PaginatedResult[YourListItem]
```

### Step 2: Add Output Schema to Tool Definition

Update the tool definition to include the output schema and annotations:

```go
// Before
func YourTool(config *config.McpServerConfig, client *client.YourService) (tool mcp.Tool, handler server.ToolHandlerFunc) {
    return mcp.NewTool("your_tool",
            mcp.WithDescription("..."),
            mcp.WithString("param_id", mcp.Required(), mcp.Description("...")),
            common.WithScope(config, true),
        ),
        handler
}

// After
func YourTool(config *config.McpServerConfig, client *client.YourService) (tool mcp.Tool, handler server.ToolHandlerFunc) {
    return mcp.NewTool("your_tool",
            mcp.WithDescription("..."),
            mcp.WithString("param_id", mcp.Required(), mcp.Description("...")),
            common.WithScope(config, true),
            // MCP v2: Add output schema for typed responses
            mcp.WithOutputSchema[schemas.YourToolOutput](),
            // MCP v2: Add tool annotations
            mcp.WithReadOnlyHintAnnotation(true),      // For read-only operations
            mcp.WithIdempotentHintAnnotation(true),    // For idempotent operations
            // mcp.WithDestructiveHintAnnotation(true), // For destructive operations
        ),
        handler
}
```

### Step 3: Update Error Handling

Replace `mcp.NewToolResultError()` with structured error helpers:

```go
// Before
if err != nil {
    return mcp.NewToolResultError(err.Error()), nil
}

// After - For missing parameters
if err != nil {
    return schemas.NewMissingParamError("param_name"), nil
}

// After - For invalid parameters
if err != nil {
    return schemas.NewInvalidParamError("param_name", err.Error()), nil
}

// After - For scope errors
if err != nil {
    return schemas.NewScopeError(err.Error()), nil
}
```

### Step 4: Return Structured Results

Replace `mcp.NewToolResultText()` with `schemas.NewStructuredResult()`:

```go
// Before
r, err := json.Marshal(data)
if err != nil {
    return nil, fmt.Errorf("failed to marshal: %w", err)
}
return mcp.NewToolResultText(string(r)), nil

// After - For single items
output := schemas.YourToolOutput{
    ID:          data.ID,
    Name:        data.Name,
    Description: data.Description,
}
return schemas.NewStructuredResult(output)

// After - For paginated lists
items := make([]schemas.YourListItem, 0, len(data.Content))
for _, d := range data.Content {
    items = append(items, schemas.YourListItem{
        ID:   d.ID,
        Name: d.Name,
    })
}
output := schemas.NewPaginatedResult(items, page, size, data.TotalElements)
return schemas.NewStructuredResult(output)
```

## Tool Annotations Reference

| Annotation | Use When |
|------------|----------|
| `WithReadOnlyHintAnnotation(true)` | Tool only reads data, no side effects |
| `WithIdempotentHintAnnotation(true)` | Same request returns same result |
| `WithDestructiveHintAnnotation(true)` | Tool deletes or destroys resources |
| `WithOpenWorldHintAnnotation(true)` | Tool interacts with external systems |

## Best Practices

1. **Use enums for constrained values**:
   ```go
   mcp.WithString("status",
       mcp.Enum("Running", "Success", "Failed", "Aborted"),
       mcp.Description("Execution status filter"),
   )
   ```

2. **Use snake_case for JSON field names** in output schemas (MCP convention)

3. **Include both structured and text content** for backward compatibility (handled automatically by `NewStructuredResult`)

4. **Format timestamps as RFC3339** strings, not Unix timestamps

5. **Use nullable pointers** for optional nested objects

## Files to Update

When migrating a tool:
1. `common/pkg/schemas/{module}.go` - Add output schema types
2. `common/pkg/tools/{module}.go` - Update tool definition and handler
3. Run `make format` to ensure proper formatting

## Example Migration: Pipeline Tools

See `common/pkg/tools/pipelines.go` for complete examples of migrated tools:
- `get_pipeline` - Single item response
- `list_pipelines` - Paginated list response
- `list_executions` - List with enum parameter
- `get_execution` - Complex nested response
