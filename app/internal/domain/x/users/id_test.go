package users

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Role1776/mcp-x/app/internal/domain/x"
)

func TestNewUserID(t *testing.T) {
	cases := []struct {
		name      string
		raw       string
		want      string
		wantErrIs error
	}{
		{name: "a snowflake", raw: "1234567890123456789", want: "1234567890123456789"},
		{name: "a short id", raw: "12", want: "12"},
		{name: "surrounding whitespace is trimmed", raw: " 12 ", want: "12"},
		{name: "empty string is rejected", raw: "", wantErrIs: x.ErrInvalidUserID},
		{name: "blank string is rejected", raw: "   ", wantErrIs: x.ErrInvalidUserID},
		{name: "a handle is not an id", raw: "jack", wantErrIs: x.ErrInvalidUserID},
		{name: "a trailing letter is rejected", raw: "1a", wantErrIs: x.ErrInvalidUserID},
		{name: "a negative number is rejected", raw: "-1", wantErrIs: x.ErrInvalidUserID},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			id, err := NewUserID(c.raw)

			if c.wantErrIs != nil {
				require.ErrorIs(t, err, c.wantErrIs)
				assert.Empty(t, id.String())

				return
			}

			require.NoError(t, err)
			assert.Equal(t, c.want, id.String())
		})
	}
}

func TestNewUserIDs(t *testing.T) {
	cases := []struct {
		name      string
		raw       []string
		want      []string
		wantErrIs error
	}{
		{name: "several ids keep their order", raw: []string{"3", "1"}, want: []string{"3", "1"}},
		{name: "a nil slice is rejected", raw: nil, wantErrIs: x.ErrInvalidRequest},
		{name: "an empty slice is rejected", raw: []string{}, wantErrIs: x.ErrInvalidRequest},
		{name: "one bad id fails the whole call", raw: []string{"1", "jack"}, wantErrIs: x.ErrInvalidUserID},
		{name: "an empty element fails the whole call", raw: []string{"1", ""}, wantErrIs: x.ErrInvalidUserID},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ids, err := NewUserIDs(c.raw)

			if c.wantErrIs != nil {
				require.ErrorIs(t, err, c.wantErrIs)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, len(c.want), ids.Len())
			assert.Equal(t, c.want, ids.Strings())
		})
	}
}

func TestUserIDsDedupe(t *testing.T) {
	cases := []struct {
		name string
		raw  []string
		want []string
	}{
		{name: "duplicates collapse to the first occurrence", raw: []string{"1", "2", "1"}, want: []string{"1", "2"}},
		{name: "an already unique slice is unchanged", raw: []string{"1", "2"}, want: []string{"1", "2"}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ids, err := NewUserIDs(c.raw)
			require.NoError(t, err)

			assert.Equal(t, c.want, ids.Dedupe().Strings())
		})
	}
}

func TestUserIDsStringsDoesNotShareState(t *testing.T) {
	ids, err := NewUserIDs([]string{"1", "2"})
	require.NoError(t, err)

	first := ids.Strings()
	first[0] = "999"

	assert.Equal(t, []string{"1", "2"}, ids.Strings())
}
