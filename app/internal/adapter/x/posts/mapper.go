package posts

import (
	"github.com/michimani/gotwi"
	"github.com/michimani/gotwi/resources"

	"github.com/Role1776/mcp-x/app/internal/domain/x/posts"
	"github.com/Role1776/mcp-x/app/internal/domain/x/users"
)

type includes struct {
	Users []resources.User
	Media []resources.Media
}

func toPosts(data []resources.Tweet, inc includes) posts.Posts {
	authors := make(map[string]string, len(inc.Users))
	for _, user := range inc.Users {
		authors[gotwi.StringValue(user.ID)] = gotwi.StringValue(user.Username)
	}

	result := make(posts.Posts, 0, len(data))
	for _, tweet := range data {
		post, err := toPost(tweet, authors)
		if err != nil {
			continue
		}
		result = append(result, post)
	}

	return result
}

func toPost(tweet resources.Tweet, authors map[string]string) (posts.Post, error) {
	authorID := gotwi.StringValue(tweet.AuthorID)

	props := posts.PostProps{
		ID:             gotwi.StringValue(tweet.ID),
		Text:           text(tweet),
		AuthorID:       authorID,
		AuthorUsername: authors[authorID],
		CreatedAt:      gotwi.TimeValue(tweet.CreatedAt),
		Lang:           gotwi.StringValue(tweet.Lang),
		ConversationID: gotwi.StringValue(tweet.ConversationID),
		Metrics:        metrics(tweet.PublicMetrics),
		URLs:           urls(tweet.Entities),
		MediaKeys:      mediaKeys(tweet.Attachments),
		Referenced:     referenced(tweet.ReferencedTweets),
	}

	return posts.NewPost(props)
}

// text prefers note_tweet: for posts longer than 280 characters the API sends a
// truncated "text" and puts the whole thing in note_tweet.
func text(tweet resources.Tweet) string {
	if tweet.NoteTweet != nil && tweet.NoteTweet.Text != nil {
		return gotwi.StringValue(tweet.NoteTweet.Text)
	}

	return gotwi.StringValue(tweet.Text)
}

func metrics(m *resources.TweetPublicMetrics) posts.Metrics {
	if m == nil {
		return posts.Metrics{}
	}

	return posts.Metrics{
		Retweets: gotwi.IntValue(m.RetweetCount),
		Replies:  gotwi.IntValue(m.ReplyCount),
		Likes:    gotwi.IntValue(m.LikeCount),
		Quotes:   gotwi.IntValue(m.QuoteCount),
	}
}

func urls(entities *resources.TweetEntities) []string {
	if entities == nil {
		return nil
	}

	result := make([]string, 0, len(entities.URLs))
	for _, u := range entities.URLs {
		expanded := gotwi.StringValue(u.ExpandedURL)
		if expanded == "" {
			expanded = gotwi.StringValue(u.URL)
		}
		if expanded != "" {
			result = append(result, expanded)
		}
	}

	return result
}

func mediaKeys(attachments *resources.TweetAttachments) []string {
	if attachments == nil {
		return nil
	}

	return attachments.MediaKeys
}

func referenced(refs []resources.ReferencedTweet) []posts.Reference {
	if len(refs) == 0 {
		return nil
	}

	result := make([]posts.Reference, 0, len(refs))
	for _, ref := range refs {
		result = append(result, posts.Reference{
			Type: gotwi.StringValue(ref.Type),
			ID:   gotwi.StringValue(ref.ID),
		})
	}

	return result
}

func toPartialFailures(errs []resources.PartialError) []posts.PartialFailure {
	if len(errs) == 0 {
		return nil
	}

	result := make([]posts.PartialFailure, 0, len(errs))
	for _, e := range errs {
		result = append(result, posts.PartialFailure{
			ResourceID: gotwi.StringValue(e.ResourceID),
			Title:      gotwi.StringValue(e.Title),
			Detail:     gotwi.StringValue(e.Detail),
		})
	}

	return result
}

func toUsers(data []resources.User) users.Users {
	result := make(users.Users, 0, len(data))
	for _, u := range data {
		user, err := users.NewUser(users.UserProps{
			ID:              gotwi.StringValue(u.ID),
			Username:        gotwi.StringValue(u.Username),
			Name:            gotwi.StringValue(u.Name),
			Description:     gotwi.StringValue(u.Description),
			ProfileImageURL: gotwi.StringValue(u.ProfileImageURL),
			Verified:        gotwi.BoolValue(u.Verified),
			Protected:       gotwi.BoolValue(u.Protected),
			FollowersCount:  publicMetric(u.PublicMetrics, followers),
			FollowingCount:  publicMetric(u.PublicMetrics, following),
			PostCount:       publicMetric(u.PublicMetrics, tweets),
			CreatedAt:       gotwi.TimeValue(u.CreatedAt),
		})
		if err != nil {
			continue
		}
		result = append(result, user)
	}

	return result
}

const (
	followers = iota
	following
	tweets
)

func publicMetric(m *resources.UserPublicMetrics, kind int) int {
	if m == nil {
		return 0
	}

	switch kind {
	case followers:
		return gotwi.IntValue(m.FollowersCount)
	case following:
		return gotwi.IntValue(m.FollowingCount)
	case tweets:
		return gotwi.IntValue(m.TweetCount)
	default:
		return 0
	}
}
