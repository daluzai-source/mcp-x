package users

import (
	"strings"

	"github.com/Role1776/mcp-x/app/internal/domain/x"
)

type Username struct {
	value string
}

func NewUsername(raw string) (Username, error) {
	value := strings.TrimPrefix(strings.TrimSpace(raw), "@")

	if !x.IsUsername(value) {
		return Username{}, x.ErrInvalidUsername
	}

	return Username{value: value}, nil
}

func (u Username) String() string {
	return u.value
}

func (u Username) Normalized() string {
	return strings.ToLower(u.value)
}

type Usernames []Username

func NewUsernames(raw []string) (Usernames, error) {
	names := make(Usernames, 0, len(raw))
	for _, r := range raw {
		name, err := NewUsername(r)
		if err != nil {
			return nil, err
		}
		names = append(names, name)
	}

	if len(names) == 0 {
		return nil, x.ErrInvalidRequest
	}

	return names.Dedupe(), nil
}

func (u Usernames) Len() int {
	return len(u)
}

func (u Usernames) Strings() []string {
	result := make([]string, 0, len(u))
	for _, name := range u {
		result = append(result, name.value)
	}

	return result
}

func (u Usernames) Dedupe() Usernames {
	seen := make(map[string]struct{}, len(u))
	result := make(Usernames, 0, len(u))

	for _, name := range u {
		key := name.Normalized()
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		result = append(result, name)
	}

	return result
}
