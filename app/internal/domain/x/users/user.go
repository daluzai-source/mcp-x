package users

import (
	"fmt"
	"time"

	"github.com/Role1776/mcp-x/app/internal/domain/x"
	"github.com/Role1776/mcp-x/app/internal/pkg/validator"
)

const (
	xBaseURL = "https://x.com/%s"
)

type User struct {
	id              UserID
	username        Username
	name            string
	description     string
	profileImageURL string
	verified        bool
	protected       bool
	followersCount  int
	followingCount  int
	postCount       int
	createdAt       time.Time
}

type UserProps struct {
	ID              string
	Username        string
	Name            string
	Description     string
	ProfileImageURL string
	Verified        bool
	Protected       bool
	FollowersCount  int `validate:"gte=0"`
	FollowingCount  int `validate:"gte=0"`
	PostCount       int `validate:"gte=0"`
	CreatedAt       time.Time
}

func NewUser(props UserProps) (User, error) {
	id, err := NewUserID(props.ID)
	if err != nil {
		return User{}, err
	}

	username, err := NewUsername(props.Username)
	if err != nil {
		return User{}, err
	}

	if err := validator.Validate(props); err != nil {
		return User{}, x.ErrInvalidRequest
	}

	return User{
		id:              id,
		username:        username,
		name:            props.Name,
		description:     props.Description,
		profileImageURL: props.ProfileImageURL,
		verified:        props.Verified,
		protected:       props.Protected,
		followersCount:  props.FollowersCount,
		followingCount:  props.FollowingCount,
		postCount:       props.PostCount,
		createdAt:       props.CreatedAt,
	}, nil
}

type UserView struct {
	ID              string
	Username        string
	Name            string
	Description     string
	ProfileImageURL string
	Verified        bool
	Protected       bool
	FollowersCount  int
	FollowingCount  int
	PostCount       int
	CreatedAt       time.Time
	URL             string
}

func (u User) View() UserView {
	return UserView{
		ID:              u.id.String(),
		Username:        u.username.String(),
		Name:            u.name,
		Description:     u.description,
		ProfileImageURL: u.profileImageURL,
		Verified:        u.verified,
		Protected:       u.protected,
		FollowersCount:  u.followersCount,
		FollowingCount:  u.followingCount,
		PostCount:       u.postCount,
		CreatedAt:       u.createdAt,
		URL:             fmt.Sprintf(xBaseURL, u.username.String()),
	}
}

func (u User) ID() UserID {
	return u.id
}

func (u User) Username() Username {
	return u.username
}

type Users []User

func (u Users) Len() int {
	return len(u)
}

func (u Users) IsEmpty() bool {
	return len(u) == 0
}

func (u Users) Views() []UserView {
	views := make([]UserView, 0, len(u))
	for _, user := range u {
		views = append(views, user.View())
	}

	return views
}

func (u Users) Limit(n int) Users {
	if n <= 0 || len(u) == 0 {
		return Users{}
	}

	result := make(Users, min(len(u), n))
	copy(result, u)

	return result
}

func (u Users) Dedupe() Users {
	seen := make(map[string]struct{}, len(u))
	result := make(Users, 0, len(u))

	for _, user := range u {
		key := user.id.String()
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		result = append(result, user)
	}

	return result
}

func (u Users) ByUsername(name Username) (User, bool) {
	for _, user := range u {
		if user.username.Normalized() == name.Normalized() {
			return user, true
		}
	}

	return User{}, false
}

func (u Users) IDsByUsername() map[string]string {
	result := make(map[string]string, len(u))
	for _, user := range u {
		result[user.username.Normalized()] = user.id.String()
	}

	return result
}
