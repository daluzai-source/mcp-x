package apierr

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/michimani/gotwi"
	"github.com/michimani/gotwi/resources"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Role1776/mcp-x/app/internal/domain/x"
)

func apiError(status int, problemType, detail string) error {
	return &gotwi.GotwiError{
		OnAPI: true,
		Non2XXError: resources.Non2XXError{
			StatusCode: status,
			Type:       problemType,
			Detail:     detail,
		},
	}
}

// rawAPIError reproduces the shape X returns when gotwi cannot parse the
// problem details: Type and Detail stay empty and the body arrives verbatim.
func rawAPIError(status int, body string) error {
	return &gotwi.GotwiError{
		OnAPI: true,
		Non2XXError: resources.Non2XXError{
			StatusCode: status,
			APIErrors:  []resources.ErrorInformation{{Message: body}},
		},
	}
}

func TestToDomain(t *testing.T) {
	cases := []struct {
		name string
		err  error
		want error
	}{
		{name: "nil stays nil", err: nil, want: nil},
		{
			name: "read-only token",
			err:  apiError(403, "https://api.x.com/2/problems/oauth1-permissions", ""),
			want: x.ErrReadOnlyToken,
		},
		{
			name: "usage capped",
			err:  apiError(429, "https://api.x.com/2/problems/usage-capped", ""),
			want: x.ErrUsageCapped,
		},
		{
			name: "credits depleted by type",
			err:  apiError(402, "https://api.x.com/2/problems/credits-depleted", ""),
			want: x.ErrCreditsDepleted,
		},
		{
			name: "credits depleted from an unparsed body",
			err: rawAPIError(402, `{"detail":"credits depleted","status":402,`+
				`"title":"Payment Required","type":"https://api.x.com/2/problems/credits-depleted"}`),
			want: x.ErrCreditsDepleted,
		},
		{
			name: "unparsed body still classifies other problem types",
			err: rawAPIError(403, `{"detail":"usage capped","status":403,`+
				`"title":"Usage capped","type":"https://api.x.com/2/problems/usage-capped"}`),
			want: x.ErrUsageCapped,
		},
		{
			name: "rate limited by type",
			err:  apiError(429, "https://api.x.com/2/problems/rate-limit-exceeded", ""),
			want: x.ErrRateLimited,
		},
		{
			name: "duplicate by type",
			err:  apiError(403, "https://api.x.com/2/problems/duplicate-rules", ""),
			want: x.ErrDuplicatePost,
		},
		{
			name: "duplicate by detail",
			err:  apiError(403, "", "You are not allowed to create a Tweet with duplicate content."),
			want: x.ErrDuplicatePost,
		},
		{
			name: "not authorized for resource",
			err:  apiError(403, "https://api.x.com/2/problems/not-authorized-for-resource", ""),
			want: x.ErrForbidden,
		},
		{
			name: "client forbidden",
			err:  apiError(403, "https://api.x.com/2/problems/client-forbidden", ""),
			want: x.ErrForbidden,
		},
		{
			name: "resource not found",
			err:  apiError(404, "https://api.x.com/2/problems/resource-not-found", ""),
			want: x.ErrPostNotFound,
		},
		{
			name: "invalid request",
			err:  apiError(400, "https://api.x.com/2/problems/invalid-request", ""),
			want: x.ErrInvalidRequest,
		},
		{name: "401 without a type", err: apiError(401, "", ""), want: x.ErrUnauthorized},
		{name: "403 without a type", err: apiError(403, "", ""), want: x.ErrForbidden},
		{name: "404 without a type", err: apiError(404, "", ""), want: x.ErrPostNotFound},
		{name: "402 without a type", err: apiError(402, "", ""), want: x.ErrCreditsDepleted},
		{name: "429 without a type", err: apiError(429, "", ""), want: x.ErrRateLimited},
		{name: "400 without a type", err: apiError(400, "", ""), want: x.ErrInvalidRequest},
		{name: "422 without a type", err: apiError(422, "", ""), want: x.ErrInvalidRequest},
		{name: "500", err: apiError(500, "", ""), want: x.ErrUpstreamUnavailable},
		{name: "503", err: apiError(503, "", ""), want: x.ErrUpstreamUnavailable},
		{name: "unmapped status", err: apiError(418, "", ""), want: x.ErrUpstreamUnavailable},
		{name: "not a gotwi error", err: errors.New("boom"), want: x.ErrUpstreamUnavailable},
		{name: "deadline is passed through", err: context.DeadlineExceeded, want: context.DeadlineExceeded},
		{name: "cancel is passed through", err: context.Canceled, want: context.Canceled},
		{
			name: "wrapped deadline is passed through",
			err:  fmt.Errorf("call: %w", context.DeadlineExceeded),
			want: context.DeadlineExceeded,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := ToDomain(c.err)

			if c.want == nil {
				require.NoError(t, got)

				return
			}

			require.ErrorIs(t, got, c.want)
		})
	}
}

func TestNotFoundAs(t *testing.T) {
	t.Run("substitutes the zone error on 404", func(t *testing.T) {
		got := NotFoundAs(apiError(404, "", ""), x.ErrUserNotFound)

		require.ErrorIs(t, got, x.ErrUserNotFound)
		assert.NotErrorIs(t, got, x.ErrPostNotFound)
	})

	t.Run("leaves other errors alone", func(t *testing.T) {
		got := NotFoundAs(apiError(401, "", ""), x.ErrUserNotFound)

		require.ErrorIs(t, got, x.ErrUnauthorized)
	})
}
