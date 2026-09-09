package lists

import (
	"context"
	"log/slog"
	"time"

	"github.com/Role1776/mcp-x/app/internal/domain/x/lists"
	"github.com/Role1776/mcp-x/app/internal/domain/x/posts"
	"github.com/Role1776/mcp-x/app/internal/domain/x/users"
)

const statusDone = "done"

type listsAdapter interface {
	Create(ctx context.Context, draft lists.Draft) (string, string, error)
	Update(ctx context.Context, id lists.ListID, patch lists.Patch) error
	Delete(ctx context.Context, id lists.ListID) error
	Get(ctx context.Context, id lists.ListID) (lists.List, error)
	Owned(ctx context.Context, userID string, maxResults int, token string) (lists.Lists, string, error)
	Posts(ctx context.Context, id lists.ListID, maxResults int, token string) (posts.Posts, string, error)
	MemberAdd(ctx context.Context, id lists.ListID, userID string) error
	MemberRemove(ctx context.Context, id lists.ListID, userID string) error
	Members(ctx context.Context, id lists.ListID, maxResults int, token string) (users.Users, string, error)
	Follow(ctx context.Context, id lists.ListID) error
	Unfollow(ctx context.Context, id lists.ListID) error
	Pin(ctx context.Context, id lists.ListID) error
	Unpin(ctx context.Context, id lists.ListID) error
	Pinned(ctx context.Context) (lists.Lists, error)
}

// userResolver speaks plain strings so that this zone never imports another
// zone's domain package.
type userResolver interface {
	ResolveIDs(ctx context.Context, usernames []string) (map[string]string, error)
}

type Config struct {
	DefaultResults   int   `env:"LISTS_DEFAULT_RESULTS" envDefault:"25" validate:"gt=0"`
	MaxResults       int   `env:"LISTS_MAX_RESULTS" envDefault:"100" validate:"gt=0,lte=100"`
	DefaultTimeoutMs int64 `env:"LISTS_DEFAULT_TIMEOUT_MS" envDefault:"15000" validate:"gt=0"`
	MinTimeoutMs     int64 `env:"LISTS_MIN_TIMEOUT_MS" envDefault:"2000" validate:"gt=0"`
	MaxTimeoutMs     int64 `env:"LISTS_MAX_TIMEOUT_MS" envDefault:"60000" validate:"gt=0"`
}

type UseCase struct {
	lists  listsAdapter
	users  userResolver
	logger *slog.Logger
	cfg    *Config
}

func New(listsAdapter listsAdapter, users userResolver, logger *slog.Logger, cfg *Config) *UseCase {
	return &UseCase{lists: listsAdapter, users: users, logger: logger, cfg: cfg}
}

func (u *UseCase) withTimeout(ctx context.Context, requested int64) (context.Context, context.CancelFunc) {
	ms := resolveMinMax(requested, u.cfg.DefaultTimeoutMs, u.cfg.MinTimeoutMs, u.cfg.MaxTimeoutMs)

	return context.WithTimeout(ctx, time.Duration(ms)*time.Millisecond)
}
