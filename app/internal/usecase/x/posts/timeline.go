package posts

import (
	"context"
	"fmt"

	"github.com/Role1776/mcp-x/app/internal/domain/x"
	"github.com/Role1776/mcp-x/app/internal/domain/x/posts"
	dto "github.com/Role1776/mcp-x/app/internal/dto/x/posts"
	"github.com/Role1776/mcp-x/app/internal/pkg/validator"
)

func (u *UseCase) Mentions(ctx context.Context, req dto.TimelineRequest) (dto.PostsResponse, error) {
	const op = "usecase.posts.Mentions"

	return u.timeline(ctx, op, req, u.posts.Mentions)
}

func (u *UseCase) Home(ctx context.Context, req dto.TimelineRequest) (dto.PostsResponse, error) {
	const op = "usecase.posts.Home"

	return u.timeline(ctx, op, req, u.posts.Home)
}

func (u *UseCase) timeline(
	ctx context.Context,
	op string,
	req dto.TimelineRequest,
	fetch func(context.Context, int, string) (posts.Posts, string, error),
) (dto.PostsResponse, error) {
	if err := validator.Validate(req); err != nil {
		return dto.PostsResponse{}, fmt.Errorf("%s: %w", op, x.ErrInvalidRequest)
	}

	ctx, cancel := u.withTimeout(ctx, req.TimeoutMs)
	defer cancel()

	limit := resolveMax(req.MaxResults, u.cfg.DefaultResults, u.cfg.MaxResults)

	found, nextToken, err := fetch(ctx, apiResults(limit), req.PaginationToken)
	if err != nil {
		return dto.PostsResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	limited := found.Dedupe().Limit(limit)

	return dto.PostsResponse{Posts: toPosts(limited), Count: limited.Len(), NextToken: nextToken}, nil
}

func (u *UseCase) Quotes(ctx context.Context, req dto.QuotesRequest) (dto.PostsResponse, error) {
	const op = "usecase.posts.Quotes"

	if err := validator.Validate(req); err != nil {
		return dto.PostsResponse{}, fmt.Errorf("%s: %w", op, x.ErrInvalidRequest)
	}

	id, err := posts.NewPostID(req.ID)
	if err != nil {
		return dto.PostsResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	ctx, cancel := u.withTimeout(ctx, req.TimeoutMs)
	defer cancel()

	limit := resolveMax(req.MaxResults, u.cfg.DefaultResults, u.cfg.MaxResults)

	found, nextToken, err := u.posts.Quotes(ctx, id, apiResults(limit), req.PaginationToken)
	if err != nil {
		return dto.PostsResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	limited := found.Dedupe().Limit(limit)

	return dto.PostsResponse{Posts: toPosts(limited), Count: limited.Len(), NextToken: nextToken}, nil
}

func (u *UseCase) Liked(ctx context.Context, req dto.TimelineRequest) (dto.PostsResponse, error) {
	const op = "usecase.posts.Liked"

	return u.timeline(ctx, op, req, u.posts.Liked)
}

func (u *UseCase) Bookmarks(ctx context.Context, req dto.TimelineRequest) (dto.PostsResponse, error) {
	const op = "usecase.posts.Bookmarks"

	return u.timeline(ctx, op, req, u.posts.Bookmarks)
}
