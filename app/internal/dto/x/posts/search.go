package posts

type SearchRequest struct {
	Queries         []string `json:"queries" jsonschema:"X search queries, up to 512 characters each, executed in parallel; operators: from: to: @mention #tag \"exact phrase\" lang: is:retweet is:reply is:quote has:media has:links url: conversation_id:; space means AND, OR and - (NOT) are supported" validate:"required,min=1,dive,required"`
	MaxResults      int      `json:"max_results,omitempty" jsonschema:"posts to return per query, 1..100; defaults to 10" validate:"omitempty,gte=1,lte=100"`
	SortOrder       string   `json:"sort_order,omitempty" jsonschema:"recency for newest first, relevancy for best match" validate:"omitempty,oneof=recency relevancy"`
	Days            int      `json:"days,omitempty" jsonschema:"how far back to search in days, 1..7; recent search covers the last 7 days only" validate:"omitempty,gte=1,lte=7"`
	IncludeRetweets bool     `json:"include_retweets,omitempty" jsonschema:"include retweets; when false the server appends -is:retweet to every query"`
	TimeoutMs       int64    `json:"timeout_ms,omitempty" jsonschema:"timeout for the whole call in milliseconds"`
}

type SearchResponse struct {
	Results  []SearchResult `json:"results"`
	Metadata Metadata       `json:"metadata"`
}

type SearchResult struct {
	Query       string `json:"query"`
	Status      string `json:"status"`
	Count       int    `json:"count"`
	Posts       []Post `json:"posts"`
	NextToken   string `json:"next_token,omitempty"`
	TotalTimeMs int64  `json:"total_time_ms"`
}

type Metadata struct {
	TotalRequestTimeMs int64 `json:"total_request_time_ms"`
}

type CountRequest struct {
	Queries     []string `json:"queries" jsonschema:"X search queries; same operators as x_posts_search" validate:"required,min=1,dive,required"`
	Granularity string   `json:"granularity,omitempty" jsonschema:"bucket size: minute, hour or day; defaults to hour" validate:"omitempty,oneof=minute hour day"`
	Days        int      `json:"days,omitempty" jsonschema:"how far back to count in days, 1..7" validate:"omitempty,gte=1,lte=7"`
	TimeoutMs   int64    `json:"timeout_ms,omitempty"`
}

type CountResponse struct {
	Results  []CountResult `json:"results"`
	Metadata Metadata      `json:"metadata"`
}

type CountResult struct {
	Query   string   `json:"query"`
	Status  string   `json:"status"`
	Total   int      `json:"total"`
	Buckets []Bucket `json:"buckets,omitempty"`
}

type Bucket struct {
	Start string `json:"start"`
	End   string `json:"end"`
	Count int    `json:"count"`
}
