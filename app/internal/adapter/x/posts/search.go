package posts

import (
	"context"
	"fmt"
	"time"

	"github.com/michimani/gotwi"
	"github.com/michimani/gotwi/tweet/searchtweet"
	"github.com/michimani/gotwi/tweet/searchtweet/types"
	"github.com/michimani/gotwi/tweet/tweetcount"
	counttypes "github.com/michimani/gotwi/tweet/tweetcount/types"

	"github.com/Role1776/mcp-x/app/internal/adapter/x/apierr"
	"github.com/Role1776/mcp-x/app/internal/domain/x/posts"
)

type SearchParams struct {
	Query      posts.Query
	MaxResults int
	SortOrder  string
	StartTime  time.Time
	NextToken  string
}

func (a *Adapter) Search(ctx context.Context, p SearchParams) (posts.Posts, string, error) {
	const op = "adapter.posts.Search"

	in := &types.ListRecentInput{
		Query:       p.Query.String(),
		MaxResults:  types.ListMaxResults(p.MaxResults),
		Expansions:  expansions,
		TweetFields: tweetFields,
		UserFields:  userFields,
		MediaFields: mediaFields,
		NextToken:   p.NextToken,
	}

	if !p.StartTime.IsZero() {
		in.StartTime = gotwi.Time(p.StartTime)
	}

	if p.SortOrder != "" {
		in.SortOrder = types.ListSortOrder(p.SortOrder)
	}

	out, err := searchtweet.ListRecent(ctx, a.client, in)
	if err != nil {
		return nil, "", fmt.Errorf("%s: %w", op, apierr.ToDomain(err))
	}

	return toPosts(out.Data, includes{Users: out.Includes.Users, Media: out.Includes.Media}),
		gotwi.StringValue(out.Meta.NextToken), nil
}

type CountParams struct {
	Query       posts.Query
	Granularity string
	StartTime   time.Time
}

func (a *Adapter) Count(ctx context.Context, p CountParams) (posts.Counts, error) {
	const op = "adapter.posts.Count"

	in := &counttypes.ListRecentInput{
		Query:       p.Query.String(),
		Granularity: counttypes.TweetCountsGranularity(p.Granularity),
	}

	if !p.StartTime.IsZero() {
		in.StartTime = gotwi.Time(p.StartTime)
	}

	out, err := tweetcount.ListRecent(ctx, a.client, in)
	if err != nil {
		return posts.Counts{}, fmt.Errorf("%s: %w", op, apierr.ToDomain(err))
	}

	buckets := make([]posts.Count, 0, len(out.Data))
	for _, bucket := range out.Data {
		buckets = append(buckets, posts.Count{
			Start: gotwi.TimeValue(bucket.Start),
			End:   gotwi.TimeValue(bucket.End),
			Total: gotwi.IntValue(bucket.TweetCount),
		})
	}

	return posts.Counts{Buckets: buckets, Total: gotwi.IntValue(out.Meta.TotalTweetCount)}, nil
}
