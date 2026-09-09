package posts

import (
	"strings"
	"unicode/utf8"

	"github.com/Role1776/mcp-x/app/internal/domain/x"
)

const (
	maxQueryLength    = 512
	noRetweetOperator = "-is:retweet"
)

type Query struct {
	value string
}

func NewQuery(raw string) (Query, error) {
	value := strings.TrimSpace(raw)

	if value == "" {
		return Query{}, x.ErrEmptyQuery
	}

	if utf8.RuneCountInString(value) > maxQueryLength {
		return Query{}, x.ErrQueryTooLong
	}

	return Query{value: value}, nil
}

// WithoutRetweets appends the operator instead of letting the caller build the
// string: a retweet duplicates the original post, costs a paid read and pushes
// fresh results out of the result window.
func (q Query) WithoutRetweets() (Query, error) {
	if strings.Contains(q.value, noRetweetOperator) || strings.Contains(q.value, "is:retweet") {
		return q, nil
	}

	value := q.value + " " + noRetweetOperator
	if utf8.RuneCountInString(value) > maxQueryLength {
		return Query{}, x.ErrQueryTooLong
	}

	return Query{value: value}, nil
}

func (q Query) String() string {
	return q.value
}

type Queries []Query

func NewQueries(raw []string) (Queries, error) {
	queries := make(Queries, 0, len(raw))
	for _, r := range raw {
		query, err := NewQuery(r)
		if err != nil {
			return nil, err
		}
		queries = append(queries, query)
	}

	if len(queries) == 0 {
		return nil, x.ErrEmptyQuery
	}

	return queries, nil
}

func (q Queries) Len() int {
	return len(q)
}

func (q Queries) Strings() []string {
	result := make([]string, 0, len(q))
	for _, query := range q {
		result = append(result, query.String())
	}

	return result
}

func (q Queries) WithoutRetweets() (Queries, error) {
	result := make(Queries, 0, len(q))
	for _, query := range q {
		filtered, err := query.WithoutRetweets()
		if err != nil {
			return nil, err
		}
		result = append(result, filtered)
	}

	return result, nil
}
