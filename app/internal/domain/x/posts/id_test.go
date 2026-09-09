package posts

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Role1776/mcp-x/app/internal/domain/x"
)

func TestNewPostID(t *testing.T) {
	cases := []struct {
		name    string
		raw     string
		want    string
		wantErr error
	}{
		{name: "a plain snowflake", raw: "1234567890123456789", want: "1234567890123456789"},
		{name: "a short id", raw: "1", want: "1"},
		{name: "surrounding whitespace is trimmed", raw: "  123  ", want: "123"},
		{name: "empty string is rejected", raw: "", wantErr: x.ErrInvalidPostID},
		{name: "blank string is rejected", raw: "   ", wantErr: x.ErrInvalidPostID},
		{name: "letters are rejected", raw: "abc", wantErr: x.ErrInvalidPostID},
		{name: "a trailing letter is rejected", raw: "12a", wantErr: x.ErrInvalidPostID},
		{name: "a negative number is rejected", raw: "-1", wantErr: x.ErrInvalidPostID},
		{name: "a decimal is rejected", raw: "1.5", wantErr: x.ErrInvalidPostID},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			id, err := NewPostID(c.raw)

			if c.wantErr != nil {
				require.ErrorIs(t, err, c.wantErr)
				assert.True(t, id.IsZero())

				return
			}

			require.NoError(t, err)
			assert.Equal(t, c.want, id.String())
			assert.False(t, id.IsZero())
		})
	}
}

func TestPostIDZeroValue(t *testing.T) {
	var id PostID

	assert.True(t, id.IsZero())
	assert.Empty(t, id.String())
}

func TestNewPostIDs(t *testing.T) {
	cases := []struct {
		name    string
		raw     []string
		want    []string
		wantErr error
	}{
		{name: "several ids keep their order", raw: []string{"3", "1", "2"}, want: []string{"3", "1", "2"}},
		{name: "duplicates are collapsed by the constructor", raw: []string{"1", "1"}, want: []string{"1"}},
		{name: "a nil slice is rejected", raw: nil, wantErr: x.ErrInvalidRequest},
		{name: "an empty slice is rejected", raw: []string{}, wantErr: x.ErrInvalidRequest},
		{name: "one bad id fails the whole call", raw: []string{"1", "nope", "3"}, wantErr: x.ErrInvalidPostID},
		{name: "an empty element fails the whole call", raw: []string{"1", ""}, wantErr: x.ErrInvalidPostID},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ids, err := NewPostIDs(c.raw)

			if c.wantErr != nil {
				require.ErrorIs(t, err, c.wantErr)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, len(c.want), ids.Len())
			assert.Equal(t, c.want, ids.Strings())
		})
	}
}

func TestPostIDsDedupe(t *testing.T) {
	cases := []struct {
		name string
		raw  []string
		want []string
	}{
		{name: "duplicates collapse to the first occurrence", raw: []string{"1", "2", "1", "3"}, want: []string{"1", "2", "3"}},
		{name: "an already unique slice is unchanged", raw: []string{"1", "2"}, want: []string{"1", "2"}},
		{name: "every element identical", raw: []string{"7", "7", "7"}, want: []string{"7"}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ids, err := NewPostIDs(c.raw)
			require.NoError(t, err)

			assert.Equal(t, c.want, ids.Dedupe().Strings())
		})
	}
}

func TestPostIDsStringsDoesNotShareState(t *testing.T) {
	ids, err := NewPostIDs([]string{"1", "2"})
	require.NoError(t, err)

	first := ids.Strings()
	first[0] = "999"

	assert.Equal(t, []string{"1", "2"}, ids.Strings())
}
