package posts

import (
	"context"
	"fmt"

	"github.com/michimani/gotwi/tweet/tweetlookup"
	"github.com/michimani/gotwi/tweet/tweetlookup/types"

	"github.com/Role1776/mcp-x/app/internal/adapter/x/apierr"
	"github.com/Role1776/mcp-x/app/internal/domain/x"
	"github.com/Role1776/mcp-x/app/internal/domain/x/posts"
)

func (a *Adapter) Lookup(ctx context.Context, ids posts.PostIDs) (posts.Posts, []posts.PartialFailure, error) {
	const op = "adapter.posts.Lookup"

	out, err := tweetlookup.List(ctx, a.client, &types.ListInput{
		IDs:         ids.Strings(),
		Expansions:  expansions,
		TweetFields: tweetFields,
		UserFields:  userFields,
		MediaFields: mediaFields,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("%s: %w", op, apierr.NotFoundAs(err, x.ErrPostNotFound))
	}

	return toPosts(out.Data, includes{Users: out.Includes.Users, Media: out.Includes.Media}),
		toPartialFailures(out.Errors), nil
}
