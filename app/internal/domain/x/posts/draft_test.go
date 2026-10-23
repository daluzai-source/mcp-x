package posts

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Role1776/mcp-x/app/internal/domain/x"
)

const (
	mediaID      = "1146654567674912769"
	otherMediaID = "1146654567674912770"
	postID       = "1234567890123456789"
	otherPostID  = "1234567890123456780"
)

func validPoll() *PollProps {
	return &PollProps{Options: []string{"yes", "no"}, DurationMinutes: 60}
}

func TestNewDraftRejects(t *testing.T) {
	cases := []struct {
		name      string
		props     DraftProps
		wantErrIs error
	}{
		{name: "neither text nor media", props: DraftProps{}, wantErrIs: x.ErrEmptyText},
		{name: "whitespace only text and no media", props: DraftProps{Text: "   "}, wantErrIs: x.ErrEmptyText},
		{
			name:      "text above the limit",
			props:     DraftProps{Text: strings.Repeat("a", MaxTextLength+1)},
			wantErrIs: x.ErrTextTooLong,
		},
		{
			name:      "note text above the note limit",
			props:     DraftProps{Text: strings.Repeat("a", MaxNoteTextLength+1), Note: true},
			wantErrIs: x.ErrTextTooLong,
		},
		{
			name:      "more media than a post can carry",
			props:     DraftProps{Text: "hi", MediaIDs: []string{mediaID, mediaID, mediaID, mediaID, mediaID}},
			wantErrIs: x.ErrTooManyMedia,
		},
		{
			name:      "a media id that is not a snowflake",
			props:     DraftProps{Text: "hi", MediaIDs: []string{"not-an-id"}},
			wantErrIs: x.ErrInvalidMediaID,
		},
		{
			name:      "an empty media id",
			props:     DraftProps{Text: "hi", MediaIDs: []string{""}},
			wantErrIs: x.ErrInvalidMediaID,
		},
		{
			name:      "a reply target that is not a snowflake",
			props:     DraftProps{Text: "hi", ReplyToID: "nope"},
			wantErrIs: x.ErrInvalidPostID,
		},
		{
			name:      "a quote target that is not a snowflake",
			props:     DraftProps{Text: "hi", QuoteID: "nope"},
			wantErrIs: x.ErrInvalidPostID,
		},
		{
			name:      "a poll with too few options",
			props:     DraftProps{Text: "hi", Poll: &PollProps{Options: []string{"only"}, DurationMinutes: 60}},
			wantErrIs: x.ErrInvalidPoll,
		},
		{
			name: "a poll with too many options",
			props: DraftProps{Text: "hi", Poll: &PollProps{
				Options: []string{"a", "b", "c", "d", "e"}, DurationMinutes: 60,
			}},
			wantErrIs: x.ErrInvalidPoll,
		},
		{
			name:      "a poll with a blank option",
			props:     DraftProps{Text: "hi", Poll: &PollProps{Options: []string{"a", "  "}, DurationMinutes: 60}},
			wantErrIs: x.ErrInvalidPoll,
		},
		{
			name:      "a poll shorter than the minimum duration",
			props:     DraftProps{Text: "hi", Poll: &PollProps{Options: []string{"a", "b"}, DurationMinutes: 4}},
			wantErrIs: x.ErrInvalidPoll,
		},
		{
			name:      "a poll longer than the maximum duration",
			props:     DraftProps{Text: "hi", Poll: &PollProps{Options: []string{"a", "b"}, DurationMinutes: 10081}},
			wantErrIs: x.ErrInvalidPoll,
		},
		{
			name:      "a poll combined with media",
			props:     DraftProps{Text: "hi", MediaIDs: []string{mediaID}, Poll: validPoll()},
			wantErrIs: x.ErrInvalidRequest,
		},
		{
			name:      "reply settings outside the enum",
			props:     DraftProps{Text: "hi", ReplySettings: "everyone"},
			wantErrIs: x.ErrInvalidReplySettings,
		},
		{
			name:      "reply settings in the wrong case",
			props:     DraftProps{Text: "hi", ReplySettings: "Following"},
			wantErrIs: x.ErrInvalidReplySettings,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, err := NewDraft(c.props)

			require.ErrorIs(t, err, c.wantErrIs)
		})
	}
}

