package posts

import (
	"context"
	"log/slog"
	"time"

	adapter "github.com/Role1776/mcp-x/app/internal/adapter/x/posts"
	"github.com/Role1776/mcp-x/app/internal/domain/x/posts"
	"github.com/Role1776/mcp-x/app/internal/domain/x/users"
)

const (
	statusSuccess  = "success"
	statusFailed   = "failed"
	statusTimeout  = "timeout"
	statusNotFound = "not found"
	statusDone     = "done"

	// apiMinResults is the floor the search and timeline endpoints accept; asking
	// for fewer is a 400, so the request is widened and the extra posts are
	// dropped locally by Posts.Limit.
	apiMinResults = 10
)

type postsAdapter interface {
	Search(ctx context.Context, p adapter.SearchParams) (posts.Posts, string, error)
	Count(ctx context.Context, p adapter.CountParams) (posts.Counts, error)
	Lookup(ctx context.Context, ids posts.PostIDs) (posts.Posts, []posts.PartialFailure, error)
	UserTimeline(ctx context.Context, p adapter.TimelineParams) (posts.Posts, string, error)
	Mentions(ctx context.Context, maxResults int, token string) (posts.Posts, string, error)
	Home(ctx context.Context, maxResults int, token string) (posts.Posts, string, error)
	Quotes(ctx context.Context, id posts.PostID, maxResults int, token string) (posts.Posts, string, error)
	Create(ctx context.Context, draft posts.Draft) (string, string, error)
	Delete(ctx context.Context, id posts.PostID) error
	Like(ctx context.Context, id posts.PostID) error
	Unlike(ctx context.Context, id posts.PostID) error
	LikedBy(ctx context.Context, id posts.PostID, maxResults int, token string) (users.Users, error)
	Liked(ctx context.Context, maxResults int, token string) (posts.Posts, string, error)
	Repost(ctx context.Context, id posts.PostID) error
	Unrepost(ctx context.Context, id posts.PostID) error
	RepostedBy(ctx context.Context, id posts.PostID, maxResults int, token string) (users.Users, error)
	Bookmarks(ctx context.Context, maxResults int, token string) (posts.Posts, string, error)
}

// userResolver speaks plain strings so that this zone never imports another
// zone's domain package.
type userResolver interface {
	ResolveIDs(ctx context.Context, usernames []string) (map[string]string, error)
}

type Config struct {
	MaxQueries       int   `env:"POSTS_MAX_QUERIES" envDefault:"5" validate:"gt=0"`
	MaxUsernames     int   `env:"POSTS_MAX_USERNAMES" envDefault:"5" validate:"gt=0"`
	MaxIDs           int   `env:"POSTS_MAX_IDS" envDefault:"100" validate:"gt=0,lte=100"`
	DefaultResults   int   `env:"POSTS_DEFAULT_RESULTS" envDefault:"10" validate:"gt=0"`
	MaxResults       int   `env:"POSTS_MAX_RESULTS" envDefault:"100" validate:"gt=0,lte=100"`
	DefaultDays      int   `env:"POSTS_DEFAULT_DAYS" envDefault:"7" validate:"gt=0,lte=7"`
	MaxDays          int   `env:"POSTS_MAX_DAYS" envDefault:"7" validate:"gt=0,lte=7"`
	DefaultTimeoutMs int64 `env:"POSTS_DEFAULT_TIMEOUT_MS" envDefault:"15000" validate:"gt=0"`
	MinTimeoutMs     int64 `env:"POSTS_MIN_TIMEOUT_MS" envDefault:"2000" validate:"gt=0"`
	MaxTimeoutMs     int64 `env:"POSTS_MAX_TIMEOUT_MS" envDefault:"60000" validate:"gt=0"`
}

type UseCase struct {
	posts  postsAdapter
	users  userResolver
	logger *slog.Logger
	cfg    *Config
}

func New(postsAdapter postsAdapter, users userResolver, logger *slog.Logger, cfg *Config) *UseCase {
	return &UseCase{posts: postsAdapter, users: users, logger: logger, cfg: cfg}
}

func (u *UseCase) withTimeout(ctx context.Context, requested int64) (context.Context, context.CancelFunc) {
	ms := resolveMinMax(requested, u.cfg.DefaultTimeoutMs, u.cfg.MinTimeoutMs, u.cfg.MaxTimeoutMs)

	return context.WithTimeout(ctx, time.Duration(ms)*time.Millisecond)
}
