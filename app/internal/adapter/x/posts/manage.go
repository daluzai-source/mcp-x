package posts

import (
	"context"
	"fmt"

	"github.com/michimani/gotwi"
	"github.com/michimani/gotwi/tweet/managetweet"
	"github.com/michimani/gotwi/tweet/managetweet/types"

	"github.com/Role1776/mcp-x/app/internal/adapter/x/apierr"
	"github.com/Role1776/mcp-x/app/internal/domain/x"
	"github.com/Role1776/mcp-x/app/internal/domain/x/posts"
)

func (a *Adapter) Create(ctx context.Context, draft posts.Draft) (string, string, error) {
	const op = "adapter.posts.Create"

	in := &types.CreateInput{}

	if text := draft.Text(); text != "" {
		in.Text = gotwi.String(text)
	}

	if replyTo := draft.ReplyTo(); !replyTo.IsZero() {
		in.Reply = &types.CreateInputReply{InReplyToTweetID: replyTo.String()}
	}

	if quote := draft.Quote(); !quote.IsZero() {
		in.QuoteTweetID = gotwi.String(quote.String())
	}

	if mediaIDs := draft.MediaIDs(); len(mediaIDs) > 0 {
		in.Media = &types.CreateInputMedia{MediaIDs: mediaIDs}
	}

	if poll, ok := draft.Poll(); ok {
		in.Poll = &types.CreateInputPoll{
			Options:         poll.Options(),
			DurationMinutes: gotwi.Int(poll.DurationMinutes()),
		}
	}

	if settings := draft.ReplySettings(); settings != "" {
		in.ReplySettings = gotwi.String(settings)
	}

	out, err := managetweet.Create(ctx, a.client, in)
	if err != nil {
		return "", "", fmt.Errorf("%s: %w", op, apierr.ToDomain(err))
	}

	return gotwi.StringValue(out.Data.ID), gotwi.StringValue(out.Data.Text), nil
}

func (a *Adapter) Delete(ctx context.Context, id posts.PostID) error {
	const op = "adapter.posts.Delete"

	out, err := managetweet.Delete(ctx, a.client, &types.DeleteInput{ID: id.String()})
	if err != nil {
		return fmt.Errorf("%s: %w", op, apierr.NotFoundAs(err, x.ErrPostNotFound))
	}

	if !gotwi.BoolValue(out.Data.Deleted) {
		return fmt.Errorf("%s: %w", op, x.ErrForbidden)
	}

	return nil
}
