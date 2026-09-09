package lists

import (
	"context"
	"fmt"

	"github.com/michimani/gotwi"
	listlookup "github.com/michimani/gotwi/list/listlookup"
	"github.com/michimani/gotwi/list/listlookup/types"
	"github.com/michimani/gotwi/list/listtweetlookup"
	tweettypes "github.com/michimani/gotwi/list/listtweetlookup/types"

	"github.com/Role1776/mcp-x/app/internal/adapter/x/apierr"
	"github.com/Role1776/mcp-x/app/internal/domain/x"
	"github.com/Role1776/mcp-x/app/internal/domain/x/lists"
	"github.com/Role1776/mcp-x/app/internal/domain/x/posts"
)

func (a *Adapter) Get(ctx context.Context, id lists.ListID) (lists.List, error) {
	const op = "adapter.lists.Get"

	out, err := listlookup.Get(ctx, a.client, &types.GetInput{
		ID:         id.String(),
		ListFields: listFields,
	})
	if err != nil {
		return lists.List{}, fmt.Errorf("%s: %w", op, apierr.NotFoundAs(err, x.ErrListNotFound))
	}

	// X answers a missing list with 200 and an empty data object, so the absent
	// id — not the transport error — is what marks it as not found.
	if gotwi.StringValue(out.Data.ID) == "" {
		return lists.List{}, fmt.Errorf("%s: %w", op, x.ErrListNotFound)
	}

	list, err := toList(out.Data)
	if err != nil {
		return lists.List{}, fmt.Errorf("%s: %w", op, err)
	}

	return list, nil
}

func (a *Adapter) Owned(ctx context.Context, userID string, maxResults int, token string) (lists.Lists, string, error) {
	const op = "adapter.lists.Owned"

	if userID == "" {
		userID = a.selfID
	}

	out, err := listlookup.ListOwned(ctx, a.client, &types.ListOwnedInput{
		ID:              userID,
		ListFields:      listFields,
		MaxResults:      types.ListOwnedMaxResults(maxResults),
		PaginationToken: token,
	})
	if err != nil {
		return nil, "", fmt.Errorf("%s: %w", op, apierr.NotFoundAs(err, x.ErrUserNotFound))
	}

	return toLists(out.Data), gotwi.StringValue(out.Meta.NextToken), nil
}

func (a *Adapter) Posts(ctx context.Context, id lists.ListID, maxResults int, token string) (posts.Posts, string, error) {
	const op = "adapter.lists.Posts"

	out, err := listtweetlookup.List(ctx, a.client, &tweettypes.ListInput{
		ID:              id.String(),
		MaxResults:      tweettypes.ListMaxResults(maxResults),
		Expansions:      expansions,
		TweetFields:     tweetFields,
		UserFields:      userFields,
		PaginationToken: token,
	})
	if err != nil {
		return nil, "", fmt.Errorf("%s: %w", op, apierr.NotFoundAs(err, x.ErrListNotFound))
	}

	return toPosts(out.Data, out.Includes.Users), gotwi.StringValue(out.Meta.NextToken), nil
}
