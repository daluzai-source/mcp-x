package users

import (
	"time"

	"github.com/Role1776/mcp-x/app/internal/domain/x/users"
	dto "github.com/Role1776/mcp-x/app/internal/dto/x/users"
)

func toUsers(collection users.Users) []dto.User {
	result := make([]dto.User, 0, collection.Len())
	for _, v := range collection.Views() {
		result = append(result, toUser(v))
	}

	return result
}

func toUser(v users.UserView) dto.User {
	user := dto.User{
		ID:              v.ID,
		Username:        v.Username,
		Name:            v.Name,
		Description:     v.Description,
		ProfileImageURL: v.ProfileImageURL,
		Verified:        v.Verified,
		Protected:       v.Protected,
		FollowersCount:  v.FollowersCount,
		FollowingCount:  v.FollowingCount,
		PostCount:       v.PostCount,
		URL:             v.URL,
	}

	if !v.CreatedAt.IsZero() {
		user.CreatedAt = v.CreatedAt.Format(time.RFC3339)
	}

	return user
}
