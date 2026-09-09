package users

import (
	"github.com/michimani/gotwi"
	"github.com/michimani/gotwi/resources"

	domainusers "github.com/Role1776/mcp-x/app/internal/domain/x/users"
)

func toUsers(data []resources.User) domainusers.Users {
	result := make(domainusers.Users, 0, len(data))
	for _, u := range data {
		user, err := toUser(u)
		if err != nil {
			continue
		}
		result = append(result, user)
	}

	return result
}

func toUser(u resources.User) (domainusers.User, error) {
	props := domainusers.UserProps{
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

	return domainusers.NewUser(props)
}
