package lists

import (
	"context"
	"fmt"

	"github.com/michimani/gotwi"
	"github.com/michimani/gotwi/list/listmember"
	"github.com/michimani/gotwi/list/listmember/types"

	"github.com/Role1776/mcp-x/app/internal/adapter/x/apierr"
	"github.com/Role1776/mcp-x/app/internal/domain/x"
	domainlists "github.com/Role1776/mcp-x/app/internal/domain/x/lists"
	"github.com/Role1776/mcp-x/app/internal/domain/x/users"
)

func (a *Adapter) MemberAdd(ctx context.Context, id domainlists.ListID, userID string) error {
	const op = "adapter.lists.MemberAdd"

	if _, err := listmember.Create(ctx, a.client, &types.CreateInput{ID: id.String(), UserID: userID}); err != nil {
		return fmt.Errorf("%s: %w", op, apierr.NotFoundAs(err, x.ErrListNotFound))
	}

	return nil
}

func (a *Adapter) MemberRemove(ctx context.Context, id domainlists.ListID, userID string) error {
	const op = "adapter.lists.MemberRemove"

	if _, err := listmember.Delete(ctx, a.client, &types.DeleteInput{ID: id.String(), UserID: userID}); err != nil {
		return fmt.Errorf("%s: %w", op, apierr.NotFoundAs(err, x.ErrListNotFound))
	}

	return nil
}

func (a *Adapter) Members(ctx context.Context, id domainlists.ListID, maxResults int, token string) (users.Users, string, error) {
	const op = "adapter.lists.Members"

	out, err := listmember.List(ctx, a.client, &types.ListInput{
		ID:              id.String(),
		UserFields:      userFields,
		MaxResults:      types.ListMembersGetMaxResults(maxResults),
		PaginationToken: token,
	})
	if err != nil {
		return nil, "", fmt.Errorf("%s: %w", op, apierr.NotFoundAs(err, x.ErrListNotFound))
	}

	return toUsers(out.Data), gotwi.StringValue(out.Meta.NextToken), nil
}
