package posts

import (
	"time"

	"github.com/Role1776/mcp-x/app/internal/domain/x/posts"
	"github.com/Role1776/mcp-x/app/internal/domain/x/users"
	dto "github.com/Role1776/mcp-x/app/internal/dto/x/posts"
)

func toPosts(collection posts.Posts) []dto.Post {
	result := make([]dto.Post, 0, collection.Len())
	for _, v := range collection.Views() {
		result = append(result, toPost(v))
	}

	return result
}

func toPost(v posts.PostView) dto.Post {
	post := dto.Post{
		ID:             v.ID,
		Text:           v.Text,
		AuthorID:       v.AuthorID,
		AuthorUsername: v.AuthorUsername,
		Lang:           v.Lang,
		ConversationID: v.ConversationID,
		Metrics: dto.Metrics{
			Retweets: v.Metrics.Retweets,
			Replies:  v.Metrics.Replies,
			Likes:    v.Metrics.Likes,
			Quotes:   v.Metrics.Quotes,
		},
		URLs:      v.URLs,
		MediaKeys: v.MediaKeys,
		URL:       v.URL,
	}

	if !v.CreatedAt.IsZero() {
		post.CreatedAt = v.CreatedAt.Format(time.RFC3339)
	}

	for _, ref := range v.Referenced {
		post.Referenced = append(post.Referenced, dto.Reference{Type: ref.Type, ID: ref.ID})
	}

	return post
}

func toUsers(collection users.Users) []dto.User {
	result := make([]dto.User, 0, collection.Len())
	for _, v := range collection.Views() {
		result = append(result, dto.User{
			ID:             v.ID,
			Username:       v.Username,
			Name:           v.Name,
			Description:    v.Description,
			Verified:       v.Verified,
			Protected:      v.Protected,
			FollowersCount: v.FollowersCount,
			URL:            v.URL,
		})
	}

	return result
}
