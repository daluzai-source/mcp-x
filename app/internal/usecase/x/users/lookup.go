package users

import (
	"context"
	"fmt"

	"github.com/Role1776/mcp-x/app/internal/domain/x"
	"github.com/Role1776/mcp-x/app/internal/domain/x/users"
	dto "github.com/Role1776/mcp-x/app/internal/dto/x/users"
	"github.com/Role1776/mcp-x/app/internal/pkg/validator"
)

func (u *UseCase) Lookup(ctx context.Context, req dto.LookupRequest) (dto.LookupResponse, error) {
	const op = "usecase.users.Lookup"

	if err := validator.Validate(req); err != nil {
		return dto.LookupResponse{}, fmt.Errorf("%s: %w", op, x.ErrInvalidRequest)
	}

	if (len(req.Usernames) == 0) == (len(req.IDs) == 0) {
		return dto.LookupResponse{}, fmt.Errorf("%s: %w", op, x.ErrInvalidRequest)
	}

	if len(req.Usernames) > u.cfg.MaxUsernames || len(req.IDs) > u.cfg.MaxIDs {
		return dto.LookupResponse{}, fmt.Errorf("%s: %w", op, x.ErrTooManyUsers)
	}

	ctx, cancel := u.withTimeout(ctx, req.TimeoutMs)
	defer cancel()

	if len(req.Usernames) > 0 {
		return u.lookupByUsernames(ctx, op, req.Usernames)
	}

	return u.lookupByIDs(ctx, op, req.IDs)
}

func (u *UseCase) lookupByUsernames(ctx context.Context, op string, raw []string) (dto.LookupResponse, error) {
	names, err := users.NewUsernames(raw)
	if err != nil {
		return dto.LookupResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	found, err := u.users.ByUsernames(ctx, names)
	if err != nil {
		return dto.LookupResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	byUsername := found.IDsByUsername()

	notFound := make([]string, 0)
	for _, name := range names {
		if _, ok := byUsername[name.Normalized()]; !ok {
			notFound = append(notFound, name.String())
		}
	}

	return dto.LookupResponse{Users: toUsers(found), Count: found.Len(), NotFound: notFound}, nil
}

func (u *UseCase) lookupByIDs(ctx context.Context, op string, raw []string) (dto.LookupResponse, error) {
	ids, err := users.NewUserIDs(raw)
	if err != nil {
		return dto.LookupResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	found, err := u.users.ByIDs(ctx, ids)
	if err != nil {
		return dto.LookupResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	present := make(map[string]struct{}, found.Len())
	for _, view := range found.Views() {
		present[view.ID] = struct{}{}
	}

	notFound := make([]string, 0)
	for _, id := range ids.Strings() {
		if _, ok := present[id]; !ok {
			notFound = append(notFound, id)
		}
	}

	return dto.LookupResponse{Users: toUsers(found), Count: found.Len(), NotFound: notFound}, nil
}

func (u *UseCase) Me(ctx context.Context, req dto.MeRequest) (dto.MeResponse, error) {
	const op = "usecase.users.Me"

	ctx, cancel := u.withTimeout(ctx, req.TimeoutMs)
	defer cancel()

	me, err := u.users.Me(ctx)
	if err != nil {
		return dto.MeResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	return dto.MeResponse{User: toUser(me.View())}, nil
}

func (u *UseCase) resolve(ctx context.Context, raw string) (users.User, error) {
	name, err := users.NewUsername(raw)
	if err != nil {
		return users.User{}, err
	}

	found, err := u.users.ByUsernames(ctx, users.Usernames{name})
	if err != nil {
		return users.User{}, err
	}

	user, ok := found.ByUsername(name)
	if !ok {
		return users.User{}, x.ErrUserNotFound
	}

	return user, nil
}
