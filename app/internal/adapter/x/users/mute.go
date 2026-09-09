package users

import (
	"context"
	"fmt"

	"github.com/michimani/gotwi"
	"github.com/michimani/gotwi/user/mute"
	"github.com/michimani/gotwi/user/mute/types"

	"github.com/Role1776/mcp-x/app/internal/adapter/x/apierr"
	"github.com/Role1776/mcp-x/app/internal/domain/x"
	domainusers "github.com/Role1776/mcp-x/app/internal/domain/x/users"
)

func (a *Adapter) Mute(ctx context.Context, target domainusers.UserID) error {
	const op = "adapter.users.Mute"

	if _, err := mute.Create(ctx, a.client, &types.CreateInput{ID: a.selfID, TargetID: target.String()}); err != nil {
		return fmt.Errorf("%s: %w", op, apierr.NotFoundAs(err, x.ErrUserNotFound))
	}

	return nil
}

func (a *Adapter) Unmute(ctx context.Context, target domainusers.UserID) error {
	const op = "adapter.users.Unmute"

	if _, err := mute.Delete(ctx, a.client, &types.DeleteInput{
		SourceUserID: a.selfID,
		TargetID:     target.String(),
	}); err != nil {
		return fmt.Errorf("%s: %w", op, apierr.NotFoundAs(err, x.ErrUserNotFound))
	}

	return nil
}

func (a *Adapter) Muted(ctx context.Context, maxResults int, token string) (domainusers.Users, string, error) {
	const op = "adapter.users.Muted"

	out, err := mute.Lists(ctx, a.client, &types.ListsInput{
		ID:              a.selfID,
		MaxResults:      types.ListMaxResults(maxResults),
		UserFields:      userFields,
		PaginationToken: token,
	})
	if err != nil {
		return nil, "", fmt.Errorf("%s: %w", op, apierr.ToDomain(err))
	}

	return toUsers(out.Data), gotwi.StringValue(out.Meta.NextToken), nil
}
