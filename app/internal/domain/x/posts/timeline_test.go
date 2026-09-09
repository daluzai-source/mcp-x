package posts

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Role1776/mcp-x/app/internal/domain/x"
)

func TestNewWindow(t *testing.T) {
	cases := []struct {
		name      string
		days      int
		wantErrIs error
	}{
		{name: "the smallest window", days: 1},
		{name: "one below the limit", days: MaxWindowDays - 1},
		{name: "exactly at the limit", days: MaxWindowDays},
		{name: "one above the limit is rejected", days: MaxWindowDays + 1, wantErrIs: x.ErrInvalidRequest},
		{name: "zero is rejected", days: 0, wantErrIs: x.ErrInvalidRequest},
		{name: "a negative window is rejected", days: -1, wantErrIs: x.ErrInvalidRequest},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			window, err := NewWindow(c.days)

			if c.wantErrIs != nil {
				require.ErrorIs(t, err, c.wantErrIs)
				assert.Equal(t, 0, window.Days())

				return
			}

			require.NoError(t, err)
			assert.Equal(t, c.days, window.Days())
		})
	}
}

func TestWindowStartTime(t *testing.T) {
	now := time.Date(2026, time.March, 10, 12, 0, 0, 0, time.UTC)

	cases := []struct {
		name string
		days int
		want time.Time
	}{
		{name: "one day back", days: 1, want: now.AddDate(0, 0, -1)},
		{name: "the full window back", days: MaxWindowDays, want: now.AddDate(0, 0, -MaxWindowDays)},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			window, err := NewWindow(c.days)
			require.NoError(t, err)

			start := window.StartTime(now)

			assert.Equal(t, c.want, start)
			assert.True(t, start.Before(now))
		})
	}
}

func TestNewExclude(t *testing.T) {
	cases := []struct {
		name         string
		values       []string
		wantReplies  bool
		wantRetweets bool
		wantStrings  []string
		wantErrIs    error
	}{
		{name: "nothing is excluded by default", values: nil, wantStrings: []string{}},
		{name: "an empty slice excludes nothing", values: []string{}, wantStrings: []string{}},
		{name: "replies only", values: []string{"replies"}, wantReplies: true, wantStrings: []string{"replies"}},
		{name: "retweets only", values: []string{"retweets"}, wantRetweets: true, wantStrings: []string{"retweets"}},
		{
			name:         "both kinds",
			values:       []string{"replies", "retweets"},
			wantReplies:  true,
			wantRetweets: true,
			wantStrings:  []string{"replies", "retweets"},
		},
		{
			name:         "order of the input does not matter",
			values:       []string{"retweets", "replies"},
			wantReplies:  true,
			wantRetweets: true,
			wantStrings:  []string{"replies", "retweets"},
		},
		{
			name:        "a repeated value is reported once",
			values:      []string{"replies", "replies"},
			wantReplies: true,
			wantStrings: []string{"replies"},
		},
		{name: "an unknown value is rejected", values: []string{"quotes"}, wantErrIs: x.ErrInvalidRequest},
		{name: "the enum is case sensitive", values: []string{"Replies"}, wantErrIs: x.ErrInvalidRequest},
		{name: "an empty value is rejected", values: []string{""}, wantErrIs: x.ErrInvalidRequest},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			exclude, err := NewExclude(c.values)

			if c.wantErrIs != nil {
				require.ErrorIs(t, err, c.wantErrIs)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, c.wantReplies, exclude.Replies())
			assert.Equal(t, c.wantRetweets, exclude.Retweets())
			assert.ElementsMatch(t, c.wantStrings, exclude.Strings())
		})
	}
}

func TestExcludeStringsDoesNotShareState(t *testing.T) {
	exclude, err := NewExclude([]string{"replies", "retweets"})
	require.NoError(t, err)

	first := exclude.Strings()
	first[0] = "mutated"

	assert.ElementsMatch(t, []string{"replies", "retweets"}, exclude.Strings())
}
