package posts

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/michimani/gotwi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Role1776/mcp-x/app/internal/domain/x"
	"github.com/Role1776/mcp-x/app/internal/domain/x/posts"
	"github.com/Role1776/mcp-x/app/internal/domain/x/users"
)

func newAdapter(t *testing.T, status int, fixture string) *Adapter {
	t.Helper()

	body, err := os.ReadFile("testdata/" + fixture)
	require.NoError(t, err)

	client, err := gotwi.NewClient(&gotwi.NewClientInput{
		HTTPClient: gotwi.NewMockHTTPClient(&gotwi.MockInput{
			ResponseStatusCode: status,
			ResponseHeader:     map[string][]string{"Content-Type": {"application/json"}},
			ResponseBody:       io.NopCloser(bytes.NewReader(body)),
		}),
		AuthenticationMethod: gotwi.AuthenMethodOAuth1UserContext,
		APIKey:               "key",
		APIKeySecret:         "secret",
		OAuthToken:           "token",
		OAuthTokenSecret:     "token-secret",
	})
	require.NoError(t, err)

	selfID, err := users.NewUserID("44196397")
	require.NoError(t, err)

	return New(client, selfID)
}

func TestSearchMapping(t *testing.T) {
	adapter := newAdapter(t, http.StatusOK, "search.json")

	query, err := posts.NewQuery("golang")
	require.NoError(t, err)

	found, nextToken, err := adapter.Search(context.Background(), SearchParams{Query: query, MaxResults: 10})
	require.NoError(t, err)

	require.Equal(t, 3, found.Len())
	assert.Equal(t, "b26v89c19zqg8o3fpds1", nextToken)

	views := found.Views()

	t.Run("joins the author from includes", func(t *testing.T) {
		assert.Equal(t, "44196397", views[0].AuthorID)
		assert.Equal(t, "elonmusk", views[0].AuthorUsername)
		assert.Equal(t, "https://x.com/elonmusk/status/1800000000000000001", views[0].URL)
	})

	t.Run("dereferences metrics, urls and media keys", func(t *testing.T) {
		assert.Equal(t, posts.Metrics{Retweets: 12, Replies: 3, Likes: 100, Quotes: 1}, views[0].Metrics)
		assert.Equal(t, []string{"https://go.dev/blog"}, views[0].URLs)
		assert.Equal(t, []string{"3_1800000000000000009"}, views[0].MediaKeys)
		assert.Equal(t, []posts.Reference{{Type: "replied_to", ID: "1799999999999999999"}}, views[0].Referenced)
		assert.Equal(t, "en", views[0].Lang)
		assert.False(t, views[0].CreatedAt.IsZero())
	})

	t.Run("prefers note_tweet over the truncated text", func(t *testing.T) {
		assert.Equal(t, "The full long text of the note tweet, well past the 280 character cut.", views[1].Text)
	})

	t.Run("keeps a post whose author is missing from includes", func(t *testing.T) {
		assert.Equal(t, "99999999", views[1].AuthorID)
		assert.Empty(t, views[1].AuthorUsername)
	})

	t.Run("survives a post with every optional field absent", func(t *testing.T) {
		assert.Equal(t, "1800000000000000003", views[2].ID)
		assert.Empty(t, views[2].AuthorID)
		assert.Equal(t, posts.Metrics{}, views[2].Metrics)
		assert.Empty(t, views[2].URLs)
		assert.True(t, views[2].CreatedAt.IsZero())
	})
}

func TestLookupPartialErrors(t *testing.T) {
	adapter := newAdapter(t, http.StatusOK, "lookup_partial.json")

	ids, err := posts.NewPostIDs([]string{"1800000000000000001", "1800000000000000002"})
	require.NoError(t, err)

	found, failures, err := adapter.Lookup(context.Background(), ids)
	require.NoError(t, err)

	assert.Equal(t, 1, found.Len())
	require.Len(t, failures, 1)
	assert.Equal(t, "1800000000000000002", failures[0].ResourceID)
	assert.Equal(t, "Not Found Error", failures[0].Title)
	assert.Contains(t, failures[0].Detail, "Could not find tweet")
}

func TestSearchMapsAPIErrors(t *testing.T) {
	adapter := newAdapter(t, http.StatusTooManyRequests, "search.json")

	query, err := posts.NewQuery("golang")
	require.NoError(t, err)

	_, _, err = adapter.Search(context.Background(), SearchParams{Query: query, MaxResults: 10})
	require.ErrorIs(t, err, x.ErrRateLimited)
}
