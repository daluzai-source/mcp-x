package users

import (
	"context"
	"log/slog"
	"time"

	"github.com/Role1776/mcp-x/app/internal/domain/x/users"
)

const (
	statusDone = "done"
)

type usersAdapter interface {
	ByUsernames(ctx context.Context, names users.Usernames) (users.Users, error)
	ByIDs(ctx context.Context, ids users.UserIDs) (users.Users, error)
	Me(ctx context.Context) (users.User, error)
	Follow(ctx context.Context, target users.UserID) error
	Unfollow(ctx context.Context, target users.UserID) error
	Followers(ctx context.Context, id users.UserID, maxResults int, token string) (users.Users, string, error)
	Following(ctx context.Context, id users.UserID, maxResults int, token string) (users.Users, string, error)
	Blocked(ctx context.Context, maxResults int, token string) (users.Users, string, error)
	Mute(ctx context.Context, target users.UserID) error
	Unmute(ctx context.Context, target users.UserID) error
	Muted(ctx context.Context, maxResults int, token string) (users.Users, string, error)
}

type Config struct {
	MaxUsernames     int   `env:"USERS_MAX_USERNAMES" envDefault:"100" validate:"gt=0,lte=100"`
	MaxIDs           int   `env:"USERS_MAX_IDS" envDefault:"100" validate:"gt=0,lte=100"`
	DefaultResults   int   `env:"USERS_DEFAULT_RESULTS" envDefault:"100" validate:"gt=0"`
	MaxResults       int   `env:"USERS_MAX_RESULTS" envDefault:"1000" validate:"gt=0"`
	DefaultTimeoutMs int64 `env:"USERS_DEFAULT_TIMEOUT_MS" envDefault:"15000" validate:"gt=0"`
	MinTimeoutMs     int64 `env:"USERS_MIN_TIMEOUT_MS" envDefault:"2000" validate:"gt=0"`
	MaxTimeoutMs     int64 `env:"USERS_MAX_TIMEOUT_MS" envDefault:"60000" validate:"gt=0"`
}

type UseCase struct {
	users  usersAdapter
	logger *slog.Logger
	cfg    *Config
}

func New(usersAdapter usersAdapter, logger *slog.Logger, cfg *Config) *UseCase {
	return &UseCase{users: usersAdapter, logger: logger, cfg: cfg}
}

func (u *UseCase) withTimeout(ctx context.Context, requested int64) (context.Context, context.CancelFunc) {
	ms := resolveMinMax(requested, u.cfg.DefaultTimeoutMs, u.cfg.MinTimeoutMs, u.cfg.MaxTimeoutMs)

	return context.WithTimeout(ctx, time.Duration(ms)*time.Millisecond)
}
