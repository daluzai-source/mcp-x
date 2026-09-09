package posts

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Role1776/mcp-x/app/internal/domain/x"
)

func newTestPost(t *testing.T, id string) Post {
	t.Helper()

	post, err := NewPost(PostProps{ID: id, Text: "text " + id, AuthorUsername: "jack"})
	require.NoError(t, err)

	return post
}

func TestNewPost(t *testing.T) {
	cases := []struct {
		name      string
		props     PostProps
		wantErrIs error
	}{
		{name: "the id alone is enough", props: PostProps{ID: postID}},
		{
			name: "every field populated",
			props: PostProps{
				ID:             postID,
				Text:           "hello",
				AuthorID:       "1146654567674912769",
				AuthorUsername: "jack",
				CreatedAt:      time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC),
				Lang:           "en",
				ConversationID: postID,
				Metrics:        Metrics{Likes: 3},
				URLs:           []string{"https://example.com"},
				MediaKeys:      []string{"3_1"},
				Referenced:     []Reference{{Type: "replied_to", ID: otherPostID}},
			},
		},
		{name: "an absent author is allowed", props: PostProps{ID: postID, AuthorID: "", AuthorUsername: ""}},
		{name: "an empty id is rejected", props: PostProps{ID: ""}, wantErrIs: x.ErrInvalidPostID},
		{name: "a malformed id is rejected", props: PostProps{ID: "nope"}, wantErrIs: x.ErrInvalidPostID},
		{
			name:      "a malformed author id is rejected",
			props:     PostProps{ID: postID, AuthorID: "nope"},
			wantErrIs: x.ErrInvalidUserID,
		},
		{
			name:      "a malformed author username is rejected",
			props:     PostProps{ID: postID, AuthorUsername: "ja ck"},
			wantErrIs: x.ErrInvalidUsername,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			post, err := NewPost(c.props)

			if c.wantErrIs != nil {
				require.ErrorIs(t, err, c.wantErrIs)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, c.props.ID, post.ID().String())
		})
	}
}

func TestPostView(t *testing.T) {
	t.Run("a view mirrors the props and carries a permalink", func(t *testing.T) {
		post, err := NewPost(PostProps{
			ID:             postID,
			Text:           "hello",
			AuthorUsername: "jack",
			Lang:           "en",
			Metrics:        Metrics{Likes: 3, Replies: 1},
		})
		require.NoError(t, err)

		view := post.View()

		assert.Equal(t, postID, view.ID)
		assert.Equal(t, "hello", view.Text)
		assert.Equal(t, "jack", view.AuthorUsername)
		assert.Equal(t, "en", view.Lang)
		assert.Equal(t, Metrics{Likes: 3, Replies: 1}, view.Metrics)
		assert.Equal(t, "https://x.com/jack/status/"+postID, view.URL)
	})

	t.Run("an unknown author falls back to the anonymous permalink", func(t *testing.T) {
		post, err := NewPost(PostProps{ID: postID, Text: "hello"})
		require.NoError(t, err)

		assert.Equal(t, "https://x.com/i/status/"+postID, post.View().URL)
	})

	t.Run("a view does not expose the internal slices", func(t *testing.T) {
		post, err := NewPost(PostProps{
			ID:        postID,
			URLs:      []string{"https://example.com"},
			MediaKeys: []string{"3_1"},
		})
		require.NoError(t, err)

		view := post.View()
		view.URLs[0] = "mutated"
		view.MediaKeys[0] = "mutated"

		assert.Equal(t, []string{"https://example.com"}, post.View().URLs)
		assert.Equal(t, []string{"3_1"}, post.View().MediaKeys)
	})

	t.Run("props slices are copied on construction", func(t *testing.T) {
		urls := []string{"https://example.com"}

		post, err := NewPost(PostProps{ID: postID, URLs: urls})
		require.NoError(t, err)

		urls[0] = "mutated"

		assert.Equal(t, []string{"https://example.com"}, post.View().URLs)
	})
}

func TestPostsLenAndIsEmpty(t *testing.T) {
	cases := []struct {
		name        string
		posts       Posts
		wantLen     int
		wantIsEmpty bool
	}{
		{name: "a nil collection", posts: nil, wantLen: 0, wantIsEmpty: true},
		{name: "an empty collection", posts: Posts{}, wantLen: 0, wantIsEmpty: true},
		{name: "one post", posts: Posts{newTestPost(t, postID)}, wantLen: 1},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			assert.Equal(t, c.wantLen, c.posts.Len())
			assert.Equal(t, c.wantIsEmpty, c.posts.IsEmpty())
		})
	}
}

func TestPostsLimit(t *testing.T) {
	all := Posts{newTestPost(t, "1"), newTestPost(t, "2"), newTestPost(t, "3")}

	cases := []struct {
		name    string
		posts   Posts
		n       int
		wantLen int
	}{
		{name: "fewer than available", posts: all, n: 2, wantLen: 2},
		{name: "exactly as many as available", posts: all, n: 3, wantLen: 3},
		{name: "more than available", posts: all, n: 10, wantLen: 3},
		{name: "zero returns nothing", posts: all, n: 0, wantLen: 0},
		{name: "a negative limit returns nothing", posts: all, n: -1, wantLen: 0},
		{name: "an empty collection stays empty", posts: nil, n: 5, wantLen: 0},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			assert.Equal(t, c.wantLen, c.posts.Limit(c.n).Len())
		})
	}
}

func TestPostsLimitKeepsTheLeadingPosts(t *testing.T) {
	limited := Posts{newTestPost(t, "1"), newTestPost(t, "2"), newTestPost(t, "3")}.Limit(2)

	require.Equal(t, 2, limited.Len())
	assert.Equal(t, "1", limited[0].ID().String())
	assert.Equal(t, "2", limited[1].ID().String())
}

func TestPostsDedupe(t *testing.T) {
	first := newTestPost(t, "1")
	second := newTestPost(t, "2")

	cases := []struct {
		name    string
		posts   Posts
		wantIDs []string
	}{
		{name: "a repeated id collapses to the first occurrence", posts: Posts{first, second, first}, wantIDs: []string{"1", "2"}},
		{name: "distinct posts are untouched", posts: Posts{first, second}, wantIDs: []string{"1", "2"}},
		{name: "an empty collection stays empty", posts: nil, wantIDs: []string{}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			deduped := c.posts.Dedupe()

			ids := make([]string, 0, deduped.Len())
			for _, post := range deduped {
				ids = append(ids, post.ID().String())
			}

			assert.Equal(t, c.wantIDs, ids)
		})
	}
}

func TestPostsViews(t *testing.T) {
	views := Posts{newTestPost(t, "1"), newTestPost(t, "2")}.Views()

	require.Len(t, views, 2)
	assert.Equal(t, "1", views[0].ID)
	assert.Equal(t, "2", views[1].ID)
	assert.Empty(t, Posts(nil).Views())
}

func TestPostsMarkdown(t *testing.T) {
	t.Run("an empty collection renders nothing", func(t *testing.T) {
		assert.Empty(t, Posts(nil).Markdown())
	})

	t.Run("every post is represented", func(t *testing.T) {
		markdown := Posts{newTestPost(t, "1"), newTestPost(t, "2")}.Markdown()

		assert.Contains(t, markdown, "text 1")
		assert.Contains(t, markdown, "text 2")
	})

	t.Run("a single post carries its text and its link", func(t *testing.T) {
		markdown := newTestPost(t, postID).Markdown()

		assert.Contains(t, markdown, "text "+postID)
		assert.True(t, strings.Contains(markdown, postID))
	})
}
