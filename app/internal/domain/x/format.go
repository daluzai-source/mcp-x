package x

import "regexp"

var (
	// X object ids are decimal uint64 snowflakes. Capped at 19 digits: the
	// largest uint64 has 20, so any match is guaranteed to fit without
	// checking for overflow separately.
	snowflakePattern = regexp.MustCompile(`^[0-9]{1,19}$`)

	usernamePattern = regexp.MustCompile(`^[A-Za-z0-9_]{1,15}$`)
)

// IsSnowflakeID reports whether s is a well-formed X object id. Shared by every
// id value object so that posts, users, lists and media cannot drift apart.
func IsSnowflakeID(s string) bool {
	return snowflakePattern.MatchString(s)
}

// IsUsername reports whether s is a well-formed X handle, without the leading @.
func IsUsername(s string) bool {
	return usernamePattern.MatchString(s)
}
