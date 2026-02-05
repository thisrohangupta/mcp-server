package schemas

// PaginatedResult is a generic wrapper for paginated responses.
// This schema is used by all list operations.
type PaginatedResult[T any] struct {
	Items      []T             `json:"items"`
	Pagination PaginationInfo  `json:"pagination"`
	Metadata   *ResultMetadata `json:"metadata,omitempty"`
}

// PaginationInfo contains pagination details for list responses.
type PaginationInfo struct {
	Page          int     `json:"page"`
	PageSize      int     `json:"page_size"`
	TotalItems    int     `json:"total_items"`
	TotalPages    int     `json:"total_pages"`
	HasMore       bool    `json:"has_more"`
	NextPageToken *string `json:"next_page_token,omitempty"`
}

// ResultMetadata contains optional metadata about the result.
type ResultMetadata struct {
	RequestID     string `json:"request_id,omitempty"`
	ExecutionTime string `json:"execution_time,omitempty"`
}

// SingleResult is a wrapper for single-item responses.
type SingleResult[T any] struct {
	Item     T               `json:"item"`
	Metadata *ResultMetadata `json:"metadata,omitempty"`
}

// URLResult is a simple result containing a URL.
type URLResult struct {
	URL      string          `json:"url"`
	Metadata *ResultMetadata `json:"metadata,omitempty"`
}

// Scope represents Harness scope information included in responses.
type Scope struct {
	AccountID string `json:"account_id"`
	OrgID     string `json:"org_id,omitempty"`
	ProjectID string `json:"project_id,omitempty"`
}

// NewPaginatedResult creates a paginated result from list data.
func NewPaginatedResult[T any](items []T, page, pageSize, totalItems int) PaginatedResult[T] {
	totalPages := 0
	if pageSize > 0 {
		totalPages = (totalItems + pageSize - 1) / pageSize
	}

	return PaginatedResult[T]{
		Items: items,
		Pagination: PaginationInfo{
			Page:       page,
			PageSize:   pageSize,
			TotalItems: totalItems,
			TotalPages: totalPages,
			HasMore:    page < totalPages-1,
		},
	}
}

// NewSingleResult creates a single-item result wrapper.
func NewSingleResult[T any](item T) SingleResult[T] {
	return SingleResult[T]{Item: item}
}
