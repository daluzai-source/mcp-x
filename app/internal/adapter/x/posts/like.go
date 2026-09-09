package posts

import (
	"context"
	"fmt"

	"github.com/michimani/gotwi"
	"github.com/michimani/gotwi/tweet/like"
	"github.com/michimani/gotwi/tweet/like/types"

	"github.com/Role1776/mcp-x/app/internal/adapter/x/apierr"
	"github.com/Role1776/mcp-x/app/internal/domain/x"
	"github.com/Role1776/mcp-x/app/internal/domain/x/posts"
	"github.com/Role1776/mcp-x/app/internal/domain/x/users"
)

func (a *Adapter) Like(ctx context.Context, id posts.PostID) error {
	const op = "adapter.posts.Like"

	if _, err := like.Create(ctx, a.client, &types.CreateInput{ID: a.selfID, TweetID: id.String()}); err != nil {
		return fmt.Errorf("%s: %w", op, apierr.NotFoundAs(err, x.ErrPostNotFound))
	}

	return nil
}

func (a *Adapter) Unlike(ctx context.Context, id posts.PostID) error {
	const op = "adapter.posts.Unlike"

	if _, err := like.Delete(ctx, a.client, &types.DeleteInput{ID: a.selfID, TweetID: id.String()}); err != nil {
		return fmt.Errorf("%s: %w", op, apierr.NotFoundAs(err, x.ErrPostNotFound))
	}

	return nil
}

func (a *Adapter) LikedBy(ctx context.Context, id posts.PostID, maxResults int, token string) (users.Users, error) {
	const op = "adapter.posts.LikedBy"

	out, err := like.ListUsers(ctx, a.client, &types.ListUsersInput{
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

func (a *Adapter) Liked(ctx context.Context, maxResults int, token string) (posts.Posts, string, error) {
	const op = "adapter.posts.Liked"

	out, err := like.List(ctx, a.client, &types.ListInput{
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
