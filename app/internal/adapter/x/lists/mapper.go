package lists

import (
	"github.com/michimani/gotwi"
	"github.com/michimani/gotwi/resources"

	"github.com/Role1776/mcp-x/app/internal/domain/x/lists"
	"github.com/Role1776/mcp-x/app/internal/domain/x/posts"
	"github.com/Role1776/mcp-x/app/internal/domain/x/users"
)

func toLists(data []resources.List) lists.Lists {
	result := make(lists.Lists, 0, len(data))
	for _, l := range data {
		list, err := toList(l)
		if err != nil {
			continue
		}
		result = append(result, list)
	}

	return result
}

func toList(l resources.List) (lists.List, error) {
	return lists.NewList(lists.ListProps{
		ID:            gotwi.StringValue(l.ID),
		Name:          gotwi.StringValue(l.Name),
		Description:   gotwi.StringValue(l.Description),
		Private:       gotwi.BoolValue(l.Private),
		OwnerID:       gotwi.StringValue(l.OwnerID),
		MemberCount:   gotwi.IntValue(l.MemberCount),
		FollowerCount: gotwi.IntValue(l.FollowerCount),
		CreatedAt:     gotwi.TimeValue(l.CreatedAt),
	})
}

func toUsers(data []resources.User) users.Users {
	result := make(users.Users, 0, len(data))
	for _, u := range data {
		props := users.UserProps{
			ID:              gotwi.StringValue(u.ID),
			Username:        gotwi.StringValue(u.Username),
			Name:            gotwi.StringValue(u.Name),
			Description:     gotwi.StringValue(u.Description),
			ProfileImageURL: gotwi.StringValue(u.ProfileImageURL),
			Verified:        gotwi.BoolValue(u.Verified),
			Protected:       gotwi.BoolValue(u.Protected),
			CreatedAt:       gotwi.TimeValue(u.CreatedAt),
		}

		if u.PublicMetrics != nil {
			props.FollowersCount = gotwi.IntValue(u.PublicMetrics.FollowersCount)
			props.FollowingCount = gotwi.IntValue(u.PublicMetrics.FollowingCount)
			props.PostCount = gotwi.IntValue(u.PublicMetrics.TweetCount)
		}

		user, err := users.NewUser(props)
		if err != nil {
			continue
		}
		result = append(result, user)
	}

	return result
}

func toPosts(data []resources.Tweet, included []resources.User) posts.Posts {
	authors := make(map[string]string, len(included))
	for _, user := range included {
		authors[gotwi.StringValue(user.ID)] = gotwi.StringValue(user.Username)
	}

	result := make(posts.Posts, 0, len(data))
	for _, tweet := range data {
		authorID := gotwi.StringValue(tweet.AuthorID)

		props := posts.PostProps{
			ID:             gotwi.StringValue(tweet.ID),
			Text:           tweetText(tweet),
			AuthorID:       authorID,
			AuthorUsername: authors[authorID],
			CreatedAt:      gotwi.TimeValue(tweet.CreatedAt),
			Lang:           gotwi.StringValue(tweet.Lang),
			ConversationID: gotwi.StringValue(tweet.ConversationID),
		}

		if tweet.PublicMetrics != nil {
			props.Metrics = posts.Metrics{
				Retweets: gotwi.IntValue(tweet.PublicMetrics.RetweetCount),
				Replies:  gotwi.IntValue(tweet.PublicMetrics.ReplyCount),
				Likes:    gotwi.IntValue(tweet.PublicMetrics.LikeCount),
				Quotes:   gotwi.IntValue(tweet.PublicMetrics.QuoteCount),
			}
		}

		if tweet.Attachments != nil {
			props.MediaKeys = tweet.Attachments.MediaKeys
		}

		post, err := posts.NewPost(props)
		if err != nil {
			continue
		}
		result = append(result, post)
	}

	return result
}

func tweetText(tweet resources.Tweet) string {
	if tweet.NoteTweet != nil && tweet.NoteTweet.Text != nil {
		return gotwi.StringValue(tweet.NoteTweet.Text)
	}

	return gotwi.StringValue(tweet.Text)
}
