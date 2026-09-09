package lists

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Role1776/mcp-x/app/internal/domain/x"
)

func TestNewListID(t *testing.T) {
	cases := []struct {
		name      string
		raw       string
		want      string
		wantErrIs error
	}{
		{name: "a snowflake", raw: "1234567890123456789", want: "1234567890123456789"},
		{name: "surrounding whitespace is trimmed", raw: " 12 ", want: "12"},
		{name: "empty string is rejected", raw: "", wantErrIs: x.ErrInvalidListID},
		{name: "blank string is rejected", raw: "  ", wantErrIs: x.ErrInvalidListID},
		{name: "letters are rejected", raw: "list", wantErrIs: x.ErrInvalidListID},
		{name: "a trailing letter is rejected", raw: "12a", wantErrIs: x.ErrInvalidListID},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			id, err := NewListID(c.raw)

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

func TestNewOwnerID(t *testing.T) {
	cases := []struct {
		name      string
		raw       string
		want      string
		wantErrIs error
	}{
		{name: "a snowflake", raw: "1234567890123456789", want: "1234567890123456789"},
		{name: "an empty string is rejected", raw: "", wantErrIs: x.ErrInvalidUserID},
		{name: "a blank string is rejected", raw: "  ", wantErrIs: x.ErrInvalidUserID},
		{name: "letters are rejected", raw: "jack", wantErrIs: x.ErrInvalidUserID},
		{name: "a negative number is rejected", raw: "-1", wantErrIs: x.ErrInvalidUserID},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			id, err := NewOwnerID(c.raw)

			if c.wantErrIs != nil {
				require.ErrorIs(t, err, c.wantErrIs)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, c.want, id.String())
		})
	}
}
