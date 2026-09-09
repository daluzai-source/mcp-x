package users

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Role1776/mcp-x/app/internal/domain/x"
)

func newTestUser(t *testing.T, id, username string) User {
	t.Helper()

	user, err := NewUser(UserProps{ID: id, Username: username})
	require.NoError(t, err)

	return user
}

func TestNewUser(t *testing.T) {
	cases := []struct {
		name      string
		props     UserProps
		wantErrIs error
	}{
		{name: "id and username are enough", props: UserProps{ID: "12", Username: "jack"}},
		{
			name: "every field populated",
			props: UserProps{
				ID:              "1234567890123456789",
				Username:        "jack",
				Name:            "Jack",
				Description:     "hi",
				ProfileImageURL: "https://example.com/a.jpg",
				Verified:        true,
				Protected:       true,
				FollowersCount:  10,
				FollowingCount:  20,
				PostCount:       30,
				CreatedAt:       time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{name: "zero counts are valid", props: UserProps{ID: "12", Username: "jack"}},
		{name: "an empty id is rejected", props: UserProps{ID: "", Username: "jack"}, wantErrIs: x.ErrInvalidUserID},
		{name: "a malformed id is rejected", props: UserProps{ID: "nope", Username: "jack"}, wantErrIs: x.ErrInvalidUserID},
		{name: "an empty username is rejected", props: UserProps{ID: "12", Username: ""}, wantErrIs: x.ErrInvalidUsername},
		{
			name:      "a malformed username is rejected",
			props:     UserProps{ID: "12", Username: "ja ck"},
			wantErrIs: x.ErrInvalidUsername,
		},
		{
			name:      "a negative follower count is rejected",
			props:     UserProps{ID: "12", Username: "jack", FollowersCount: -1},
			wantErrIs: x.ErrInvalidRequest,
		},
		{
			name:      "a negative post count is rejected",
			props:     UserProps{ID: "12", Username: "jack", PostCount: -1},
			wantErrIs: x.ErrInvalidRequest,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			user, err := NewUser(c.props)

			if c.wantErrIs != nil {
				require.ErrorIs(t, err, c.wantErrIs)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, c.props.ID, user.ID().String())
			assert.Equal(t, c.props.Username, user.Username().String())
		})
	}
}

func TestUserView(t *testing.T) {
	user, err := NewUser(UserProps{
		ID:             "12",
		Username:       "JaCk",
		Name:           "Jack",
		Verified:       true,
		FollowersCount: 7,
	})
	require.NoError(t, err)

	view := user.View()

	assert.Equal(t, "12", view.ID)
	assert.Equal(t, "JaCk", view.Username)
	assert.Equal(t, "Jack", view.Name)
	assert.True(t, view.Verified)
	assert.Equal(t, 7, view.FollowersCount)
	assert.Equal(t, "https://x.com/JaCk", view.URL)
}

func TestUsersLenAndIsEmpty(t *testing.T) {
	cases := []struct {
		name        string
		users       Users
		wantLen     int
		wantIsEmpty bool
	}{
		{name: "a nil collection", users: nil, wantLen: 0, wantIsEmpty: true},
		{name: "an empty collection", users: Users{}, wantLen: 0, wantIsEmpty: true},
		{name: "one user", users: Users{newTestUser(t, "1", "jack")}, wantLen: 1},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			assert.Equal(t, c.wantLen, c.users.Len())
			assert.Equal(t, c.wantIsEmpty, c.users.IsEmpty())
		})
	}
}

func TestUsersLimit(t *testing.T) {
	all := Users{newTestUser(t, "1", "a"), newTestUser(t, "2", "b"), newTestUser(t, "3", "c")}

	cases := []struct {
		name    string
		n       int
		wantLen int
	}{
		{name: "fewer than available", n: 2, wantLen: 2},
		{name: "exactly as many as available", n: 3, wantLen: 3},
		{name: "more than available", n: 10, wantLen: 3},
		{name: "zero returns nothing", n: 0, wantLen: 0},
		{name: "a negative limit returns nothing", n: -1, wantLen: 0},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			assert.Equal(t, c.wantLen, all.Limit(c.n).Len())
		})
	}
}

func TestUsersDedupe(t *testing.T) {
	first := newTestUser(t, "1", "jack")
	second := newTestUser(t, "2", "jill")

	cases := []struct {
		name    string
		users   Users
		wantIDs []string
	}{
		{name: "a repeated id collapses to the first occurrence", users: Users{first, second, first}, wantIDs: []string{"1", "2"}},
		{name: "distinct users are untouched", users: Users{first, second}, wantIDs: []string{"1", "2"}},
		{name: "an empty collection stays empty", users: nil, wantIDs: []string{}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			deduped := c.users.Dedupe()

			ids := make([]string, 0, deduped.Len())
			for _, user := range deduped {
				ids = append(ids, user.ID().String())
			}

			assert.Equal(t, c.wantIDs, ids)
		})
	}
}

func TestUsersByUsername(t *testing.T) {
	found := Users{newTestUser(t, "1", "JaCk"), newTestUser(t, "2", "jill")}

	cases := []struct {
		name    string
		lookFor string
		wantID  string
		wantOk  bool
	}{
		{name: "an exact match", lookFor: "jill", wantID: "2", wantOk: true},
		{name: "a handle typed in another case still matches", lookFor: "jack", wantID: "1", wantOk: true},
		{name: "an absent handle", lookFor: "nobody", wantOk: false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			name, err := NewUsername(c.lookFor)
			require.NoError(t, err)

			user, ok := found.ByUsername(name)

			require.Equal(t, c.wantOk, ok)

			if !c.wantOk {
				return
			}

			assert.Equal(t, c.wantID, user.ID().String())
		})
	}
}

func TestUsersByUsernameOnAnEmptyCollection(t *testing.T) {
	name, err := NewUsername("jack")
	require.NoError(t, err)

	_, ok := Users(nil).ByUsername(name)

	assert.False(t, ok)
}

func TestUsersIDsByUsername(t *testing.T) {
	t.Run("keys are normalized handles", func(t *testing.T) {
		index := Users{newTestUser(t, "1", "JaCk"), newTestUser(t, "2", "jill")}.IDsByUsername()

		assert.Equal(t, map[string]string{"jack": "1", "jill": "2"}, index)
	})

	t.Run("an empty collection yields a usable empty map", func(t *testing.T) {
		index := Users(nil).IDsByUsername()

		require.NotNil(t, index)
		assert.Empty(t, index)
	})
}

func TestUsersViews(t *testing.T) {
	views := Users{newTestUser(t, "1", "jack"), newTestUser(t, "2", "jill")}.Views()

	require.Len(t, views, 2)
	assert.Equal(t, "jack", views[0].Username)
	assert.Empty(t, Users(nil).Views())
}
