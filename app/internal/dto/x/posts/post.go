package posts

type Metrics struct {
	Retweets int `json:"retweets"`
	Replies  int `json:"replies"`
	Likes    int `json:"likes"`
	Quotes   int `json:"quotes"`
}

type Reference struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

type Post struct {
	ID             string      `json:"id"`
	Text           string      `json:"text"`
	AuthorID       string      `json:"author_id,omitempty"`
	AuthorUsername string      `json:"author_username,omitempty"`
	CreatedAt      string      `json:"created_at,omitempty"`
	Lang           string      `json:"lang,omitempty"`
	ConversationID string      `json:"conversation_id,omitempty"`
	Metrics        Metrics     `json:"metrics"`
	URLs           []string    `json:"urls,omitempty"`
	MediaKeys      []string    `json:"media_keys,omitempty"`
	Referenced     []Reference `json:"referenced,omitempty"`
	URL            string      `json:"url"`
}

type User struct {
	ID             string `json:"id"`
	Username       string `json:"username"`
	Name           string `json:"name,omitempty"`
	Description    string `json:"description,omitempty"`
	Verified       bool   `json:"verified"`
	Protected      bool   `json:"protected"`
	FollowersCount int    `json:"followers_count"`
	URL            string `json:"url"`
}
