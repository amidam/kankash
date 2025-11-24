package push

type Payload struct {
	ObjectKind   string         `json:"object_kind"`
	EventName    string         `json:"event_name"`
	Before       string         `json:"before"`
	After        string         `json:"after"`
	Ref          string         `json:"ref"`
	RefProtected bool           `json:"ref_protected"`
	CheckoutSHA  string         `json:"checkout_sha"`
	Message      string         `json:"message"`
	UserID       int64          `json:"user_id"`
	UserName     string         `json:"user_name"`
	UserUsername string         `json:"user_username"`
	UserEmail    string         `json:"user_email"`
	UserAvatar   string         `json:"user_avatar"`
	ProjectID    int64          `json:"project_id"`
	Project      ProjectInfo    `json:"project"`
	Commits      []Commit       `json:"commits"`
	TotalCommits int64          `json:"total_commits_count"`
	PushOptions  map[string]any `json:"push_options"`
	Repository   RepositoryInfo `json:"repository"`
}

type ProjectInfo struct {
	ID                int64   `json:"id"`
	Name              string  `json:"name"`
	Description       string  `json:"description"`
	WebURL            string  `json:"web_url"`
	AvatarURL         *string `json:"avatar_url"`
	GitSSHURL         string  `json:"git_ssh_url"`
	GitHTTPURL        string  `json:"git_http_url"`
	Namespace         string  `json:"namespace"`
	VisibilityLevel   int     `json:"visibility_level"`
	PathWithNamespace string  `json:"path_with_namespace"`
	DefaultBranch     string  `json:"default_branch"`
	CIConfigPath      *string `json:"ci_config_path"`
	Homepage          string  `json:"homepage"`
	URL               string  `json:"url"`
	SSHURL            string  `json:"ssh_url"`
	HTTPURL           string  `json:"http_url"`
}

type Commit struct {
	ID        string       `json:"id"`
	Message   string       `json:"message"`
	Title     string       `json:"title"`
	Timestamp string       `json:"timestamp"`
	URL       string       `json:"url"`
	Author    CommitAuthor `json:"author"`
	Added     []string     `json:"added"`
	Modified  []string     `json:"modified"`
	Removed   []string     `json:"removed"`
}

type CommitAuthor struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type RepositoryInfo struct {
	Name            string `json:"name"`
	URL             string `json:"url"`
	Description     string `json:"description"`
	Homepage        string `json:"homepage"`
	GitHTTPURL      string `json:"git_http_url"`
	GitSSHURL       string `json:"git_ssh_url"`
	VisibilityLevel int    `json:"visibility_level"`
}
