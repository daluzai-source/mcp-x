package x

import "errors"

var (
	ErrInvalidRequest = errors.New("invalid request")

	ErrEmptyQuery   = errors.New("empty query")
	ErrQueryTooLong = errors.New("query too long")

	ErrInvalidUsername = errors.New("invalid username")
	ErrInvalidUserID   = errors.New("invalid user id")
	ErrInvalidPostID   = errors.New("invalid post id")
	ErrInvalidListID   = errors.New("invalid list id")
	ErrInvalidMediaID  = errors.New("invalid media id")

	ErrTooManyQueries = errors.New("too many queries")
	ErrTooManyUsers   = errors.New("too many users")
	ErrTooManyIDs     = errors.New("too many ids")
	ErrTooManyMedia   = errors.New("too many media")

	ErrEmptyText              = errors.New("empty text")
	ErrTextTooLong            = errors.New("text too long")
	ErrInvalidPoll            = errors.New("invalid poll")
	ErrInvalidReplySettings   = errors.New("invalid reply settings")
	ErrListNameTooLong        = errors.New("list name too long")
	ErrListDescriptionTooLong = errors.New("list description too long")

	ErrMediaEmpty           = errors.New("media is empty")
	ErrMediaTooLarge        = errors.New("media too large")
	ErrUnsupportedMedia     = errors.New("unsupported media type")
	ErrMediaPathOutsideRoot = errors.New("media path outside root")

	ErrUserNotFound          = errors.New("user not found")
	ErrPostNotFound          = errors.New("post not found")
	ErrListNotFound          = errors.New("list not found")
	ErrMediaNotFound         = errors.New("media not found")
	ErrMediaProcessingFailed = errors.New("media processing failed")

	ErrAllQueriesFailed = errors.New("all queries failed")
	ErrAllUsersFailed   = errors.New("all users failed")
	ErrDuplicatePost    = errors.New("duplicate post")

	ErrUnauthorized        = errors.New("unauthorized")
	ErrForbidden           = errors.New("forbidden")
	ErrReadOnlyToken       = errors.New("read-only access token")
	ErrRateLimited         = errors.New("rate limited")
	ErrUsageCapped         = errors.New("usage capped")
	ErrCreditsDepleted     = errors.New("api credits depleted")
	ErrUpstreamUnavailable = errors.New("upstream unavailable")
)
