package posts

type LookupRequest struct {
	IDs []string `json:"ids" jsonschema:"post ids to fetch, up to 100 in one call; batch them instead of calling once per id" validate:"required,min=1,max=100,dive,required"`
}

type LookupResponse struct {
	Results []LookupResult `json:"results"`
	Count   int            `json:"count"`
}

type LookupResult struct {
	ID     string `json:"id"`
	Status string `json:"status"`
	Post   *Post  `json:"post,omitempty"`
	Detail string `json:"detail,omitempty"`
}

type ByUserRequest struct {
	Usernames  []string `json:"usernames" jsonschema:"X usernames without the @, fetched in parallel" validate:"required,min=1,dive,required"`
	MaxResults int      `json:"max_results,omitempty" jsonschema:"posts per user, 1..100" validate:"omitempty,gte=1,lte=100"`
	Exclude    []string `json:"exclude,omitempty" jsonschema:"drop replies and/or retweets from the timeline" validate:"omitempty,dive,oneof=replies retweets"`
	Days       int      `json:"days,omitempty" jsonschema:"only posts newer than this many days, 1..7" validate:"omitempty,gte=1,lte=7"`
	TimeoutMs  int64    `json:"timeout_ms,omitempty"`
}

type ByUserResponse struct {
	Results  []ByUserResult `json:"results"`
	Metadata Metadata       `json:"metadata"`
}

type ByUserResult struct {
	Username  string `json:"username"`
	Status    string `json:"status"`
	Count     int    `json:"count"`
	Posts     []Post `json:"posts"`
	NextToken string `json:"next_token,omitempty"`
}

type TimelineRequest struct {
	MaxResults      int    `json:"max_results,omitempty" jsonschema:"posts to return, 1..100" validate:"omitempty,gte=1,lte=100"`
	PaginationToken string `json:"pagination_token,omitempty" jsonschema:"next_token from the previous call; omit for the first page"`
	TimeoutMs       int64  `json:"timeout_ms,omitempty"`
}

type QuotesRequest struct {
	ID              string `json:"id" jsonschema:"id of the post whose quotes to list" validate:"required"`
	MaxResults      int    `json:"max_results,omitempty" validate:"omitempty,gte=1,lte=100"`
	PaginationToken string `json:"pagination_token,omitempty"`
	TimeoutMs       int64  `json:"timeout_ms,omitempty"`
}

type PostsResponse struct {
	Posts     []Post `json:"posts"`
	Count     int    `json:"count"`
	NextToken string `json:"next_token,omitempty"`
}

type UsersResponse struct {
	Users     []User `json:"users"`
	Count     int    `json:"count"`
	NextToken string `json:"next_token,omitempty"`
}

type PostUsersRequest struct {
	ID              string `json:"id" jsonschema:"post id" validate:"required"`
	MaxResults      int    `json:"max_results,omitempty" validate:"omitempty,gte=1,lte=100"`
	PaginationToken string `json:"pagination_token,omitempty"`
	TimeoutMs       int64  `json:"timeout_ms,omitempty"`
}
