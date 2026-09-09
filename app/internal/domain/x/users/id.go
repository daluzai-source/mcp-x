package users

import (
	"strings"

	"github.com/Role1776/mcp-x/app/internal/domain/x"
)

type UserID struct {
	value string
}

func NewUserID(raw string) (UserID, error) {
	id := strings.TrimSpace(raw)

	if !x.IsSnowflakeID(id) {
		return UserID{}, x.ErrInvalidUserID
	}

	return UserID{value: id}, nil
}

func (u UserID) String() string {
	return u.value
}

type UserIDs []UserID

func NewUserIDs(raw []string) (UserIDs, error) {
	ids := make(UserIDs, 0, len(raw))
	for _, r := range raw {
		id, err := NewUserID(r)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	if len(ids) == 0 {
		return nil, x.ErrInvalidRequest
	}

	return ids.Dedupe(), nil
}

func (u UserIDs) Len() int {
	return len(u)
}

func (u UserIDs) Strings() []string {
	result := make([]string, 0, len(u))
	for _, id := range u {
		result = append(result, id.value)
	}

	return result
}

func (u UserIDs) Dedupe() UserIDs {
	seen := make(map[string]struct{}, len(u))
	result := make(UserIDs, 0, len(u))

	for _, id := range u {
		if _, ok := seen[id.value]; ok {
			continue
		}
		seen[id.value] = struct{}{}
		result = append(result, id)

	}

	return result
}
