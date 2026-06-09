package schemas

// RepositoryOutput represents a repository in structured output format.
// Used by: get_repository
type RepositoryOutput struct {
	ID             string `json:"id"`
	Identifier     string `json:"identifier"`
	Name           string `json:"name,omitempty"`
	Description    string `json:"description,omitempty"`
	Path           string `json:"path,omitempty"`
	DefaultBranch  string `json:"default_branch,omitempty"`
	IsPublic       bool   `json:"is_public,omitempty"`
	IsEmpty        bool   `json:"is_empty,omitempty"`
	GitURL         string `json:"git_url,omitempty"`
	NumForks       int    `json:"num_forks,omitempty"`
	NumPulls       int    `json:"num_pulls,omitempty"`
	NumClosedPulls int    `json:"num_closed_pulls,omitempty"`
	NumOpenPulls   int    `json:"num_open_pulls,omitempty"`
	NumMergedPulls int    `json:"num_merged_pulls,omitempty"`
	CreatedAt      string `json:"created_at,omitempty"`
	UpdatedAt      string `json:"updated_at,omitempty"`
}

// RepositoryListItem represents a repository in list responses.
// Used by: list_repositories
type RepositoryListItem struct {
	ID            string `json:"id"`
	Identifier    string `json:"identifier"`
	Name          string `json:"name,omitempty"`
	Description   string `json:"description,omitempty"`
	Path          string `json:"path,omitempty"`
	DefaultBranch string `json:"default_branch,omitempty"`
	IsPublic      bool   `json:"is_public,omitempty"`
}

// RepositoryListOutput is the structured output for list_repositories.
type RepositoryListOutput = PaginatedResult[RepositoryListItem]

// PullRequestOutput represents a pull request in structured output format.
// Used by: get_pull_request
type PullRequestOutput struct {
	Number           int               `json:"number"`
	Title            string            `json:"title"`
	State            string            `json:"state"`
	Description      string            `json:"description,omitempty"`
	SourceBranch     string            `json:"source_branch"`
	TargetBranch     string            `json:"target_branch"`
	SourceRepoID     int               `json:"source_repo_id,omitempty"`
	Author           *PullRequestUser  `json:"author,omitempty"`
	IsDraft          bool              `json:"is_draft,omitempty"`
	MergeCheckStatus string            `json:"merge_check_status,omitempty"`
	MergeMethod      string            `json:"merge_method,omitempty"`
	MergeSHA         string            `json:"merge_sha,omitempty"`
	Stats            *PullRequestStats `json:"stats,omitempty"`
	CreatedAt        string            `json:"created_at,omitempty"`
	UpdatedAt        string            `json:"updated_at,omitempty"`
	MergedAt         string            `json:"merged_at,omitempty"`
}

// PullRequestUser represents a user associated with a pull request.
type PullRequestUser struct {
	ID          int    `json:"id,omitempty"`
	UID         string `json:"uid,omitempty"`
	DisplayName string `json:"display_name,omitempty"`
	Email       string `json:"email,omitempty"`
}

// PullRequestStats contains statistics for a pull request.
type PullRequestStats struct {
	Commits         int `json:"commits,omitempty"`
	FilesChanged    int `json:"files_changed,omitempty"`
	Additions       int `json:"additions,omitempty"`
	Deletions       int `json:"deletions,omitempty"`
	Conversations   int `json:"conversations,omitempty"`
	UnresolvedCount int `json:"unresolved_count,omitempty"`
}

// PullRequestListItem represents a pull request in list responses.
// Used by: list_pull_requests
type PullRequestListItem struct {
	Number       int              `json:"number"`
	Title        string           `json:"title"`
	State        string           `json:"state"`
	SourceBranch string           `json:"source_branch"`
	TargetBranch string           `json:"target_branch"`
	Author       *PullRequestUser `json:"author,omitempty"`
	IsDraft      bool             `json:"is_draft,omitempty"`
	CreatedAt    string           `json:"created_at,omitempty"`
}

// PullRequestListOutput is the structured output for list_pull_requests.
type PullRequestListOutput = PaginatedResult[PullRequestListItem]

// PullRequestCreateOutput represents the result of creating a pull request.
// Used by: create_pull_request
type PullRequestCreateOutput struct {
	Number       int    `json:"number"`
	Title        string `json:"title"`
	State        string `json:"state"`
	SourceBranch string `json:"source_branch"`
	TargetBranch string `json:"target_branch"`
	URL          string `json:"url,omitempty"`
}

// PullRequestActionOutput represents the result of a PR action.
// Used by: merge_pull_request, close_pull_request, reopen_pull_request
type PullRequestActionOutput struct {
	Number  int    `json:"number"`
	State   string `json:"state"`
	Message string `json:"message,omitempty"`
}
