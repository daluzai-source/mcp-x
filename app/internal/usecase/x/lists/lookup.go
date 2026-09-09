package lists

import (
	"context"
	"fmt"
	"strings"

	"github.com/Role1776/mcp-x/app/internal/domain/x"
	"github.com/Role1776/mcp-x/app/internal/domain/x/lists"
	dto "github.com/Role1776/mcp-x/app/internal/dto/x/lists"
	"github.com/Role1776/mcp-x/app/internal/pkg/validator"
)

func (u *UseCase) Get(ctx context.Context, req dto.ListIDRequest) (dto.ListResponse, error) {
	const op = "usecase.lists.Get"

	if err := validator.Validate(req); err != nil {
		return dto.ListResponse{}, fmt.Errorf("%s: %w", op, x.ErrInvalidRequest)
	}

	id, err := lists.NewListID(req.ID)
	if err != nil {
		return dto.ListResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	ctx, cancel := u.withTimeout(ctx, req.TimeoutMs)
	defer cancel()

	list, err := u.lists.Get(ctx, id)
	if err != nil {
		return dto.ListResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	return dto.ListResponse{List: toList(list.View())}, nil
}

func (u *UseCase) Owned(ctx context.Context, req dto.OwnedRequest) (dto.ListsResponse, error) {
	const op = "usecase.lists.Owned"

	if err := validator.Validate(req); err != nil {
		return dto.ListsResponse{}, fmt.Errorf("%s: %w", op, x.ErrInvalidRequest)
	}

	ctx, cancel := u.withTimeout(ctx, req.TimeoutMs)
	defer cancel()

	var userID string
	if req.Username != "" {
		ids, err := u.users.ResolveIDs(ctx, []string{req.Username})
		if err != nil {
			return dto.ListsResponse{}, fmt.Errorf("%s: %w", op, err)
		}

		id, ok := ids[strings.ToLower(strings.TrimPrefix(strings.TrimSpace(req.Username), "@"))]
		if !ok {
			return dto.ListsResponse{}, fmt.Errorf("%s: %w", op, x.ErrUserNotFound)
		}
		userID = id
	}

	limit := resolveMax(req.MaxResults, u.cfg.DefaultResults, u.cfg.MaxResults)

	found, nextToken, err := u.lists.Owned(ctx, userID, limit, req.PaginationToken)
	if err != nil {
		return dto.ListsResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	limited := found.Dedupe().Limit(limit)

	return dto.ListsResponse{Lists: toLists(limited), Count: limited.Len(), NextToken: nextToken}, nil
}

func (u *UseCase) Posts(ctx context.Context, req dto.ListPageRequest) (dto.PostsResponse, error) {
	const op = "usecase.lists.Posts"

	if err := validator.Validate(req); err != nil {
		return dto.PostsResponse{}, fmt.Errorf("%s: %w", op, x.ErrInvalidRequest)
	}

	id, err := lists.NewListID(req.ID)
	if err != nil {
		return dto.PostsResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	ctx, cancel := u.withTimeout(ctx, req.TimeoutMs)
	defer cancel()

	limit := resolveMax(req.MaxResults, u.cfg.DefaultResults, u.cfg.MaxResults)

	found, nextToken, err := u.lists.Posts(ctx, id, limit, req.PaginationToken)
	if err != nil {
		return dto.PostsResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	limited := found.Dedupe().Limit(limit)

	return dto.PostsResponse{Posts: toPosts(limited), Count: limited.Len(), NextToken: nextToken}, nil
}

func (u *UseCase) Pinned(ctx context.Context, req dto.SelfPageRequest) (dto.ListsResponse, error) {
	const op = "usecase.lists.Pinned"

	ctx, cancel := u.withTimeout(ctx, req.TimeoutMs)
	defer cancel()

	found, err := u.lists.Pinned(ctx)
	if err != nil {
		return dto.ListsResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	return dto.ListsResponse{Lists: toLists(found), Count: found.Len()}, nil
}
