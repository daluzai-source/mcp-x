package lists

import (
	"strings"

	"github.com/Role1776/mcp-x/app/internal/domain/x"
)

type ListID struct {
	value string
}

func NewListID(raw string) (ListID, error) {
	value := strings.TrimSpace(raw)

	if !x.IsSnowflakeID(value) {
		return ListID{}, x.ErrInvalidListID
	}

	return ListID{value: value}, nil
}

func (l ListID) String() string {
	return l.value
}

type OwnerID struct {
	value string
}

func NewOwnerID(raw string) (OwnerID, error) {
	value := strings.TrimSpace(raw)

	if !x.IsSnowflakeID(value) {
		return OwnerID{}, x.ErrInvalidUserID
	}

	return OwnerID{value: value}, nil
}

func (o OwnerID) String() string {
	return o.value
}
