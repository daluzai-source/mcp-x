package posts

import (
	"context"
	"fmt"

	"github.com/Role1776/mcp-x/app/internal/domain/x"
	"github.com/Role1776/mcp-x/app/internal/domain/x/posts"
	dto "github.com/Role1776/mcp-x/app/internal/dto/x/posts"
	"github.com/Role1776/mcp-x/app/internal/pkg/validator"
)

func (u *UseCase) Create(ctx context.Context, req dto.CreateRequest) (dto.CreateResponse, error) {
	const op = "usecase.posts.Create"

	if err := validator.Validate(req); err != nil {
		return dto.CreateResponse{}, fmt.Errorf("%s: %w", op, x.ErrInvalidRequest)
	}

	props := posts.DraftProps{
		Text:          req.Text,
		ReplyToID:     req.ReplyToID,
		QuoteID:       req.QuoteID,
		MediaIDs:      req.MediaIDs,
		ReplySettings: req.ReplySettings,
	}

	if req.Poll != nil {
		props.Poll = &posts.PollProps{
			Options:         req.Poll.Options,
			DurationMinutes: req.Poll.DurationMinutes,
		}
	}

	draft, err := posts.NewDraft(props)
	if err != nil {
		return dto.CreateResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	ctx, cancel := u.withTimeout(ctx, 0)
	defer cancel()

	id, text, err := u.posts.Create(ctx, draft)
	if err != nil {
		return dto.CreateResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	return dto.CreateResponse{ID: id, Text: text, URL: "https://x.com/i/status/" + id}, nil
}

func (u *UseCase) Delete(ctx context.Context, req dto.PostIDRequest) (dto.ActionResponse, error) {
	const op = "usecase.posts.Delete"

	return u.act(ctx, op, req, u.posts.Delete)
}

func (u *UseCase) Like(ctx context.Context, req dto.PostIDRequest) (dto.ActionResponse, error) {
	const op = "usecase.posts.Like"

	return u.act(ctx, op, req, u.posts.Like)
}

func (u *UseCase) Unlike(ctx context.Context, req dto.PostIDRequest) (dto.ActionResponse, error) {
	const op = "usecase.posts.Unlike"

	return u.act(ctx, op, req, u.posts.Unlike)
}

func (u *UseCase) Repost(ctx context.Context, req dto.PostIDRequest) (dto.ActionResponse, error) {
	const op = "usecase.posts.Repost"

	return u.act(ctx, op, req, u.posts.Repost)
}

func (u *UseCase) Unrepost(ctx context.Context, req dto.PostIDRequest) (dto.ActionResponse, error) {
	const op = "usecase.posts.Unrepost"

	return u.act(ctx, op, req, u.posts.Unrepost)
}

func (u *UseCase) act(
	ctx context.Context,
	op string,
	req dto.PostIDRequest,
	action func(context.Context, posts.PostID) error,
) (dto.ActionResponse, error) {
	if err := validator.Validate(req); err != nil {
		return dto.ActionResponse{}, fmt.Errorf("%s: %w", op, x.ErrInvalidRequest)
	}

	id, err := posts.NewPostID(req.ID)
	if err != nil {
		return dto.ActionResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	ctx, cancel := u.withTimeout(ctx, 0)
	defer cancel()

	if err := action(ctx, id); err != nil {
		return dto.ActionResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	return dto.ActionResponse{ID: id.String(), Status: statusDone}, nil
}

func (u *UseCase) LikedBy(ctx context.Context, req dto.PostUsersRequest) (dto.UsersResponse, error) {
	const op = "usecase.posts.LikedBy"

	return u.postUsers(ctx, op, req, u.posts.LikedBy)
}

func (u *UseCase) RepostedBy(ctx context.Context, req dto.PostUsersRequest) (dto.UsersResponse, error) {
	const op = "usecase.posts.RepostedBy"

	return u.postUsers(ctx, op, req, u.posts.RepostedBy)
}
