package models

import "time"

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	Name      string    `json:"name"`
	Bio       string    `json:"bio"`
	AvatarURL string    `json:"avatar_url"`
	Location  string    `json:"location"`
	Website   string    `json:"website"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Repository struct {
	ID              int       `json:"id"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	OwnerID         int       `json:"owner_id"`
	OwnerName       string    `json:"owner_name"`
	IsPrivate       bool      `json:"is_private"`
	IsFork          bool      `json:"is_fork"`
	ParentID        *int      `json:"parent_id,omitempty"`
	DefaultBranch   string    `json:"default_branch"`
	StarsCount      int       `json:"stars_count"`
	ForksCount      int       `json:"forks_count"`
	WatchersCount   int       `json:"watchers_count"`
	OpenIssuesCount int       `json:"open_issues_count"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type Commit struct {
	ID           int       `json:"id"`
	RepositoryID int       `json:"repository_id"`
	Hash         string    `json:"hash"`
	Message      string    `json:"message"`
	Author       string    `json:"author"`
	CreatedAt    time.Time `json:"created_at"`
}

type File struct {
	ID           int       `json:"id"`
	RepositoryID int       `json:"repository_id"`
	CommitID     int       `json:"commit_id"`
	Path         string    `json:"path"`
	Content      string    `json:"content"`
	Size         int       `json:"size"`
	CreatedAt    time.Time `json:"created_at"`
}

type Issue struct {
	ID            int       `json:"id"`
	RepositoryID  int       `json:"repository_id"`
	Number        int       `json:"number"`
	Title         string    `json:"title"`
	Body          string    `json:"body"`
	State         string    `json:"state"`
	AuthorID      int       `json:"author_id"`
	AuthorName    string    `json:"author_name"`
	AssigneeID    *int      `json:"assignee_id,omitempty"`
	AssigneeName  string    `json:"assignee_name,omitempty"`
	CommentsCount int       `json:"comments_count"`
	Labels        []Label   `json:"labels,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	ClosedAt      *time.Time `json:"closed_at,omitempty"`
}

type IssueComment struct {
	ID         int       `json:"id"`
	IssueID    int       `json:"issue_id"`
	AuthorID   int       `json:"author_id"`
	AuthorName string    `json:"author_name"`
	Body       string    `json:"body"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type PullRequest struct {
	ID            int       `json:"id"`
	RepositoryID  int       `json:"repository_id"`
	Number        int       `json:"number"`
	Title         string    `json:"title"`
	Body          string    `json:"body"`
	State         string    `json:"state"`
	AuthorID      int       `json:"author_id"`
	AuthorName    string    `json:"author_name"`
	HeadBranch    string    `json:"head_branch"`
	BaseBranch    string    `json:"base_branch"`
	Merged        bool      `json:"merged"`
	MergedAt      *time.Time `json:"merged_at,omitempty"`
	CommentsCount int       `json:"comments_count"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	ClosedAt      *time.Time `json:"closed_at,omitempty"`
}

type PRComment struct {
	ID            int       `json:"id"`
	PullRequestID int       `json:"pull_request_id"`
	AuthorID      int       `json:"author_id"`
	AuthorName    string    `json:"author_name"`
	Body          string    `json:"body"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type Star struct {
	ID           int       `json:"id"`
	UserID       int       `json:"user_id"`
	RepositoryID int       `json:"repository_id"`
	CreatedAt    time.Time `json:"created_at"`
}

type Watcher struct {
	ID           int       `json:"id"`
	UserID       int       `json:"user_id"`
	RepositoryID int       `json:"repository_id"`
	CreatedAt    time.Time `json:"created_at"`
}

type Collaborator struct {
	ID           int       `json:"id"`
	RepositoryID int       `json:"repository_id"`
	UserID       int       `json:"user_id"`
	Username     string    `json:"username"`
	Permission   string    `json:"permission"`
	CreatedAt    time.Time `json:"created_at"`
}

type Branch struct {
	ID           int       `json:"id"`
	RepositoryID int       `json:"repository_id"`
	Name         string    `json:"name"`
	CommitHash   string    `json:"commit_hash"`
	IsDefault    bool      `json:"is_default"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Label struct {
	ID           int       `json:"id"`
	RepositoryID int       `json:"repository_id"`
	Name         string    `json:"name"`
	Color        string    `json:"color"`
	Description  string    `json:"description"`
	CreatedAt    time.Time `json:"created_at"`
}

type Notification struct {
	ID            int       `json:"id"`
	UserID        int       `json:"user_id"`
	Type          string    `json:"type"`
	RepositoryID  *int      `json:"repository_id,omitempty"`
	IssueID       *int      `json:"issue_id,omitempty"`
	PullRequestID *int      `json:"pull_request_id,omitempty"`
	Read          bool      `json:"read"`
	CreatedAt     time.Time `json:"created_at"`
}

type AuthRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type CreateIssueRequest struct {
	Title      string `json:"title"`
	Body       string `json:"body"`
	AssigneeID *int   `json:"assignee_id,omitempty"`
	Labels     []int  `json:"labels,omitempty"`
}

type CreatePullRequestRequest struct {
	Title      string `json:"title"`
	Body       string `json:"body"`
	HeadBranch string `json:"head_branch"`
	BaseBranch string `json:"base_branch"`
}

type SearchResult struct {
	Repositories []Repository `json:"repositories,omitempty"`
	Users        []User       `json:"users,omitempty"`
	Issues       []Issue      `json:"issues,omitempty"`
}
