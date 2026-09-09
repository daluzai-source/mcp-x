package posts

import (
	"strings"

	"github.com/Role1776/mcp-x/app/internal/domain/x"
)

type AuthorID struct {
	value string
}

func NewAuthorID(raw string) (AuthorID, error) {
	value := strings.TrimSpace(raw)

	if !x.IsSnowflakeID(value) {
		return AuthorID{}, x.ErrInvalidUserID
	}

	return AuthorID{value: value}, nil
}

func (a AuthorID) String() string {
	return a.value
}

type AuthorUsername struct {
	value string
}

func NewAuthorUsername(raw string) (AuthorUsername, error) {
	value := strings.TrimPrefix(strings.TrimSpace(raw), "@")

	if !x.IsUsername(value) {
		return AuthorUsername{}, x.ErrInvalidUsername
	}

	return AuthorUsername{value: value}, nil
}

func (a AuthorUsername) String() string {
	return a.value
}
