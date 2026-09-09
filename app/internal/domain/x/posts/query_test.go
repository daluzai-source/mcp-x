package posts

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Role1776/mcp-x/app/internal/domain/x"
)

func TestNewQuery(t *testing.T) {
	cases := []struct {
		name      string
		raw       string
		want      string
		wantErrIs error
	}{
		{name: "a plain term", raw: "golang", want: "golang"},
		{name: "operators are kept verbatim", raw: `from:jack "exact phrase" -is:retweet`, want: `from:jack "exact phrase" -is:retweet`},
		{name: "surrounding whitespace is trimmed", raw: "  golang  ", want: "golang"},
		{name: "one below the limit", raw: strings.Repeat("a", maxQueryLength-1), want: strings.Repeat("a", maxQueryLength-1)},
		{name: "exactly at the limit", raw: strings.Repeat("a", maxQueryLength), want: strings.Repeat("a", maxQueryLength)},
		{name: "one above the limit is rejected", raw: strings.Repeat("a", maxQueryLength+1), wantErrIs: x.ErrQueryTooLong},
		{name: "the limit counts runes and not bytes", raw: strings.Repeat("🙂", maxQueryLength), want: strings.Repeat("🙂", maxQueryLength)},
		{name: "empty string is rejected", raw: "", wantErrIs: x.ErrEmptyQuery},
		{name: "whitespace only is rejected", raw: "   ", wantErrIs: x.ErrEmptyQuery},
		{name: "tabs and newlines only are rejected", raw: "\t\n", wantErrIs: x.ErrEmptyQuery},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			query, err := NewQuery(c.raw)

			if c.wantErrIs != nil {
				require.ErrorIs(t, err, c.wantErrIs)
				assert.Empty(t, query.String())

				return
			}

			require.NoError(t, err)
			assert.Equal(t, c.want, query.String())
		})
	}
}

func TestQueryWithoutRetweets(t *testing.T) {
	t.Run("appends the operator", func(t *testing.T) {
		query, err := NewQuery("golang")
		require.NoError(t, err)

		filtered, err := query.WithoutRetweets()
		require.NoError(t, err)

		assert.Contains(t, filtered.String(), "-is:retweet")
		assert.Contains(t, filtered.String(), "golang")
	})

	t.Run("leaves the receiver untouched", func(t *testing.T) {
		query, err := NewQuery("golang")
		require.NoError(t, err)

		_, err = query.WithoutRetweets()
		require.NoError(t, err)

		assert.Equal(t, "golang", query.String())
	})

	t.Run("a query at the length limit no longer fits the operator", func(t *testing.T) {
		query, err := NewQuery(strings.Repeat("a", maxQueryLength))
		require.NoError(t, err)

		_, err = query.WithoutRetweets()

		require.ErrorIs(t, err, x.ErrQueryTooLong)
	})
}

func TestNewQueries(t *testing.T) {
	cases := []struct {
		name      string
		raw       []string
		want      []string
		wantErrIs error
	}{
		{name: "several queries keep their order", raw: []string{"b", "a"}, want: []string{"b", "a"}},
		{name: "a nil slice is rejected", raw: nil, wantErrIs: x.ErrEmptyQuery},
		{name: "an empty slice is rejected", raw: []string{}, wantErrIs: x.ErrEmptyQuery},
		{name: "an empty element fails the whole call", raw: []string{"ok", ""}, wantErrIs: x.ErrEmptyQuery},
		{name: "a blank element fails the whole call", raw: []string{"ok", "  "}, wantErrIs: x.ErrEmptyQuery},
		{
			name:      "an over-long element fails the whole call",
			raw:       []string{"ok", strings.Repeat("a", maxQueryLength+1)},
			wantErrIs: x.ErrQueryTooLong,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			queries, err := NewQueries(c.raw)

			if c.wantErrIs != nil {
				require.ErrorIs(t, err, c.wantErrIs)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, len(c.want), queries.Len())
			assert.Equal(t, c.want, queries.Strings())
		})
	}
}

func TestQueriesWithoutRetweets(t *testing.T) {
	t.Run("every query gets the operator", func(t *testing.T) {
		queries, err := NewQueries([]string{"go", "rust"})
		require.NoError(t, err)

		filtered, err := queries.WithoutRetweets()
		require.NoError(t, err)

		require.Equal(t, 2, filtered.Len())
		for _, value := range filtered.Strings() {
			assert.Contains(t, value, "-is:retweet")
		}
	})

	t.Run("one query that no longer fits fails the whole set", func(t *testing.T) {
		queries, err := NewQueries([]string{"go", strings.Repeat("a", maxQueryLength)})
		require.NoError(t, err)

		_, err = queries.WithoutRetweets()

		require.ErrorIs(t, err, x.ErrQueryTooLong)
	})
}

func TestQueriesStringsDoesNotShareState(t *testing.T) {
	queries, err := NewQueries([]string{"go", "rust"})
	require.NoError(t, err)

	first := queries.Strings()
	first[0] = "mutated"

	assert.Equal(t, []string{"go", "rust"}, queries.Strings())
}
