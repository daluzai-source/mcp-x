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

func (u *UseCase) Members(ctx context.Context, req dto.ListPageRequest) (dto.UsersResponse, error) {
	const op = "usecase.lists.Members"

	if err := validator.Validate(req); err != nil {
		return dto.UsersResponse{}, fmt.Errorf("%s: %w", op, x.ErrInvalidRequest)
	}

	id, err := lists.NewListID(req.ID)
	if err != nil {
		return dto.UsersResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	ctx, cancel := u.withTimeout(ctx, req.TimeoutMs)
	defer cancel()

	limit := resolveMax(req.MaxResults, u.cfg.DefaultResults, u.cfg.MaxResults)

	found, nextToken, err := u.lists.Members(ctx, id, limit, req.PaginationToken)
	if err != nil {
		return dto.UsersResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	limited := found.Dedupe().Limit(limit)

	return dto.UsersResponse{Users: toUsers(limited), Count: limited.Len(), NextToken: nextToken}, nil
}

func (u *UseCase) MemberAdd(ctx context.Context, req dto.MemberRequest) (dto.ActionResponse, error) {
	const op = "usecase.lists.MemberAdd"

	return u.member(ctx, op, req, u.lists.MemberAdd)
}

func (u *UseCase) MemberRemove(ctx context.Context, req dto.MemberRequest) (dto.ActionResponse, error) {
	const op = "usecase.lists.MemberRemove"

	return u.member(ctx, op, req, u.lists.MemberRemove)
}

func (u *UseCase) member(
	ctx context.Context,
	op string,
	req dto.MemberRequest,
	action func(context.Context, lists.ListID, string) error,
) (dto.ActionResponse, error) {
	if err := validator.Validate(req); err != nil {
		return dto.ActionResponse{}, fmt.Errorf("%s: %w", op, x.ErrInvalidRequest)
	}

	id, err := lists.NewListID(req.ID)
	if err != nil {
		return dto.ActionResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	ctx, cancel := u.withTimeout(ctx, req.TimeoutMs)
	defer cancel()

	ids, err := u.users.ResolveIDs(ctx, []string{req.Username})
	if err != nil {
		return dto.ActionResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	userID, ok := ids[strings.ToLower(strings.TrimPrefix(strings.TrimSpace(req.Username), "@"))]
	if !ok {
		return dto.ActionResponse{}, fmt.Errorf("%s: %w", op, x.ErrUserNotFound)
	}

	if err := action(ctx, id, userID); err != nil {
		return dto.ActionResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	return dto.ActionResponse{ID: id.String(), Status: statusDone}, nil
}
