package lists

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Role1776/mcp-x/app/internal/domain/x"
)

func ptr[T any](value T) *T {
	return &value
}

func TestNewDraft(t *testing.T) {
	cases := []struct {
		name            string
		props           DraftProps
		wantName        string
		wantDescription string
		wantPrivate     bool
		wantErrIs       error
	}{
		{name: "a minimal draft", props: DraftProps{Name: "reading"}, wantName: "reading"},
		{
			name:            "name and description are trimmed",
			props:           DraftProps{Name: "  reading  ", Description: "  books  "},
			wantName:        "reading",
			wantDescription: "books",
		},
		{
			name:            "a whitespace-only description becomes empty",
			props:           DraftProps{Name: "reading", Description: "   "},
			wantName:        "reading",
			wantDescription: "",
		},
		{
			name:        "the private flag is carried through",
			props:       DraftProps{Name: "reading", Private: true},
			wantName:    "reading",
			wantPrivate: true,
		},
		{
			name:     "a name one below the limit",
			props:    DraftProps{Name: strings.Repeat("a", MaxNameLength-1)},
			wantName: strings.Repeat("a", MaxNameLength-1),
		},
		{
			name:     "a name exactly at the limit",
			props:    DraftProps{Name: strings.Repeat("a", MaxNameLength)},
			wantName: strings.Repeat("a", MaxNameLength),
		},
		{
			name:     "the name limit counts runes and not bytes",
			props:    DraftProps{Name: strings.Repeat("🙂", MaxNameLength)},
			wantName: strings.Repeat("🙂", MaxNameLength),
		},
		{
			name:            "a description exactly at the limit",
			props:           DraftProps{Name: "reading", Description: strings.Repeat("a", MaxDescriptionLength)},
			wantName:        "reading",
			wantDescription: strings.Repeat("a", MaxDescriptionLength),
		},
		{
			name:      "a name above the limit is rejected",
			props:     DraftProps{Name: strings.Repeat("a", MaxNameLength+1)},
			wantErrIs: x.ErrListNameTooLong,
		},
		{
			name:      "a description above the limit is rejected",
			props:     DraftProps{Name: "reading", Description: strings.Repeat("a", MaxDescriptionLength+1)},
			wantErrIs: x.ErrListDescriptionTooLong,
		},
		{name: "an empty name is rejected", props: DraftProps{Name: ""}, wantErrIs: x.ErrInvalidRequest},
		{name: "a whitespace-only name is rejected", props: DraftProps{Name: "   "}, wantErrIs: x.ErrInvalidRequest},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			draft, err := NewDraft(c.props)

			if c.wantErrIs != nil {
				require.ErrorIs(t, err, c.wantErrIs)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, c.wantName, draft.Name())
			assert.Equal(t, c.wantDescription, draft.Description())
			assert.Equal(t, c.wantPrivate, draft.Private())
		})
	}
}

func TestNewPatch(t *testing.T) {
	cases := []struct {
		name      string
		props     PatchProps
		wantErrIs error
	}{
		{name: "an empty patch is rejected", props: PatchProps{}, wantErrIs: x.ErrInvalidRequest},
		{name: "only the name", props: PatchProps{Name: ptr("reading")}},
		{name: "only the description", props: PatchProps{Description: ptr("books")}},
		{name: "only the privacy flag", props: PatchProps{Private: ptr(true)}},
		{name: "every field at once", props: PatchProps{
			Name: ptr("reading"), Description: ptr("books"), Private: ptr(false),
		}},
		{name: "a name at the limit", props: PatchProps{Name: ptr(strings.Repeat("a", MaxNameLength))}},
		{
			name:      "a name above the limit is rejected",
			props:     PatchProps{Name: ptr(strings.Repeat("a", MaxNameLength+1))},
			wantErrIs: x.ErrListNameTooLong,
		},
		{
			name:      "a description above the limit is rejected",
			props:     PatchProps{Description: ptr(strings.Repeat("a", MaxDescriptionLength+1))},
			wantErrIs: x.ErrListDescriptionTooLong,
		},
		{name: "clearing the name is rejected", props: PatchProps{Name: ptr("")}, wantErrIs: x.ErrInvalidRequest},
		{name: "a whitespace-only name is rejected", props: PatchProps{Name: ptr("  ")}, wantErrIs: x.ErrInvalidRequest},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, err := NewPatch(c.props)

			if c.wantErrIs != nil {
				require.ErrorIs(t, err, c.wantErrIs)

				return
			}

			require.NoError(t, err)
		})
	}
}

func TestPatchGetters(t *testing.T) {
	t.Run("unset fields stay nil", func(t *testing.T) {
		patch, err := NewPatch(PatchProps{Name: ptr("reading")})
		require.NoError(t, err)

		require.NotNil(t, patch.Name())
		assert.Equal(t, "reading", *patch.Name())
		assert.Nil(t, patch.Description())
		assert.Nil(t, patch.Private())
	})

	t.Run("an empty description clears the field rather than skipping it", func(t *testing.T) {
		patch, err := NewPatch(PatchProps{Description: ptr("")})
		require.NoError(t, err)

		require.NotNil(t, patch.Description())
		assert.Empty(t, *patch.Description())
	})

	t.Run("a false privacy flag is still a change", func(t *testing.T) {
		patch, err := NewPatch(PatchProps{Private: ptr(false)})
		require.NoError(t, err)

		require.NotNil(t, patch.Private())
		assert.False(t, *patch.Private())
	})

	t.Run("the patch does not alias the caller's pointers", func(t *testing.T) {
		name := "reading"

		patch, err := NewPatch(PatchProps{Name: &name})
		require.NoError(t, err)

		name = "mutated"

		require.NotNil(t, patch.Name())
		assert.Equal(t, "reading", *patch.Name())
	})
}
