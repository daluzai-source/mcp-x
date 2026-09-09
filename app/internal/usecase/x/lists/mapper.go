package lists

import (
	"time"

	"github.com/Role1776/mcp-x/app/internal/domain/x/lists"
	"github.com/Role1776/mcp-x/app/internal/domain/x/posts"
	"github.com/Role1776/mcp-x/app/internal/domain/x/users"
	dto "github.com/Role1776/mcp-x/app/internal/dto/x/lists"
)

func toLists(collection lists.Lists) []dto.List {
	result := make([]dto.List, 0, collection.Len())
	for _, v := range collection.Views() {
		result = append(result, toList(v))
	}

	return result
}

func toList(v lists.ListView) dto.List {
	list := dto.List{
		ID:            v.ID,
		Name:          v.Name,
		Description:   v.Description,
		Private:       v.Private,
		OwnerID:       v.OwnerID,
		MemberCount:   v.MemberCount,
		FollowerCount: v.FollowerCount,
		URL:           v.URL,
	}

	if !v.CreatedAt.IsZero() {
		list.CreatedAt = v.CreatedAt.Format(time.RFC3339)
	}

	return list
}

func toUsers(collection users.Users) []dto.User {
	result := make([]dto.User, 0, collection.Len())
	for _, v := range collection.Views() {
		result = append(result, dto.User{
			ID:             v.ID,
			Username:       v.Username,
			Name:           v.Name,
			Verified:       v.Verified,
			FollowersCount: v.FollowersCount,
			URL:            v.URL,
		})
	}

	return result
}

func toPosts(collection posts.Posts) []dto.Post {
	result := make([]dto.Post, 0, collection.Len())
	for _, v := range collection.Views() {
		post := dto.Post{
			ID:             v.ID,
			Text:           v.Text,
			AuthorID:       v.AuthorID,
			AuthorUsername: v.AuthorUsername,
			Likes:          v.Metrics.Likes,
			Reposts:        v.Metrics.Retweets,
			URL:            v.URL,
		}

		if !v.CreatedAt.IsZero() {
			post.CreatedAt = v.CreatedAt.Format(time.RFC3339)
		}

		result = append(result, post)
	}

	return result
}
