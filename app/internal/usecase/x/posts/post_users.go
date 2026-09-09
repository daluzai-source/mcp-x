package posts

import (
	"context"
	"fmt"

	"github.com/Role1776/mcp-x/app/internal/domain/x"
	"github.com/Role1776/mcp-x/app/internal/domain/x/posts"
	"github.com/Role1776/mcp-x/app/internal/domain/x/users"
	dto "github.com/Role1776/mcp-x/app/internal/dto/x/posts"
	"github.com/Role1776/mcp-x/app/internal/pkg/validator"
)

func (u *UseCase) postUsers(
	ctx context.Context,
	op string,
	req dto.PostUsersRequest,
	fetch func(context.Context, posts.PostID, int, string) (users.Users, error),
) (dto.UsersResponse, error) {
	if err := validator.Validate(req); err != nil {
		return dto.UsersResponse{}, fmt.Errorf("%s: %w", op, x.ErrInvalidRequest)
	}

	id, err := posts.NewPostID(req.ID)
	if err != nil {
		return dto.UsersResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	ctx, cancel := u.withTimeout(ctx, req.TimeoutMs)
	defer cancel()

	limit := resolveMax(req.MaxResults, u.cfg.DefaultResults, u.cfg.MaxResults)

	found, err := fetch(ctx, id, limit, req.PaginationToken)
	if err != nil {
		return dto.UsersResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	limited := found.Dedupe().Limit(limit)

	return dto.UsersResponse{Users: toUsers(limited), Count: limited.Len()}, nil
}
