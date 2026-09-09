package posts

import (
	"slices"
	"strconv"
	"strings"
	"time"
)

type Metrics struct {
	Retweets int
	Replies  int
	Likes    int
	Quotes   int
}

type Reference struct {
	Type string
	ID   string
}

type Post struct {
	id             PostID
	text           string
	authorID       AuthorID
	authorUsername AuthorUsername
	createdAt      time.Time
	lang           string
	conversationID string
	metrics        Metrics
	urls           []string
	mediaKeys      []string
	referenced     []Reference
}

type PostProps struct {
	ID             string
	Text           string
	AuthorID       string
	AuthorUsername string
	CreatedAt      time.Time
	Lang           string
	ConversationID string
	Metrics        Metrics
	URLs           []string
	MediaKeys      []string
	Referenced     []Reference
}

func NewPost(props PostProps) (Post, error) {
	id, err := NewPostID(props.ID)
	if err != nil {
		return Post{}, err
	}

	post := Post{
		id:             id,
		text:           props.Text,
		createdAt:      props.CreatedAt,
		lang:           props.Lang,
		conversationID: props.ConversationID,
		metrics:        props.Metrics,
		urls:           slices.Clone(props.URLs),
		mediaKeys:      slices.Clone(props.MediaKeys),
		referenced:     slices.Clone(props.Referenced),
	}

	if props.AuthorID != "" {
		authorID, err := NewAuthorID(props.AuthorID)
		if err != nil {
			return Post{}, err
		}
		post.authorID = authorID
	}

	if props.AuthorUsername != "" {
		authorUsername, err := NewAuthorUsername(props.AuthorUsername)
		if err != nil {
			return Post{}, err
		}
		post.authorUsername = authorUsername
	}

	return post, nil
}

type PostView struct {
	ID             string
	Text           string
	AuthorID       string
	AuthorUsername string
	CreatedAt      time.Time
	Lang           string
	ConversationID string
	Metrics        Metrics
	URLs           []string
	MediaKeys      []string
	Referenced     []Reference
	URL            string
}

func (p Post) View() PostView {
	return PostView{
		ID:             p.id.String(),
		Text:           p.text,
		AuthorID:       p.authorID.String(),
		AuthorUsername: p.authorUsername.String(),
		CreatedAt:      p.createdAt,
		Lang:           p.lang,
		ConversationID: p.conversationID,
		Metrics:        p.metrics,
		URLs:           slices.Clone(p.urls),
		MediaKeys:      slices.Clone(p.mediaKeys),
		Referenced:     slices.Clone(p.referenced),
		URL:            p.link(),
	}
}

func (p Post) ID() PostID {
	return p.id
}

func (p Post) link() string {
	author := p.authorUsername.String()
	if author == "" {
		author = "i"
	}

	return "https://x.com/" + author + "/status/" + p.id.String()
}

func (p Post) Markdown() string {
	var b strings.Builder

	if author := p.authorUsername.String(); author != "" {
		b.WriteString("@")
		b.WriteString(author)
		b.WriteString(" ")
	}

	if !p.createdAt.IsZero() {
		b.WriteString(p.createdAt.Format(time.RFC3339))
	}

	b.WriteString("\n")
	b.WriteString(p.text)
	b.WriteString("\n")
	b.WriteString(p.link())
	b.WriteString(" | likes ")
	b.WriteString(strconv.Itoa(p.metrics.Likes))
	b.WriteString(", reposts ")
	b.WriteString(strconv.Itoa(p.metrics.Retweets))

	return b.String()
}

type Posts []Post

func (p Posts) Len() int {
	return len(p)
}

func (p Posts) IsEmpty() bool {
	return len(p) == 0
}

func (p Posts) Views() []PostView {
	views := make([]PostView, 0, len(p))
	for _, post := range p {
		views = append(views, post.View())
	}

	return views
}

func (p Posts) Limit(n int) Posts {
	if n <= 0 || len(p) == 0 {
		return Posts{}
	}

	result := make(Posts, min(len(p), n))
	copy(result, p)

	return result
}

func (p Posts) Dedupe() Posts {
	seen := make(map[string]struct{}, len(p))
	result := make(Posts, 0, len(p))

	for _, post := range p {
		key := post.id.String()
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		result = append(result, post)
	}

	return result
}

func (p Posts) Markdown() string {
	parts := make([]string, 0, len(p))
	for _, post := range p {
		parts = append(parts, post.Markdown())
	}

	return strings.Join(parts, "\n\n")
}

type PartialFailure struct {
	ResourceID string
	Title      string
	Detail     string
}

type Count struct {
	Start time.Time
	End   time.Time
	Total int
}

type Counts struct {
	Buckets []Count
	Total   int
}
