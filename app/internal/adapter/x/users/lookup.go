package users

import (
	"context"
	"fmt"

	"github.com/michimani/gotwi"

	"github.com/michimani/gotwi/user/userlookup"
	"github.com/michimani/gotwi/user/userlookup/types"

	"github.com/Role1776/mcp-x/app/internal/adapter/x/apierr"
	"github.com/Role1776/mcp-x/app/internal/domain/x"
	domainusers "github.com/Role1776/mcp-x/app/internal/domain/x/users"
)

func (a *Adapter) ByUsernames(ctx context.Context, names domainusers.Usernames) (domainusers.Users, error) {
	const op = "adapter.users.ByUsernames"

	out, err := userlookup.ListByUsernames(ctx, a.client, &types.ListByUsernamesInput{
		Usernames:  names.Strings(),
		UserFields: userFields,
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, apierr.NotFoundAs(err, x.ErrUserNotFound))
	}

	return toUsers(out.Data), nil
}

func (a *Adapter) ByIDs(ctx context.Context, ids domainusers.UserIDs) (domainusers.Users, error) {
	const op = "adapter.users.ByIDs"

	out, err := userlookup.List(ctx, a.client, &types.ListInput{
		IDs:        ids.Strings(),
		UserFields: userFields,
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, apierr.NotFoundAs(err, x.ErrUserNotFound))
	}

	return toUsers(out.Data), nil
}

func (a *Adapter) Me(ctx context.Context) (domainusers.User, error) {
	const op = "adapter.users.Me"

	out, err := userlookup.GetMe(ctx, a.client, &types.GetMeInput{UserFields: userFields})
	if err != nil {
		return domainusers.User{}, fmt.Errorf("%s: %w", op, apierr.ToDomain(err))
	}

	if gotwi.StringValue(out.Data.ID) == "" {
		return domainusers.User{}, fmt.Errorf("%s: %w", op, x.ErrUserNotFound)
	}

	user, err := toUser(out.Data)
	if err != nil {
		return domainusers.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

// ResolveIDs returns plain strings so that a use case in another business zone
// can depend on this method without importing this zone's domain package.
func (a *Adapter) ResolveIDs(ctx context.Context, usernames []string) (map[string]string, error) {
	const op = "adapter.users.ResolveIDs"

	names, err := domainusers.NewUsernames(usernames)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	found, err := a.ByUsernames(ctx, names)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return found.IDsByUsername(), nil
}
