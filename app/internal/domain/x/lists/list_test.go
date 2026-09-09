package lists

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Role1776/mcp-x/app/internal/domain/x"
)

func newTestList(t *testing.T, id string) List {
	t.Helper()

	list, err := NewList(ListProps{ID: id, Name: "list " + id})
	require.NoError(t, err)

	return list
}

func TestNewList(t *testing.T) {
	cases := []struct {
		name      string
		props     ListProps
		wantErrIs error
	}{
		{name: "the id alone is enough", props: ListProps{ID: "12"}},
		{
			name: "every field populated",
			props: ListProps{
				ID:            "1234567890123456789",
				Name:          "reading",
				Description:   "books",
				Private:       true,
				OwnerID:       "1146654567674912769",
				MemberCount:   3,
				FollowerCount: 4,
				CreatedAt:     time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{name: "an absent owner is allowed", props: ListProps{ID: "12", OwnerID: ""}},
		{name: "a name at the limit", props: ListProps{ID: "12", Name: strings.Repeat("a", MaxNameLength)}},
		{
			name:  "a description at the limit",
			props: ListProps{ID: "12", Description: strings.Repeat("a", MaxDescriptionLength)},
		},
		{name: "an empty id is rejected", props: ListProps{ID: ""}, wantErrIs: x.ErrInvalidListID},
		{name: "a malformed id is rejected", props: ListProps{ID: "nope"}, wantErrIs: x.ErrInvalidListID},
		{name: "a malformed owner is rejected", props: ListProps{ID: "12", OwnerID: "jack"}, wantErrIs: x.ErrInvalidUserID},
		{
			name:      "a name above the limit is rejected",
			props:     ListProps{ID: "12", Name: strings.Repeat("a", MaxNameLength+1)},
			wantErrIs: x.ErrListNameTooLong,
		},
		{
			name:      "a description above the limit is rejected",
			props:     ListProps{ID: "12", Description: strings.Repeat("a", MaxDescriptionLength+1)},
			wantErrIs: x.ErrListDescriptionTooLong,
		},
		{
			name:      "a negative member count is rejected",
			props:     ListProps{ID: "12", MemberCount: -1},
			wantErrIs: x.ErrInvalidRequest,
		},
		{
			name:      "a negative follower count is rejected",
			props:     ListProps{ID: "12", FollowerCount: -1},
			wantErrIs: x.ErrInvalidRequest,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			list, err := NewList(c.props)

			if c.wantErrIs != nil {
				require.ErrorIs(t, err, c.wantErrIs)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, c.props.ID, list.ID().String())
		})
	}
}

func TestListView(t *testing.T) {
	t.Run("a view mirrors the props and carries a permalink", func(t *testing.T) {
		created := time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)

		list, err := NewList(ListProps{
			ID:            "12",
			Name:          "reading",
			Description:   "books",
			Private:       true,
			OwnerID:       "34",
			MemberCount:   3,
			FollowerCount: 4,
			CreatedAt:     created,
		})
		require.NoError(t, err)

		view := list.View()

		assert.Equal(t, "12", view.ID)
		assert.Equal(t, "reading", view.Name)
		assert.Equal(t, "books", view.Description)
		assert.True(t, view.Private)
		assert.Equal(t, "34", view.OwnerID)
		assert.Equal(t, 3, view.MemberCount)
		assert.Equal(t, 4, view.FollowerCount)
		assert.Equal(t, created, view.CreatedAt)
		assert.Equal(t, "https://x.com/i/lists/12", view.URL)
	})

	t.Run("an absent owner comes back as an empty string", func(t *testing.T) {
		list, err := NewList(ListProps{ID: "12"})
		require.NoError(t, err)

		assert.Empty(t, list.View().OwnerID)
	})
}

func TestListsLenAndIsEmpty(t *testing.T) {
	cases := []struct {
		name        string
		lists       Lists
		wantLen     int
		wantIsEmpty bool
	}{
		{name: "a nil collection", lists: nil, wantLen: 0, wantIsEmpty: true},
		{name: "an empty collection", lists: Lists{}, wantLen: 0, wantIsEmpty: true},
		{name: "one list", lists: Lists{newTestList(t, "1")}, wantLen: 1},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			assert.Equal(t, c.wantLen, c.lists.Len())
			assert.Equal(t, c.wantIsEmpty, c.lists.IsEmpty())
		})
	}
}

func TestListsLimit(t *testing.T) {
	all := Lists{newTestList(t, "1"), newTestList(t, "2"), newTestList(t, "3")}

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

func TestListsDedupe(t *testing.T) {
	first := newTestList(t, "1")
	second := newTestList(t, "2")

	cases := []struct {
		name    string
		lists   Lists
		wantIDs []string
	}{
		{name: "a repeated id collapses to the first occurrence", lists: Lists{first, second, first}, wantIDs: []string{"1", "2"}},
		{name: "distinct lists are untouched", lists: Lists{first, second}, wantIDs: []string{"1", "2"}},
		{name: "an empty collection stays empty", lists: nil, wantIDs: []string{}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			deduped := c.lists.Dedupe()

			ids := make([]string, 0, deduped.Len())
			for _, list := range deduped {
				ids = append(ids, list.ID().String())
			}

			assert.Equal(t, c.wantIDs, ids)
		})
	}
}

func TestListsViews(t *testing.T) {
	views := Lists{newTestList(t, "1"), newTestList(t, "2")}.Views()

	require.Len(t, views, 2)
	assert.Equal(t, "1", views[0].ID)
	assert.Empty(t, Lists(nil).Views())
}
