package posts

import (
	"context"
	"fmt"

	"github.com/michimani/gotwi"
	"github.com/michimani/gotwi/tweet/bookmark"
	"github.com/michimani/gotwi/tweet/bookmark/types"

	"github.com/Role1776/mcp-x/app/internal/adapter/x/apierr"
	"github.com/Role1776/mcp-x/app/internal/domain/x/posts"
)

func (a *Adapter) Bookmarks(ctx context.Context, maxResults int, token string) (posts.Posts, string, error) {
	const op = "adapter.posts.Bookmarks"

	out, err := bookmark.List(ctx, a.client, &types.ListInput{
		ID:              a.selfID,
		MaxResults:      types.ListMaxResults(maxResults),
		Expansions:      expansions,
		TweetFields:     tweetFields,
		UserFields:      userFields,
		MediaFields:     mediaFields,
		PaginationToken: token,
	})
	if err != nil {
		return nil, "", fmt.Errorf("%s: %w", op, apierr.ToDomain(err))
	}

	return toPosts(out.Data, includes{Media: out.Includes.Media}), gotwi.StringValue(out.Meta.NextToken), nil
}
