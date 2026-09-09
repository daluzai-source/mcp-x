package posts

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	adapter "github.com/Role1776/mcp-x/app/internal/adapter/x/posts"
	"github.com/Role1776/mcp-x/app/internal/domain/x"
	"github.com/Role1776/mcp-x/app/internal/domain/x/posts"
	dto "github.com/Role1776/mcp-x/app/internal/dto/x/posts"
	"github.com/Role1776/mcp-x/app/internal/pkg/parallel"
	"github.com/Role1776/mcp-x/app/internal/pkg/validator"
)

func (u *UseCase) Search(ctx context.Context, req dto.SearchRequest) (dto.SearchResponse, error) {
	const op = "usecase.posts.Search"

	if err := validator.Validate(req); err != nil {
		return dto.SearchResponse{}, fmt.Errorf("%s: %w", op, x.ErrInvalidRequest)
	}

	if len(req.Queries) > u.cfg.MaxQueries {
		return dto.SearchResponse{}, fmt.Errorf("%s: %w", op, x.ErrTooManyQueries)
	}

	queries, err := posts.NewQueries(req.Queries)
	if err != nil {
		return dto.SearchResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	if !req.IncludeRetweets {
		queries, err = queries.WithoutRetweets()
		if err != nil {
			return dto.SearchResponse{}, fmt.Errorf("%s: %w", op, err)
		}
	}

	window, err := posts.NewWindow(resolveMax(req.Days, u.cfg.DefaultDays, u.cfg.MaxDays))
	if err != nil {
		return dto.SearchResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	ctx, cancel := u.withTimeout(ctx, req.TimeoutMs)
	defer cancel()

	start := time.Now()
	limit := resolveMax(req.MaxResults, u.cfg.DefaultResults, u.cfg.MaxResults)

	results, errs := parallel.Map(ctx, queries, len(queries), func(query posts.Query) (dto.SearchResult, error) {
		return u.searchOne(ctx, query, window, limit, req.SortOrder)
	})

	if len(errs) >= len(queries) && len(errs) > 0 {
		u.logger.Error("all queries failed", slog.String("op", op), slog.Any("errs", errs))

		return dto.SearchResponse{}, fmt.Errorf("%s: %w", op, firstFatal(errs, x.ErrAllQueriesFailed))
	}

	return dto.SearchResponse{
		Results:  results,
		Metadata: dto.Metadata{TotalRequestTimeMs: time.Since(start).Milliseconds()},
	}, nil
}

func (u *UseCase) searchOne(ctx context.Context, query posts.Query, window posts.Window, limit int, sortOrder string) (dto.SearchResult, error) {
	const op = "usecase.posts.searchOne"

	start := time.Now()

	found, nextToken, err := u.posts.Search(ctx, adapter.SearchParams{
		Query:      query,
		MaxResults: apiResults(limit),
		SortOrder:  sortOrder,
		StartTime:  window.StartTime(time.Now()),
	})
	if err != nil {
		u.logger.Error("search failed",
			slog.String("op", op), slog.String("query", query.String()), slog.Any("err", err))

		return dto.SearchResult{
			Query:       query.String(),
			Status:      statusOf(err),
			Posts:       []dto.Post{},
			TotalTimeMs: time.Since(start).Milliseconds(),
		}, err
	}

	limited := found.Dedupe().Limit(limit)

	return dto.SearchResult{
		Query:       query.String(),
		Status:      statusSuccess,
		Count:       limited.Len(),
		Posts:       toPosts(limited),
		NextToken:   nextToken,
		TotalTimeMs: time.Since(start).Milliseconds(),
	}, nil
}

func statusOf(err error) string {
	if errors.Is(err, context.DeadlineExceeded) {
		return statusTimeout
	}

	return statusFailed
}

// firstFatal surfaces an auth or quota failure instead of the generic "all
// failed": those need a different reaction from the model than a retry.
func firstFatal(errs []error, fallback error) error {
	fatal := []error{
		x.ErrUnauthorized, x.ErrReadOnlyToken, x.ErrUsageCapped,
		x.ErrCreditsDepleted, x.ErrRateLimited, x.ErrForbidden,
	}

	for _, candidate := range fatal {
		for _, err := range errs {
			if errors.Is(err, candidate) {
				return candidate
			}
		}
	}

	return fallback
}
