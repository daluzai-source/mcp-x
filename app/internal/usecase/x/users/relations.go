package users

import (
	"context"
	"fmt"

	"github.com/Role1776/mcp-x/app/internal/domain/x"
	"github.com/Role1776/mcp-x/app/internal/domain/x/users"
	dto "github.com/Role1776/mcp-x/app/internal/dto/x/users"
	"github.com/Role1776/mcp-x/app/internal/pkg/validator"
)

func (u *UseCase) Followers(ctx context.Context, req dto.ListRequest) (dto.UsersResponse, error) {
	const op = "usecase.users.Followers"

	return u.listByUsername(ctx, op, req, u.users.Followers)
}

func (u *UseCase) Following(ctx context.Context, req dto.ListRequest) (dto.UsersResponse, error) {
	const op = "usecase.users.Following"

	return u.listByUsername(ctx, op, req, u.users.Following)
}

func (u *UseCase) listByUsername(
	ctx context.Context,
	op string,
	req dto.ListRequest,
	fetch func(context.Context, users.UserID, int, string) (users.Users, string, error),
) (dto.UsersResponse, error) {
	if err := validator.Validate(req); err != nil {
		return dto.UsersResponse{}, fmt.Errorf("%s: %w", op, x.ErrInvalidRequest)
	}

	ctx, cancel := u.withTimeout(ctx, req.TimeoutMs)
	defer cancel()

	user, err := u.resolve(ctx, req.Username)
	if err != nil {
		return dto.UsersResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	limit := resolveMax(req.MaxResults, u.cfg.DefaultResults, u.cfg.MaxResults)

	found, nextToken, err := fetch(ctx, user.ID(), limit, req.PaginationToken)
	if err != nil {
		return dto.UsersResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	limited := found.Dedupe().Limit(limit)

	return dto.UsersResponse{Users: toUsers(limited), Count: limited.Len(), NextToken: nextToken}, nil
}

func (u *UseCase) Follow(ctx context.Context, req dto.TargetRequest) (dto.ActionResponse, error) {
	const op = "usecase.users.Follow"

	return u.act(ctx, op, req, u.users.Follow)
}

func (u *UseCase) Unfollow(ctx context.Context, req dto.TargetRequest) (dto.ActionResponse, error) {
	const op = "usecase.users.Unfollow"

	return u.act(ctx, op, req, u.users.Unfollow)
}

func (u *UseCase) Mute(ctx context.Context, req dto.TargetRequest) (dto.ActionResponse, error) {
	const op = "usecase.users.Mute"

	return u.act(ctx, op, req, u.users.Mute)
}

func (u *UseCase) Unmute(ctx context.Context, req dto.TargetRequest) (dto.ActionResponse, error) {
	const op = "usecase.users.Unmute"

	return u.act(ctx, op, req, u.users.Unmute)
}

func (u *UseCase) act(
	ctx context.Context,
	op string,
	req dto.TargetRequest,
	action func(context.Context, users.UserID) error,
) (dto.ActionResponse, error) {
	if err := validator.Validate(req); err != nil {
		return dto.ActionResponse{}, fmt.Errorf("%s: %w", op, x.ErrInvalidRequest)
	}

	ctx, cancel := u.withTimeout(ctx, req.TimeoutMs)
	defer cancel()

	user, err := u.resolve(ctx, req.Username)
	if err != nil {
		return dto.ActionResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	if err := action(ctx, user.ID()); err != nil {
		return dto.ActionResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	return dto.ActionResponse{
		Username: user.Username().String(),
		UserID:   user.ID().String(),
		Status:   statusDone,
	}, nil
}

func (u *UseCase) Blocked(ctx context.Context, req dto.SelfListRequest) (dto.UsersResponse, error) {
	const op = "usecase.users.Blocked"

	return u.selfList(ctx, op, req, u.users.Blocked)
}

func (u *UseCase) Muted(ctx context.Context, req dto.SelfListRequest) (dto.UsersResponse, error) {
	const op = "usecase.users.Muted"

	return u.selfList(ctx, op, req, u.users.Muted)
}

func (u *UseCase) selfList(
	ctx context.Context,
	op string,
	req dto.SelfListRequest,
	fetch func(context.Context, int, string) (users.Users, string, error),
) (dto.UsersResponse, error) {
	if err := validator.Validate(req); err != nil {
		return dto.UsersResponse{}, fmt.Errorf("%s: %w", op, x.ErrInvalidRequest)
	}

	ctx, cancel := u.withTimeout(ctx, req.TimeoutMs)
	defer cancel()

	limit := resolveMax(req.MaxResults, u.cfg.DefaultResults, u.cfg.MaxResults)

	found, nextToken, err := fetch(ctx, limit, req.PaginationToken)
	if err != nil {
		return dto.UsersResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	limited := found.Dedupe().Limit(limit)

	return dto.UsersResponse{Users: toUsers(limited), Count: limited.Len(), NextToken: nextToken}, nil
}
