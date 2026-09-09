package posts

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	adapter "github.com/Role1776/mcp-x/app/internal/adapter/x/posts"
	"github.com/Role1776/mcp-x/app/internal/domain/x"
	"github.com/Role1776/mcp-x/app/internal/domain/x/posts"
	dto "github.com/Role1776/mcp-x/app/internal/dto/x/posts"
	"github.com/Role1776/mcp-x/app/internal/pkg/parallel"
	"github.com/Role1776/mcp-x/app/internal/pkg/validator"
)

type resolvedUser struct {
	username string
	id       string
}

func (u *UseCase) ByUser(ctx context.Context, req dto.ByUserRequest) (dto.ByUserResponse, error) {
	const op = "usecase.posts.ByUser"

	if err := validator.Validate(req); err != nil {
		return dto.ByUserResponse{}, fmt.Errorf("%s: %w", op, x.ErrInvalidRequest)
	}

	if len(req.Usernames) > u.cfg.MaxUsernames {
		return dto.ByUserResponse{}, fmt.Errorf("%s: %w", op, x.ErrTooManyUsers)
	}

	exclude, err := posts.NewExclude(req.Exclude)
	if err != nil {
		return dto.ByUserResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	window, err := posts.NewWindow(resolveMax(req.Days, u.cfg.DefaultDays, u.cfg.MaxDays))
	if err != nil {
		return dto.ByUserResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	ctx, cancel := u.withTimeout(ctx, req.TimeoutMs)
	defer cancel()

	start := time.Now()

	ids, err := u.users.ResolveIDs(ctx, req.Usernames)
	if err != nil {
		return dto.ByUserResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	resolved := make([]resolvedUser, 0, len(req.Usernames))
	missing := make([]dto.ByUserResult, 0, len(req.Usernames))

	for _, raw := range req.Usernames {
		username, err := posts.NewAuthorUsername(raw)
		if err != nil {
			return dto.ByUserResponse{}, fmt.Errorf("%s: %w", op, err)
		}

		id, ok := ids[strings.ToLower(username.String())]
		if !ok {
			missing = append(missing, dto.ByUserResult{
				Username: username.String(), Status: statusNotFound, Posts: []dto.Post{},
			})

			continue
		}

		resolved = append(resolved, resolvedUser{username: username.String(), id: id})
	}

	if len(resolved) == 0 {
		return dto.ByUserResponse{}, fmt.Errorf("%s: %w", op, x.ErrAllUsersFailed)
	}

	limit := resolveMax(req.MaxResults, u.cfg.DefaultResults, u.cfg.MaxResults)

	results, errs := parallel.Map(ctx, resolved, len(resolved), func(user resolvedUser) (dto.ByUserResult, error) {
		found, nextToken, err := u.posts.UserTimeline(ctx, adapter.TimelineParams{
			UserID:     user.id,
			MaxResults: apiResults(limit),
			Exclude:    exclude,
			StartTime:  window.StartTime(time.Now()),
		})
		if err != nil {
			u.logger.Error("timeline failed",
				slog.String("op", op), slog.String("username", user.username), slog.Any("err", err))

			return dto.ByUserResult{Username: user.username, Status: statusOf(err), Posts: []dto.Post{}}, err
		}

		limited := found.Dedupe().Limit(limit)

		return dto.ByUserResult{
			Username:  user.username,
			Status:    statusSuccess,
			Count:     limited.Len(),
			Posts:     toPosts(limited),
			NextToken: nextToken,
		}, nil
	})

	if len(errs) >= len(resolved) && len(errs) > 0 {
		return dto.ByUserResponse{}, fmt.Errorf("%s: %w", op, firstFatal(errs, x.ErrAllUsersFailed))
	}

	return dto.ByUserResponse{
		Results:  append(results, missing...),
		Metadata: dto.Metadata{TotalRequestTimeMs: time.Since(start).Milliseconds()},
	}, nil
}
