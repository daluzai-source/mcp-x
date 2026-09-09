package posts

import (
	"strings"

	"github.com/Role1776/mcp-x/app/internal/domain/x"
)

type PostID struct {
	value string
}

func NewPostID(raw string) (PostID, error) {
	value := strings.TrimSpace(raw)

	if !x.IsSnowflakeID(value) {
		return PostID{}, x.ErrInvalidPostID
	}

	return PostID{value: value}, nil
}

func (p PostID) String() string {
	return p.value
}

func (p PostID) IsZero() bool {
	return p.value == ""
}

type PostIDs []PostID

func NewPostIDs(raw []string) (PostIDs, error) {
	ids := make(PostIDs, 0, len(raw))
	for _, r := range raw {
		id, err := NewPostID(r)
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

func (p PostIDs) Len() int {
	return len(p)
}

func (p PostIDs) Strings() []string {
	result := make([]string, 0, len(p))
	for _, id := range p {
		result = append(result, id.value)
	}

	return result
}

func (p PostIDs) Dedupe() PostIDs {
	seen := make(map[string]struct{}, len(p))
	result := make(PostIDs, 0, len(p))

	for _, id := range p {
		if _, ok := seen[id.value]; ok {
			continue
		}
		seen[id.value] = struct{}{}
		result = append(result, id)
	}

	return result
}
