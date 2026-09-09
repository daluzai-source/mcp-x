package users

import (
	"context"
	"fmt"

	"github.com/michimani/gotwi"
	"github.com/michimani/gotwi/user/block"
	"github.com/michimani/gotwi/user/block/types"

	"github.com/Role1776/mcp-x/app/internal/adapter/x/apierr"
	domainusers "github.com/Role1776/mcp-x/app/internal/domain/x/users"
)

func (a *Adapter) Blocked(ctx context.Context, maxResults int, token string) (domainusers.Users, string, error) {
	const op = "adapter.users.Blocked"

	out, err := block.List(ctx, a.client, &types.ListInput{
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
