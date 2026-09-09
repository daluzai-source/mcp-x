package posts

import (
	"context"
	"fmt"

	"github.com/Role1776/mcp-x/app/internal/domain/x"
	"github.com/Role1776/mcp-x/app/internal/domain/x/posts"
	dto "github.com/Role1776/mcp-x/app/internal/dto/x/posts"
	"github.com/Role1776/mcp-x/app/internal/pkg/validator"
)

func (u *UseCase) Lookup(ctx context.Context, req dto.LookupRequest) (dto.LookupResponse, error) {
	const op = "usecase.posts.Lookup"

	if err := validator.Validate(req); err != nil {
		return dto.LookupResponse{}, fmt.Errorf("%s: %w", op, x.ErrInvalidRequest)
	}

	if len(req.IDs) > u.cfg.MaxIDs {
		return dto.LookupResponse{}, fmt.Errorf("%s: %w", op, x.ErrTooManyIDs)
	}

	ids, err := posts.NewPostIDs(req.IDs)
	if err != nil {
		return dto.LookupResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	ctx, cancel := u.withTimeout(ctx, 0)
	defer cancel()

	found, failures, err := u.posts.Lookup(ctx, ids)
	if err != nil {
		return dto.LookupResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	byID := make(map[string]dto.Post, found.Len())
	for _, view := range found.Views() {
		byID[view.ID] = toPost(view)
	}

	details := make(map[string]string, len(failures))
	for _, failure := range failures {
		details[failure.ResourceID] = failure.Detail
	}

	results := make([]dto.LookupResult, 0, ids.Len())
	for _, id := range ids.Strings() {
		post, ok := byID[id]
		if !ok {
			results = append(results, dto.LookupResult{ID: id, Status: statusNotFound, Detail: details[id]})

			continue
		}

		results = append(results, dto.LookupResult{ID: id, Status: statusSuccess, Post: &post})
	}

	return dto.LookupResponse{Results: results, Count: found.Len()}, nil
}
