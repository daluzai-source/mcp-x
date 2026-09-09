package posts

import (
	"context"
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

func (u *UseCase) Count(ctx context.Context, req dto.CountRequest) (dto.CountResponse, error) {
	const op = "usecase.posts.Count"

	if err := validator.Validate(req); err != nil {
		return dto.CountResponse{}, fmt.Errorf("%s: %w", op, x.ErrInvalidRequest)
	}

	if len(req.Queries) > u.cfg.MaxQueries {
		return dto.CountResponse{}, fmt.Errorf("%s: %w", op, x.ErrTooManyQueries)
	}

	queries, err := posts.NewQueries(req.Queries)
	if err != nil {
		return dto.CountResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	window, err := posts.NewWindow(resolveMax(req.Days, u.cfg.DefaultDays, u.cfg.MaxDays))
	if err != nil {
		return dto.CountResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	ctx, cancel := u.withTimeout(ctx, req.TimeoutMs)
	defer cancel()

	start := time.Now()

	results, errs := parallel.Map(ctx, queries, len(queries), func(query posts.Query) (dto.CountResult, error) {
		counts, err := u.posts.Count(ctx, adapter.CountParams{
			Query:       query,
			Granularity: req.Granularity,
			StartTime:   window.StartTime(time.Now()),
		})
		if err != nil {
			u.logger.Error("count failed",
				slog.String("op", op), slog.String("query", query.String()), slog.Any("err", err))

			return dto.CountResult{Query: query.String(), Status: statusOf(err)}, err
		}

		buckets := make([]dto.Bucket, 0, len(counts.Buckets))
		for _, bucket := range counts.Buckets {
			buckets = append(buckets, dto.Bucket{
				Start: bucket.Start.Format(time.RFC3339),
				End:   bucket.End.Format(time.RFC3339),
				Count: bucket.Total,
			})
		}

		return dto.CountResult{
			Query:   query.String(),
			Status:  statusSuccess,
			Total:   counts.Total,
			Buckets: buckets,
		}, nil
	})

	if len(errs) >= len(queries) && len(errs) > 0 {
		return dto.CountResponse{}, fmt.Errorf("%s: %w", op, firstFatal(errs, x.ErrAllQueriesFailed))
	}

	return dto.CountResponse{
		Results:  results,
		Metadata: dto.Metadata{TotalRequestTimeMs: time.Since(start).Milliseconds()},
	}, nil
}