func TestNewDraftAccepts(t *testing.T) {
	cases := []struct {
		name  string
		props DraftProps
	}{
		{name: "text only", props: DraftProps{Text: "hello"}},
		{name: "text one below the limit", props: DraftProps{Text: strings.Repeat("a", MaxTextLength-1)}},
		{name: "text exactly at the limit", props: DraftProps{Text: strings.Repeat("a", MaxTextLength)}},
		{
			name:  "the limit counts runes and not bytes",
			props: DraftProps{Text: strings.Repeat("🙂", MaxTextLength)},
		},
		{
			name:  "a note post exactly at the note limit",
			props: DraftProps{Text: strings.Repeat("a", MaxNoteTextLength), Note: true},
		},
		{
			// Post limit (25000) now exceeds the note limit: note-length text
			// without the note flag is accepted as a regular long post.
			name:  "note-length text without the note flag",
			props: DraftProps{Text: strings.Repeat("a", MaxNoteTextLength)},
		},
		{name: "media without text", props: DraftProps{MediaIDs: []string{mediaID}}},
		{
			name: "the maximum number of media items",
			props: DraftProps{Text: "hi", MediaIDs: []string{
				mediaID, otherMediaID, mediaID, otherMediaID,
			}},
		},
		{name: "a poll at the minimum duration", props: DraftProps{
			Text: "hi", Poll: &PollProps{Options: []string{"a", "b"}, DurationMinutes: 5},
		}},
		{name: "a poll at the maximum duration", props: DraftProps{
			Text: "hi", Poll: &PollProps{Options: []string{"a", "b"}, DurationMinutes: 10080},
		}},
		{name: "a poll with four options", props: DraftProps{
			Text: "hi", Poll: &PollProps{Options: []string{"a", "b", "c", "d"}, DurationMinutes: 60},
		}},
		{name: "a reply", props: DraftProps{Text: "hi", ReplyToID: postID}},
		{name: "a quote", props: DraftProps{Text: "hi", QuoteID: postID}},
		{name: "a quote reply", props: DraftProps{Text: "hi", ReplyToID: postID, QuoteID: otherPostID}},
		{name: "reply settings following", props: DraftProps{Text: "hi", ReplySettings: "following"}},
		{name: "reply settings mentionedUsers", props: DraftProps{Text: "hi", ReplySettings: "mentionedUsers"}},
		{name: "reply settings subscribers", props: DraftProps{Text: "hi", ReplySettings: "subscribers"}},
		{name: "reply settings verified", props: DraftProps{Text: "hi", ReplySettings: "verified"}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, err := NewDraft(c.props)

			require.NoError(t, err)
		})
	}
}

func TestDraftGetters(t *testing.T) {
	t.Run("a bare text draft leaves every optional field empty", func(t *testing.T) {
		draft, err := NewDraft(DraftProps{Text: "  hello  "})
		require.NoError(t, err)

		assert.Equal(t, "hello", draft.Text())
		assert.True(t, draft.ReplyTo().IsZero())
		assert.True(t, draft.Quote().IsZero())
		assert.Empty(t, draft.MediaIDs())
		assert.Empty(t, draft.ReplySettings())

		_, ok := draft.Poll()
		assert.False(t, ok)
	})

	t.Run("every field is carried through", func(t *testing.T) {
		draft, err := NewDraft(DraftProps{
			Text:          "hello",
			ReplyToID:     postID,
			QuoteID:       otherPostID,
			MediaIDs:      []string{mediaID, otherMediaID},
			ReplySettings: "following",
		})
		require.NoError(t, err)

		assert.Equal(t, "hello", draft.Text())
		assert.Equal(t, postID, draft.ReplyTo().String())
		assert.Equal(t, otherPostID, draft.Quote().String())
		assert.Equal(t, []string{mediaID, otherMediaID}, draft.MediaIDs())
		assert.Equal(t, "following", draft.ReplySettings())
	})

	t.Run("a poll is reported as present", func(t *testing.T) {
		draft, err := NewDraft(DraftProps{
			Text: "hi",
			Poll: &PollProps{Options: []string{"yes", "no"}, DurationMinutes: 120},
		})
		require.NoError(t, err)

		poll, ok := draft.Poll()

		require.True(t, ok)
		assert.Equal(t, []string{"yes", "no"}, poll.Options())
		assert.Equal(t, 120, poll.DurationMinutes())
	})
}

func TestDraftDoesNotShareStateWithProps(t *testing.T) {
	t.Run("media ids are copied out of the props", func(t *testing.T) {
		ids := []string{mediaID, otherMediaID}

		draft, err := NewDraft(DraftProps{Text: "hi", MediaIDs: ids})
		require.NoError(t, err)

		ids[0] = "999"

		assert.Equal(t, []string{mediaID, otherMediaID}, draft.MediaIDs())
	})

	t.Run("poll options are copied out of the props", func(t *testing.T) {
		options := []string{"yes", "no"}

		draft, err := NewDraft(DraftProps{
			Text: "hi",
			Poll: &PollProps{Options: options, DurationMinutes: 60},
		})
		require.NoError(t, err)

		options[0] = "maybe"

		poll, ok := draft.Poll()
		require.True(t, ok)
		assert.Equal(t, []string{"yes", "no"}, poll.Options())
	})

	t.Run("the media ids getter does not expose the internal slice", func(t *testing.T) {
		draft, err := NewDraft(DraftProps{Text: "hi", MediaIDs: []string{mediaID, otherMediaID}})
		require.NoError(t, err)

		first := draft.MediaIDs()
		first[0] = "999"

		assert.Equal(t, []string{mediaID, otherMediaID}, draft.MediaIDs())
	})

	t.Run("the poll options getter does not expose the internal slice", func(t *testing.T) {
		draft, err := NewDraft(DraftProps{Text: "hi", Poll: validPoll()})
		require.NoError(t, err)

		poll, ok := draft.Poll()
		require.True(t, ok)

		first := poll.Options()
		first[0] = "mutated"

		assert.Equal(t, []string{"yes", "no"}, poll.Options())
	})
}
