package users

import (
	"context"
	"fmt"

	"github.com/michimani/gotwi"
	"github.com/michimani/gotwi/user/follow"
	"github.com/michimani/gotwi/user/follow/types"

	"github.com/Role1776/mcp-x/app/internal/adapter/x/apierr"
	"github.com/Role1776/mcp-x/app/internal/domain/x"
	domainusers "github.com/Role1776/mcp-x/app/internal/domain/x/users"
)

func (a *Adapter) Follow(ctx context.Context, target domainusers.UserID) error {
	const op = "adapter.users.Follow"

	if _, err := follow.CreateFollowing(ctx, a.client, &types.CreateFollowingInput{
		ID:       a.selfID,
		TargetID: target.String(),
	}); err != nil {
		return fmt.Errorf("%s: %w", op, apierr.NotFoundAs(err, x.ErrUserNotFound))
	}

	return nil
}

func (a *Adapter) Unfollow(ctx context.Context, target domainusers.UserID) error {
	const op = "adapter.users.Unfollow"

	if _, err := follow.DeleteFollowing(ctx, a.client, &types.DeleteFollowingInput{
		SourceUserID: a.selfID,
		TargetID:     target.String(),
	}); err != nil {
		return fmt.Errorf("%s: %w", op, apierr.NotFoundAs(err, x.ErrUserNotFound))
	}

	return nil
}

func (a *Adapter) Followers(ctx context.Context, id domainusers.UserID, maxResults int, token string) (domainusers.Users, string, error) {
	const op = "adapter.users.Followers"

	out, err := follow.ListFollowers(ctx, a.client, &types.ListFollowersInput{
		ID:              id.String(),
		MaxResults:      types.ListMaxResults(maxResults),
		UserFields:      userFields,
		PaginationToken: token,
	})
	if err != nil {
		return nil, "", fmt.Errorf("%s: %w", op, apierr.NotFoundAs(err, x.ErrUserNotFound))
	}

	return toUsers(out.Data), gotwi.StringValue(out.Meta.NextToken), nil
}

func (a *Adapter) Following(ctx context.Context, id domainusers.UserID, maxResults int, token string) (domainusers.Users, string, error) {
	const op = "adapter.users.Following"

	out, err := follow.ListFollowings(ctx, a.client, &types.ListFollowingsInput{
		ID:              id.String(),
		MaxResults:      types.ListMaxResults(maxResults),
		UserFields:      userFields,
		PaginationToken: token,
	})
	if err != nil {
		return nil, "", fmt.Errorf("%s: %w", op, apierr.NotFoundAs(err, x.ErrUserNotFound))
	}

	return toUsers(out.Data), gotwi.StringValue(out.Meta.NextToken), nil
}
