package posts

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Role1776/mcp-x/app/internal/domain/x"
)

func TestNewAuthorID(t *testing.T) {
	cases := []struct {
		name      string
		raw       string
		want      string
		wantErrIs error
	}{
		{name: "a snowflake", raw: "1234567890123456789", want: "1234567890123456789"},
		{name: "an empty string is rejected", raw: "", wantErrIs: x.ErrInvalidUserID},
		{name: "a blank string is rejected", raw: "  ", wantErrIs: x.ErrInvalidUserID},
		{name: "letters are rejected", raw: "abc", wantErrIs: x.ErrInvalidUserID},
		{name: "a negative number is rejected", raw: "-7", wantErrIs: x.ErrInvalidUserID},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			id, err := NewAuthorID(c.raw)

			if c.wantErrIs != nil {
				require.ErrorIs(t, err, c.wantErrIs)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, c.want, id.String())
		})
	}
}

func TestNewAuthorUsername(t *testing.T) {
	cases := []struct {
		name      string
		raw       string
		want      string
		wantErrIs error
	}{
		{name: "a handle", raw: "jack", want: "jack"},
		{name: "case is preserved", raw: "Jack", want: "Jack"},
		{name: "an empty string is rejected", raw: "", wantErrIs: x.ErrInvalidUsername},
		{name: "an inner space is rejected", raw: "ja ck", wantErrIs: x.ErrInvalidUsername},
		{name: "above the handle limit is rejected", raw: "aaaaaaaaaaaaaaaa", wantErrIs: x.ErrInvalidUsername},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			name, err := NewAuthorUsername(c.raw)

			if c.wantErrIs != nil {
				require.ErrorIs(t, err, c.wantErrIs)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, c.want, name.String())
		})
	}
}
