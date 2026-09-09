package posts

import (
	"time"

	"github.com/Role1776/mcp-x/app/internal/domain/x"
)

const MaxWindowDays = 7

type Window struct {
	days int
}

func NewWindow(days int) (Window, error) {
	if days < 1 || days > MaxWindowDays {
		return Window{}, x.ErrInvalidRequest
	}

	return Window{days: days}, nil
}

func (w Window) Days() int {
	return w.days
}

func (w Window) StartTime(now time.Time) time.Time {
	return now.Add(-time.Duration(w.days) * 24 * time.Hour)
}

type Exclude struct {
	replies  bool
	retweets bool
}

func NewExclude(values []string) (Exclude, error) {
	var exclude Exclude
	for _, value := range values {
		switch value {
		case "replies":
			exclude.replies = true
		case "retweets":
			exclude.retweets = true
		default:
			return Exclude{}, x.ErrInvalidRequest
		}
	}

	return exclude, nil
}

func (e Exclude) Replies() bool {
	return e.replies
}

func (e Exclude) Retweets() bool {
	return e.retweets
}

func (e Exclude) Strings() []string {
	result := make([]string, 0, 2)
	if e.replies {
		result = append(result, "replies")
	}
	if e.retweets {
		result = append(result, "retweets")
	}

	return result
}
