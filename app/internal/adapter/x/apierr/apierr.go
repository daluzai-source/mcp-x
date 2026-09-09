package apierr

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/michimani/gotwi"

	"github.com/Role1776/mcp-x/app/internal/domain/x"
)

func ToDomain(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
		return err
	}

	var apiErr *gotwi.GotwiError
	if !errors.As(err, &apiErr) {
		return x.ErrUpstreamUnavailable
	}

	if mapped, ok := byType(problem(apiErr)); ok {
		return mapped
	}

	return byStatus(apiErr.StatusCode)
}

func NotFoundAs(err error, target error) error {
	if errors.Is(ToDomain(err), x.ErrPostNotFound) {
		return target
	}

	return ToDomain(err)
}

// problem collects everything the response can carry about the failure.
//
// X answers errors with Content-Type "application/problem+json", and gotwi
// only decodes a body when the type contains "application/json" — which that
// value does not match. So for every X problem response Type and Detail stay
// empty and the whole JSON lands verbatim in APIErrors[].Message. When gotwi
// does decode a body, the reverse holds: resources.ErrorInformation has no
// type or detail field, so the per-error entries carry nothing to match on.
// The two sources never compete, which is why they are simply concatenated.
func problem(apiErr *gotwi.GotwiError) string {
	var b strings.Builder

	b.WriteString(apiErr.Type)
	b.WriteString(" ")
	b.WriteString(apiErr.Detail)

	for _, e := range apiErr.APIErrors {
		b.WriteString(" ")
		b.WriteString(e.Message)
	}

	return b.String()
}

// byType matches the problem type and the human-readable detail against the
// same haystack: X puts both in one body, and the patterns are distinct enough
// that widening the match costs nothing.
func byType(problem string) (error, bool) {
	problem = strings.ToLower(problem)

	switch {
	case strings.Contains(problem, "oauth1-permissions"):
		return x.ErrReadOnlyToken, true
	case strings.Contains(problem, "usage-capped"):
		return x.ErrUsageCapped, true
	case strings.Contains(problem, "credits-depleted"):
		return x.ErrCreditsDepleted, true
	case strings.Contains(problem, "rate-limit-exceeded"):
		return x.ErrRateLimited, true
	// covers both the "duplicate-rules" problem type and the "duplicate content"
	// wording X uses when it rejects a repeated post.
	case strings.Contains(problem, "duplicate"):
		return x.ErrDuplicatePost, true
	case strings.Contains(problem, "not-authorized-for-resource"),
		strings.Contains(problem, "client-forbidden"),
		strings.Contains(problem, "unsupported-authentication"):
		return x.ErrForbidden, true
	case strings.Contains(problem, "resource-not-found"):
		return x.ErrPostNotFound, true
	case strings.Contains(problem, "invalid-request"):
		return x.ErrInvalidRequest, true
	default:
		return nil, false
	}
}

func byStatus(status int) error {
	switch {
	case status == http.StatusUnauthorized:
		return x.ErrUnauthorized
	case status == http.StatusForbidden:
		return x.ErrForbidden
	case status == http.StatusNotFound:
		return x.ErrPostNotFound
	case status == http.StatusPaymentRequired:
		return x.ErrCreditsDepleted
	case status == http.StatusTooManyRequests:
		return x.ErrRateLimited
	case status == http.StatusBadRequest || status == http.StatusUnprocessableEntity:
		return x.ErrInvalidRequest
	case status >= http.StatusInternalServerError:
		return x.ErrUpstreamUnavailable
	default:
		return x.ErrUpstreamUnavailable
	}
}
