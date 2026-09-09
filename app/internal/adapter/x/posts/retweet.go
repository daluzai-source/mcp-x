package posts

import (
	"context"
	"fmt"

	"github.com/michimani/gotwi/tweet/retweet"
	"github.com/michimani/gotwi/tweet/retweet/types"

	"github.com/Role1776/mcp-x/app/internal/adapter/x/apierr"
	"github.com/Role1776/mcp-x/app/internal/domain/x"
	"github.com/Role1776/mcp-x/app/internal/domain/x/posts"
	"github.com/Role1776/mcp-x/app/internal/domain/x/users"
)

func (a *Adapter) Repost(ctx context.Context, id posts.PostID) error {
	const op = "adapter.posts.Repost"

	if _, err := retweet.Create(ctx, a.client, &types.CreateInput{ID: a.selfID, TweetID: id.String()}); err != nil {
		return fmt.Errorf("%s: %w", op, apierr.NotFoundAs(err, x.ErrPostNotFound))
	}

	return nil
}

func (a *Adapter) Unrepost(ctx context.Context, id posts.PostID) error {
	const op = "adapter.posts.Unrepost"

	if _, err := retweet.Delete(ctx, a.client, &types.DeleteInput{ID: a.selfID, SourceTweetID: id.String()}); err != nil {
		return fmt.Errorf("%s: %w", op, apierr.NotFoundAs(err, x.ErrPostNotFound))
	}

	return nil
}

func (a *Adapter) RepostedBy(ctx context.Context, id posts.PostID, maxResults int, token string) (users.Users, error) {
	const op = "adapter.posts.RepostedBy"

	out, err := retweet.ListUsers(ctx, a.client, &types.ListUsersInput{
		ID:              id.String(),
		MaxResults:      types.ListUsersMaxResults(maxResults),
		UserFields:      listUserFields,
		PaginationToken: token,
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, apierr.NotFoundAs(err, x.ErrPostNotFound))
	}

	return toUsers(out.Data), nil
}
