package media

import (
	"context"
	"log/slog"
	"time"

	"github.com/Role1776/mcp-x/app/internal/domain/x/media"
)

const (
	statusReady   = "ready"
	statusPending = "pending"
)

type mediaAdapter interface {
	Initialize(ctx context.Context, upload media.Upload) (media.MediaID, error)
	Append(ctx context.Context, id media.MediaID, index int, chunk []byte) error
	Finalize(ctx context.Context, id media.MediaID) (media.Media, error)
	Status(ctx context.Context, id media.MediaID) (media.Media, error)
}

type Config struct {
	Root             string `env:"MEDIA_ROOT" validate:"required"`
	ChunkBytes       int    `env:"MEDIA_CHUNK_BYTES" envDefault:"4194304" validate:"gt=0"`
	PollIntervalMs   int64  `env:"MEDIA_POLL_INTERVAL_MS" envDefault:"2000" validate:"gt=0"`
	DefaultTimeoutMs int64  `env:"MEDIA_DEFAULT_TIMEOUT_MS" envDefault:"120000" validate:"gt=0"`
	MinTimeoutMs     int64  `env:"MEDIA_MIN_TIMEOUT_MS" envDefault:"5000" validate:"gt=0"`
	MaxTimeoutMs     int64  `env:"MEDIA_MAX_TIMEOUT_MS" envDefault:"600000" validate:"gt=0"`
}

type UseCase struct {
	media  mediaAdapter
	logger *slog.Logger
	cfg    *Config
}

func New(mediaAdapter mediaAdapter, logger *slog.Logger, cfg *Config) *UseCase {
	return &UseCase{media: mediaAdapter, logger: logger, cfg: cfg}
}

func (u *UseCase) withTimeout(ctx context.Context, requested int64) (context.Context, context.CancelFunc) {
	ms := requested
	if ms <= 0 {
		ms = u.cfg.DefaultTimeoutMs
	}
	ms = min(max(ms, u.cfg.MinTimeoutMs), u.cfg.MaxTimeoutMs)

	return context.WithTimeout(ctx, time.Duration(ms)*time.Millisecond)
}
