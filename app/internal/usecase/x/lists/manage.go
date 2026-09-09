package lists

import (
	"context"
	"fmt"

	"github.com/Role1776/mcp-x/app/internal/domain/x"
	"github.com/Role1776/mcp-x/app/internal/domain/x/lists"
	dto "github.com/Role1776/mcp-x/app/internal/dto/x/lists"
	"github.com/Role1776/mcp-x/app/internal/pkg/validator"
)

func (u *UseCase) Create(ctx context.Context, req dto.CreateRequest) (dto.ActionResponse, error) {
	const op = "usecase.lists.Create"

	if err := validator.Validate(req); err != nil {
		return dto.ActionResponse{}, fmt.Errorf("%s: %w", op, x.ErrInvalidRequest)
	}

	draft, err := lists.NewDraft(lists.DraftProps{
		Name:        req.Name,
		Description: req.Description,
		Private:     req.Private,
	})
	if err != nil {
		return dto.ActionResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	ctx, cancel := u.withTimeout(ctx, req.TimeoutMs)
	defer cancel()

	id, _, err := u.lists.Create(ctx, draft)
	if err != nil {
		return dto.ActionResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	return dto.ActionResponse{ID: id, Status: statusDone}, nil
}

func (u *UseCase) Update(ctx context.Context, req dto.UpdateRequest) (dto.ActionResponse, error) {
	const op = "usecase.lists.Update"

	if err := validator.Validate(req); err != nil {
		return dto.ActionResponse{}, fmt.Errorf("%s: %w", op, x.ErrInvalidRequest)
	}

	id, err := lists.NewListID(req.ID)
	if err != nil {
		return dto.ActionResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	patch, err := lists.NewPatch(lists.PatchProps{
		Name:        req.Name,
		Description: req.Description,
		Private:     req.Private,
	})

	if err != nil {
		return dto.ActionResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	ctx, cancel := u.withTimeout(ctx, req.TimeoutMs)
	defer cancel()

	if err := u.lists.Update(ctx, id, patch); err != nil {
		return dto.ActionResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	return dto.ActionResponse{ID: id.String(), Status: statusDone}, nil
}

func (u *UseCase) Delete(ctx context.Context, req dto.ListIDRequest) (dto.ActionResponse, error) {
	const op = "usecase.lists.Delete"

	return u.act(ctx, op, req, u.lists.Delete)
}

func (u *UseCase) Follow(ctx context.Context, req dto.ListIDRequest) (dto.ActionResponse, error) {
	const op = "usecase.lists.Follow"

	return u.act(ctx, op, req, u.lists.Follow)
}

func (u *UseCase) Unfollow(ctx context.Context, req dto.ListIDRequest) (dto.ActionResponse, error) {
	const op = "usecase.lists.Unfollow"

	return u.act(ctx, op, req, u.lists.Unfollow)
}

func (u *UseCase) Pin(ctx context.Context, req dto.ListIDRequest) (dto.ActionResponse, error) {
	const op = "usecase.lists.Pin"

	return u.act(ctx, op, req, u.lists.Pin)
}

func (u *UseCase) Unpin(ctx context.Context, req dto.ListIDRequest) (dto.ActionResponse, error) {
	const op = "usecase.lists.Unpin"

	return u.act(ctx, op, req, u.lists.Unpin)
}

func (u *UseCase) act(
	ctx context.Context,
	op string,
	req dto.ListIDRequest,
	action func(context.Context, lists.ListID) error,
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

	if err := action(ctx, id); err != nil {
		return dto.ActionResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	return dto.ActionResponse{ID: id.String(), Status: statusDone}, nil
}
