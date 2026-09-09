package users

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Role1776/mcp-x/app/internal/domain/x"
)

func TestNewUsername(t *testing.T) {
	cases := []struct {
		name           string
		raw            string
		want           string
		wantNormalized string
		wantErrIs      error
	}{
		{name: "a plain handle", raw: "jack", want: "jack", wantNormalized: "jack"},
		{name: "case is preserved but normalized separately", raw: "JaCk", want: "JaCk", wantNormalized: "jack"},
		{name: "a leading at sign is stripped", raw: "@jack", want: "jack", wantNormalized: "jack"},
		{name: "surrounding whitespace is trimmed", raw: "  jack  ", want: "jack", wantNormalized: "jack"},
		{name: "underscores and digits are allowed", raw: "j_1", want: "j_1", wantNormalized: "j_1"},
		{
			name:           "exactly at the handle limit",
			raw:            strings.Repeat("a", 15),
			want:           strings.Repeat("a", 15),
			wantNormalized: strings.Repeat("a", 15),
		},
		{name: "one above the handle limit is rejected", raw: strings.Repeat("a", 16), wantErrIs: x.ErrInvalidUsername},
		{name: "empty string is rejected", raw: "", wantErrIs: x.ErrInvalidUsername},
		{name: "blank string is rejected", raw: "   ", wantErrIs: x.ErrInvalidUsername},
		{name: "a bare at sign is rejected", raw: "@", wantErrIs: x.ErrInvalidUsername},
		{name: "an inner space is rejected", raw: "ja ck", wantErrIs: x.ErrInvalidUsername},
		{name: "a hyphen is rejected", raw: "ja-ck", wantErrIs: x.ErrInvalidUsername},
		{name: "non-latin letters are rejected", raw: "джек", wantErrIs: x.ErrInvalidUsername},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			name, err := NewUsername(c.raw)

			if c.wantErrIs != nil {
				require.ErrorIs(t, err, c.wantErrIs)
				assert.Empty(t, name.String())

				return
			}

			require.NoError(t, err)
			assert.Equal(t, c.want, name.String())
			assert.Equal(t, c.wantNormalized, name.Normalized())
		})
	}
}

func TestNewUsernames(t *testing.T) {
	cases := []struct {
		name      string
		raw       []string
		want      []string
		wantErrIs error
	}{
		{name: "several handles keep their order", raw: []string{"b", "a"}, want: []string{"b", "a"}},
		{name: "the original case survives", raw: []string{"JaCk"}, want: []string{"JaCk"}},
		{name: "a nil slice is rejected", raw: nil, wantErrIs: x.ErrInvalidRequest},
		{name: "an empty slice is rejected", raw: []string{}, wantErrIs: x.ErrInvalidRequest},
		{name: "one bad handle fails the whole call", raw: []string{"jack", "ja ck"}, wantErrIs: x.ErrInvalidUsername},
		{name: "an empty element fails the whole call", raw: []string{"jack", ""}, wantErrIs: x.ErrInvalidUsername},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			names, err := NewUsernames(c.raw)

			if c.wantErrIs != nil {
				require.ErrorIs(t, err, c.wantErrIs)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, len(c.want), names.Len())
			assert.Equal(t, c.want, names.Strings())
		})
	}
}

func TestUsernamesDedupe(t *testing.T) {
	cases := []struct {
		name string
		raw  []string
		want []string
	}{
		{name: "identical handles collapse", raw: []string{"jack", "jack"}, want: []string{"jack"}},
		{
			name: "handles differing only in case are the same account",
			raw:  []string{"jack", "JACK", "Jack"},
			want: []string{"jack"},
		},
		{name: "distinct handles survive", raw: []string{"jack", "jill"}, want: []string{"jack", "jill"}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			names, err := NewUsernames(c.raw)
			require.NoError(t, err)

			assert.Equal(t, c.want, names.Dedupe().Strings())
		})
	}
}

func TestUsernamesStringsDoesNotShareState(t *testing.T) {
	names, err := NewUsernames([]string{"jack", "jill"})
	require.NoError(t, err)

	first := names.Strings()
	first[0] = "mutated"

	assert.Equal(t, []string{"jack", "jill"}, names.Strings())
}
