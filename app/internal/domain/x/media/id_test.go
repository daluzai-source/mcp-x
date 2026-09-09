package media

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Role1776/mcp-x/app/internal/domain/x"
)

func TestNewMediaID(t *testing.T) {
	cases := []struct {
		name      string
		raw       string
		want      string
		wantErrIs error
	}{
		{name: "a snowflake", raw: "1146654567674912769", want: "1146654567674912769"},
		{name: "surrounding whitespace is trimmed", raw: " 12 ", want: "12"},
		{name: "empty string is rejected", raw: "", wantErrIs: x.ErrInvalidMediaID},
		{name: "blank string is rejected", raw: "  ", wantErrIs: x.ErrInvalidMediaID},
		{name: "letters are rejected", raw: "media", wantErrIs: x.ErrInvalidMediaID},
		{name: "a media key is not a media id", raw: "3_1146654567674912769", wantErrIs: x.ErrInvalidMediaID},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			id, err := NewMediaID(c.raw)

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
