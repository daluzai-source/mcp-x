package media

import (
	"strings"

	"github.com/Role1776/mcp-x/app/internal/domain/x"
)

type MediaID struct {
	value string
}

func NewMediaID(raw string) (MediaID, error) {
	value := strings.TrimSpace(raw)

	if !x.IsSnowflakeID(value) {
		return MediaID{}, x.ErrInvalidMediaID
	}

	return MediaID{value: value}, nil
}

func (m MediaID) String() string {
	return m.value
}
