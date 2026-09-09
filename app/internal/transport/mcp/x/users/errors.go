package users

import (
	"context"
	"errors"

	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/Role1776/mcp-x/app/internal/domain/x"
)

func result(err error) *mcpsdk.CallToolResult {
	var message string

	switch {
	case errors.Is(err, x.ErrReadOnlyToken):
		message = "the access token is read-only: open the X Developer Console, set the app permissions " +
			"to Read and Write, then regenerate the Access Token and Secret — an existing token keeps " +
			"the permissions it was created with"

	case errors.Is(err, x.ErrUsageCapped):
		message = "the monthly usage cap for this project is reached; every further call fails until it " +
			"resets, so do not retry"

	case errors.Is(err, x.ErrCreditsDepleted):
		message = "the project has no X API credits left; credits do not reset on their own, so every " +
			"further call fails until they are topped up in the X Developer Console — do not retry"

	case errors.Is(err, x.ErrRateLimited):
		message = "rate limit exceeded; the window resets in about 15 minutes, so do not retry immediately"

	case errors.Is(err, x.ErrUnauthorized):
		message = "X rejected the credentials; check the four keys in the server configuration"

	case errors.Is(err, x.ErrForbidden):
		message = "not authorized for this resource: it may be protected, deleted, or owned by somebody else"

	case errors.Is(err, x.ErrDuplicatePost):
		message = "X rejected this as a duplicate of a recent post; change the text before trying again"

	case errors.Is(err, x.ErrUserNotFound):
		message = "user not found"

	case errors.Is(err, x.ErrPostNotFound):
		message = "post not found"

	case errors.Is(err, x.ErrListNotFound):
		message = "list not found"

	case errors.Is(err, x.ErrMediaNotFound):
		message = "media file not found"

	case errors.Is(err, x.ErrMediaProcessingFailed):
		message = "X failed to process the uploaded media; re-encode the file and upload it again"

	case errors.Is(err, x.ErrMediaTooLarge):
		message = "media exceeds the size limit for its category: 5 MB for images, 15 MB for gifs, 512 MB for video"

	case errors.Is(err, x.ErrUnsupportedMedia):
		message = "unsupported media type for this category"

	case errors.Is(err, x.ErrMediaPathOutsideRoot):
		message = "the path is outside the directory this server is allowed to read"

	case errors.Is(err, x.ErrMediaEmpty):
		message = "the media file is empty"

	case errors.Is(err, x.ErrTextTooLong):
		message = "post text exceeds 280 characters"

	case errors.Is(err, x.ErrEmptyText):
		message = "a post needs text unless it carries media"

	case errors.Is(err, x.ErrInvalidPoll):
		message = "a poll needs 2 to 4 options and a duration between 5 and 10080 minutes"

	case errors.Is(err, x.ErrInvalidReplySettings):
		message = "reply_settings must be following, mentionedUsers, subscribers or verified"

	case errors.Is(err, x.ErrTooManyMedia):
		message = "a post carries at most 4 media items"

	case errors.Is(err, x.ErrListNameTooLong):
		message = "list name exceeds 25 characters"

	case errors.Is(err, x.ErrListDescriptionTooLong):
		message = "list description exceeds 100 characters"

	case errors.Is(err, x.ErrQueryTooLong):
		message = "query exceeds 512 characters"

	case errors.Is(err, x.ErrEmptyQuery):
		message = "query must not be empty"

	case errors.Is(err, x.ErrTooManyQueries):
		message = "too many queries in one call"

	case errors.Is(err, x.ErrTooManyUsers):
		message = "too many users in one call"

	case errors.Is(err, x.ErrTooManyIDs):
		message = "too many ids in one call; the limit is 100"

	case errors.Is(err, x.ErrInvalidUsername):
		message = "invalid username: 1 to 15 letters, digits or underscores"

	case errors.Is(err, x.ErrInvalidUserID),
		errors.Is(err, x.ErrInvalidPostID),
		errors.Is(err, x.ErrInvalidListID),
		errors.Is(err, x.ErrInvalidMediaID):
		message = "invalid id"

	case errors.Is(err, x.ErrAllQueriesFailed):
		message = "every query failed; the X API may be unreachable"

	case errors.Is(err, x.ErrAllUsersFailed):
		message = "none of the usernames could be resolved"

	case errors.Is(err, context.DeadlineExceeded):
		message = "the call timed out; retry with a larger timeout_ms or fewer inputs"

	case errors.Is(err, x.ErrUpstreamUnavailable):
		message = "the X API is unavailable"

	case errors.Is(err, x.ErrInvalidRequest):
		message = "invalid request"

	default:
		message = "internal server error"
	}

	return &mcpsdk.CallToolResult{
		IsError: true,
		Content: []mcpsdk.Content{&mcpsdk.TextContent{Text: message}},
	}
}
