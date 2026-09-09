package posts

import (
	"context"
	"fmt"
	"time"

	"github.com/michimani/gotwi"
	"github.com/michimani/gotwi/fields"
	"github.com/michimani/gotwi/tweet/quotetweet"
	quotetypes "github.com/michimani/gotwi/tweet/quotetweet/types"
	"github.com/michimani/gotwi/tweet/timeline"
	"github.com/michimani/gotwi/tweet/timeline/types"

	"github.com/Role1776/mcp-x/app/internal/adapter/x/apierr"
	"github.com/Role1776/mcp-x/app/internal/domain/x"
	"github.com/Role1776/mcp-x/app/internal/domain/x/posts"
)

type TimelineParams struct {
	UserID     string
	MaxResults int
	Exclude    posts.Exclude
	StartTime  time.Time
	Token      string
}

func (a *Adapter) UserTimeline(ctx context.Context, p TimelineParams) (posts.Posts, string, error) {
	const op = "adapter.posts.UserTimeline"

	in := &types.ListTweetsInput{
		ID:              p.UserID,
		MaxResults:      types.ListMaxResults(p.MaxResults),
		Exclude:         excludeList(p.Exclude),
		Expansions:      expansions,
		TweetFields:     tweetFields,
		UserFields:      userFields,
		MediaFields:     mediaFields,
		PaginationToken: p.Token,
	}

	if !p.StartTime.IsZero() {
		in.StartTime = gotwi.Time(p.StartTime)
	}

	out, err := timeline.ListTweets(ctx, a.client, in)
	if err != nil {
		return nil, "", fmt.Errorf("%s: %w", op, apierr.NotFoundAs(err, x.ErrUserNotFound))
	}

	return toPosts(out.Data, includes{Users: out.Includes.Users, Media: out.Includes.Media}),
		gotwi.StringValue(out.Meta.NextToken), nil
}

func (a *Adapter) Mentions(ctx context.Context, maxResults int, token string) (posts.Posts, string, error) {
	const op = "adapter.posts.Mentions"

	out, err := timeline.ListMentions(ctx, a.client, &types.ListMentionsInput{
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

	return toPosts(out.Data, includes{Users: out.Includes.Users, Media: out.Includes.Media}),
		gotwi.StringValue(out.Meta.NextToken), nil
}

func (a *Adapter) Home(ctx context.Context, maxResults int, token string) (posts.Posts, string, error) {
	const op = "adapter.posts.Home"

	out, err := timeline.ListReverseChronological(ctx, a.client, &types.ListReverseChronologicalInput{
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

	return toPosts(out.Data, includes{Users: out.Includes.Users, Media: out.Includes.Media}),
		gotwi.StringValue(out.Meta.NextToken), nil
}

func (a *Adapter) Quotes(ctx context.Context, id posts.PostID, maxResults int, token string) (posts.Posts, string, error) {
	const op = "adapter.posts.Quotes"

	out, err := quotetweet.List(ctx, a.client, &quotetypes.ListInput{
		ID:              id.String(),
		MaxResults:      quotetypes.ListMaxResults(maxResults),
		Expansions:      expansions,
		TweetFields:     tweetFields,
		UserFields:      userFields,
		MediaFields:     mediaFields,
		PaginationToken: token,
	})
	if err != nil {
		return nil, "", fmt.Errorf("%s: %w", op, apierr.NotFoundAs(err, x.ErrPostNotFound))
	}

	return toPosts(out.Data, includes{Users: out.Includes.Users}),
		gotwi.StringValue(out.Meta.NextToken), nil
}

func excludeList(exclude posts.Exclude) fields.ExcludeList {
	list := make(fields.ExcludeList, 0, 2)
	if exclude.Replies() {
		list = append(list, fields.ExcludeReplies)
	}
	if exclude.Retweets() {
		list = append(list, fields.ExcludeRetweets)
	}

	return list
}
