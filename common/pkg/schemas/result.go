// Package schemas provides structured output schemas for MCP tools.
// This package implements MCP best practices for typed tool outputs,
// enabling client-side validation and better downstream processing.
package schemas

import (
	"encoding/json"

	"github.com/mark3labs/mcp-go/mcp"
)

// NewStructuredResult creates a CallToolResult with both structured content
// and text fallback for backward compatibility with older MCP clients.
//
// Per MCP specification: "A tool returning structured content SHOULD also
// return functionally equivalent unstructured content."
func NewStructuredResult(data any) (*mcp.CallToolResult, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(string(jsonBytes)),
		},
		StructuredContent: data,
	}, nil
}

// NewStructuredResultWithText creates a CallToolResult with structured content
// and a custom text representation (useful when you want different text output).
func NewStructuredResultWithText(data any, text string) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(text),
		},
		StructuredContent: data,
	}
}

// ToolError represents a machine-readable error response.
// This enables clients to programmatically handle errors.
type ToolError struct {
	Code    string `json:"code"`              // Machine-readable error code
	Message string `json:"message"`           // Human-readable message
	Details any    `json:"details,omitempty"` // Additional context
}

// Error codes for common validation failures
const (
	ErrCodeMissingParam     = "MISSING_REQUIRED_PARAMETER"
	ErrCodeInvalidParam     = "INVALID_PARAMETER"
	ErrCodeScopeRequired    = "SCOPE_REQUIRED"
	ErrCodeNotFound         = "NOT_FOUND"
	ErrCodePermissionDenied = "PERMISSION_DENIED"
	ErrCodeRateLimited      = "RATE_LIMITED"
	ErrCodeInternal         = "INTERNAL_ERROR"
)

// NewToolError creates a structured error result.
func NewToolError(code, message string, details any) *mcp.CallToolResult {
	errData := ToolError{
		Code:    code,
		Message: message,
		Details: details,
	}

	jsonBytes, _ := json.Marshal(errData)
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(string(jsonBytes)),
		},
		StructuredContent: errData,
		IsError:           true,
	}
}

// NewMissingParamError creates an error for missing required parameters.
func NewMissingParamError(paramName string) *mcp.CallToolResult {
	return NewToolError(
		ErrCodeMissingParam,
		"missing required parameter: "+paramName,
		map[string]string{"parameter": paramName},
	)
}

// NewInvalidParamError creates an error for invalid parameter values.
func NewInvalidParamError(paramName, reason string) *mcp.CallToolResult {
	return NewToolError(
		ErrCodeInvalidParam,
		"invalid parameter "+paramName+": "+reason,
		map[string]string{"parameter": paramName, "reason": reason},
	)
}

// NewScopeError creates an error for missing scope parameters.
func NewScopeError(message string) *mcp.CallToolResult {
	return NewToolError(ErrCodeScopeRequired, message, nil)
}
