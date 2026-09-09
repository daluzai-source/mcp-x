package lists

type List struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description,omitempty"`
	Private       bool   `json:"private"`
	OwnerID       string `json:"owner_id,omitempty"`
	MemberCount   int    `json:"member_count"`
	FollowerCount int    `json:"follower_count"`
	CreatedAt     string `json:"created_at,omitempty"`
	URL           string `json:"url"`
}

type User struct {
	ID             string `json:"id"`
	Username       string `json:"username"`
	Name           string `json:"name,omitempty"`
	Verified       bool   `json:"verified"`
	FollowersCount int    `json:"followers_count"`
	URL            string `json:"url"`
}

type Post struct {
	ID             string `json:"id"`
	Text           string `json:"text"`
	AuthorID       string `json:"author_id,omitempty"`
	AuthorUsername string `json:"author_username,omitempty"`
	CreatedAt      string `json:"created_at,omitempty"`
	Likes          int    `json:"likes"`
	Reposts        int    `json:"reposts"`
	URL            string `json:"url"`
}

type CreateRequest struct {
	Name        string `json:"name" jsonschema:"list name, at most 25 characters" validate:"required,max=25"`
	Description string `json:"description,omitempty" jsonschema:"list description, at most 100 characters" validate:"omitempty,max=100"`
	Private     bool   `json:"private,omitempty" jsonschema:"private lists are visible only to you"`
	TimeoutMs   int64  `json:"timeout_ms,omitempty"`
}

type UpdateRequest struct {
	ID          string  `json:"id" jsonschema:"list id" validate:"required"`
	Name        *string `json:"name,omitempty" jsonschema:"new list name, at most 25 characters"`
	Description *string `json:"description,omitempty" jsonschema:"new description, at most 100 characters"`
	Private     *bool   `json:"private,omitempty"`
	TimeoutMs   int64   `json:"timeout_ms,omitempty"`
}

type ListIDRequest struct {
	ID        string `json:"id" jsonschema:"list id" validate:"required"`
	TimeoutMs int64  `json:"timeout_ms,omitempty"`
}

type ListPageRequest struct {
	ID              string `json:"id" jsonschema:"list id" validate:"required"`
	MaxResults      int    `json:"max_results,omitempty" jsonschema:"items per page, 1..100" validate:"omitempty,gte=1,lte=100"`
	PaginationToken string `json:"pagination_token,omitempty" jsonschema:"next_token from the previous call; omit for the first page"`
	TimeoutMs       int64  `json:"timeout_ms,omitempty"`
}

type OwnedRequest struct {
	Username        string `json:"username,omitempty" jsonschema:"whose lists to return; omit for the key owner"`
	MaxResults      int    `json:"max_results,omitempty" validate:"omitempty,gte=1,lte=100"`
	PaginationToken string `json:"pagination_token,omitempty"`
	TimeoutMs       int64  `json:"timeout_ms,omitempty"`
}

type SelfPageRequest struct {
	MaxResults      int    `json:"max_results,omitempty" validate:"omitempty,gte=1,lte=100"`
	PaginationToken string `json:"pagination_token,omitempty"`
	TimeoutMs       int64  `json:"timeout_ms,omitempty"`
}

type MemberRequest struct {
	ID        string `json:"id" jsonschema:"list id" validate:"required"`
	Username  string `json:"username" jsonschema:"X username without the @" validate:"required"`
	TimeoutMs int64  `json:"timeout_ms,omitempty"`
}

type ListResponse struct {
	List List `json:"list"`
}

type ListsResponse struct {
	Lists     []List `json:"lists"`
	Count     int    `json:"count"`
	NextToken string `json:"next_token,omitempty"`
}

type UsersResponse struct {
	Users     []User `json:"users"`
	Count     int    `json:"count"`
	NextToken string `json:"next_token,omitempty"`
}

type PostsResponse struct {
	Posts     []Post `json:"posts"`
	Count     int    `json:"count"`
	NextToken string `json:"next_token,omitempty"`
}

type ActionResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}
