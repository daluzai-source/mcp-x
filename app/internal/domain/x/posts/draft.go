package posts

import (
	"slices"
	"strings"
	"unicode/utf8"

	"github.com/Role1776/mcp-x/app/internal/domain/x"
)

const (
	MaxTextLength     = 25000
	MaxNoteTextLength = 4000
	MaxDraftMedia     = 4

	minPollOptions   = 2
	maxPollOptions   = 4
	minPollMinutes   = 5
	maxPollMinutes   = 10080
	maxMediaIDLength = 19
)

var replySettings = map[string]struct{}{
	"following":      {},
	"mentionedUsers": {},
	"subscribers":    {},
	"verified":       {},
}

type Poll struct {
	options         []string
	durationMinutes int
}

func (p Poll) Options() []string {
	return slices.Clone(p.options)
}

func (p Poll) DurationMinutes() int {
	return p.durationMinutes
}

type PollProps struct {
	Options         []string
	DurationMinutes int
}

type Draft struct {
	text          string
	replyTo       PostID
	quote         PostID
	mediaIDs      []string
	poll          *Poll
	replySettings string
}

type DraftProps struct {
	Text          string
	ReplyToID     string
	QuoteID       string
	MediaIDs      []string
	Poll          *PollProps
	ReplySettings string
	Note          bool
}

func NewDraft(props DraftProps) (Draft, error) {
	draft := Draft{
		text:          strings.TrimSpace(props.Text),
		mediaIDs:      slices.Clone(props.MediaIDs),
		replySettings: props.ReplySettings,
	}

	if len(props.MediaIDs) > MaxDraftMedia {
		return Draft{}, x.ErrTooManyMedia
	}

	if draft.text == "" && len(props.MediaIDs) == 0 {
		return Draft{}, x.ErrEmptyText
	}

	limit := MaxTextLength
	if props.Note {
		limit = MaxNoteTextLength
	}

	if utf8.RuneCountInString(draft.text) > limit {
		return Draft{}, x.ErrTextTooLong
	}

	for _, id := range props.MediaIDs {
		if !x.IsSnowflakeID(id) {
			return Draft{}, x.ErrInvalidMediaID
		}
	}

	if props.Poll != nil {
		if len(props.MediaIDs) > 0 {
			return Draft{}, x.ErrInvalidRequest
		}

		poll, err := newPoll(*props.Poll)
		if err != nil {
			return Draft{}, err
		}
		draft.poll = &poll
	}

	if props.ReplySettings != "" {
		if _, ok := replySettings[props.ReplySettings]; !ok {
			return Draft{}, x.ErrInvalidReplySettings
		}
	}

	if props.ReplyToID != "" {
		replyTo, err := NewPostID(props.ReplyToID)
		if err != nil {
			return Draft{}, err
		}
		draft.replyTo = replyTo
	}

	if props.QuoteID != "" {
		quote, err := NewPostID(props.QuoteID)
		if err != nil {
			return Draft{}, err
		}
		draft.quote = quote
	}

	return draft, nil
}

func newPoll(props PollProps) (Poll, error) {
	options := make([]string, 0, len(props.Options))
	for _, option := range props.Options {
		option = strings.TrimSpace(option)
		if option == "" {
			return Poll{}, x.ErrInvalidPoll
		}
		options = append(options, option)
	}

	if len(options) < minPollOptions || len(options) > maxPollOptions {
		return Poll{}, x.ErrInvalidPoll
	}

	if props.DurationMinutes < minPollMinutes || props.DurationMinutes > maxPollMinutes {
		return Poll{}, x.ErrInvalidPoll
	}

	return Poll{options: options, durationMinutes: props.DurationMinutes}, nil
}

func (d Draft) Text() string {
	return d.text
}

func (d Draft) ReplyTo() PostID {
	return d.replyTo
}

func (d Draft) Quote() PostID {
	return d.quote
}

func (d Draft) MediaIDs() []string {
	return slices.Clone(d.mediaIDs)
}

func (d Draft) Poll() (Poll, bool) {
	if d.poll == nil {
		return Poll{}, false
	}

	return *d.poll, true
}

func (d Draft) ReplySettings() string {
	return d.replySettings
}
