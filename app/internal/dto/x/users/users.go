package users

type User struct {
	ID              string `json:"id"`
	Username        string `json:"username"`
	Name            string `json:"name,omitempty"`
	Description     string `json:"description,omitempty"`
	ProfileImageURL string `json:"profile_image_url,omitempty"`
	Verified        bool   `json:"verified"`
	Protected       bool   `json:"protected"`
	FollowersCount  int    `json:"followers_count"`
	FollowingCount  int    `json:"following_count"`
	PostCount       int    `json:"post_count"`
	CreatedAt       string `json:"created_at,omitempty"`
	URL             string `json:"url"`
}

type LookupRequest struct {
	Usernames []string `json:"usernames,omitempty" jsonschema:"X usernames without the @, up to 100 in one call; set this or ids, not both" validate:"omitempty,max=100,dive,required"`
	IDs       []string `json:"ids,omitempty" jsonschema:"numeric user ids, up to 100 in one call; set this or usernames, not both" validate:"omitempty,max=100,dive,required"`
	TimeoutMs int64    `json:"timeout_ms,omitempty"`
}

type LookupResponse struct {
	Users    []User   `json:"users"`
	Count    int      `json:"count"`
	NotFound []string `json:"not_found,omitempty"`
}

type MeRequest struct {
	TimeoutMs int64 `json:"timeout_ms,omitempty"`
}

type MeResponse struct {
	User User `json:"user"`
}

type ListRequest struct {
	Username        string `json:"username" jsonschema:"X username without the @" validate:"required"`
	MaxResults      int    `json:"max_results,omitempty" jsonschema:"users per page, 1..1000" validate:"omitempty,gte=1,lte=1000"`
	PaginationToken string `json:"pagination_token,omitempty" jsonschema:"next_token from the previous call; omit for the first page"`
	TimeoutMs       int64  `json:"timeout_ms,omitempty"`
}

type SelfListRequest struct {
	MaxResults      int    `json:"max_results,omitempty" jsonschema:"users per page, 1..1000" validate:"omitempty,gte=1,lte=1000"`
	PaginationToken string `json:"pagination_token,omitempty"`
	TimeoutMs       int64  `json:"timeout_ms,omitempty"`
}

type UsersResponse struct {
	Users     []User `json:"users"`
	Count     int    `json:"count"`
	NextToken string `json:"next_token,omitempty"`
}

type TargetRequest struct {
	Username  string `json:"username" jsonschema:"X username without the @" validate:"required"`
	TimeoutMs int64  `json:"timeout_ms,omitempty"`
}

type ActionResponse struct {
	Username string `json:"username"`
	UserID   string `json:"user_id"`
	Status   string `json:"status"`
}
