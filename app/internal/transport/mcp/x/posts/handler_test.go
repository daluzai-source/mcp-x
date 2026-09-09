package posts

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/Role1776/mcp-x/app/internal/domain/x"
	dto "github.com/Role1776/mcp-x/app/internal/dto/x/posts"
)

func TestSearch(t *testing.T) {
	req := dto.SearchRequest{Queries: []string{"golang"}, MaxResults: 3}

	okResponse := dto.SearchResponse{
		Results: []dto.SearchResult{{
			Query:  "golang",
			Status: "success",
			Count:  1,
			Posts:  []dto.Post{{ID: "1", Text: "hello", URL: "https://x.com/i/status/1"}},
		}},
		Metadata: dto.Metadata{TotalRequestTimeMs: 12},
	}

	cases := []struct {
		name string
		res  dto.SearchResponse
		err  error

		wantIsError bool
		wantText    string
		wantRes     dto.SearchResponse
	}{
		{name: "success", res: okResponse, wantRes: okResponse},
		{name: "success with empty results", res: dto.SearchResponse{}, wantRes: dto.SearchResponse{}},
		{
			name:        "rate limited",
			err:         x.ErrRateLimited,
			wantIsError: true,
			wantText:    "rate limit exceeded; the window resets in about 15 minutes, so do not retry immediately",
		},
		{
			name:        "usage capped",
			err:         x.ErrUsageCapped,
			wantIsError: true,
			wantText: "the monthly usage cap for this project is reached; every further call fails until it " +
				"resets, so do not retry",
		},
		{
			name:        "too many queries",
			err:         x.ErrTooManyQueries,
			wantIsError: true,
			wantText:    "too many queries in one call",
		},
		{
			name:        "query too long",
			err:         x.ErrQueryTooLong,
			wantIsError: true,
			wantText:    "query exceeds 512 characters",
		},
		{
			name:        "all queries failed",
			err:         x.ErrAllQueriesFailed,
			wantIsError: true,
			wantText:    "every query failed; the X API may be unreachable",
		},
		{
			name:        "timeout",
			err:         context.DeadlineExceeded,
			wantIsError: true,
			wantText:    "the call timed out; retry with a larger timeout_ms or fewer inputs",
		},
		{
			name:        "unknown error",
			err:         errors.New("boom"),
			wantIsError: true,
			wantText:    "internal server error",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			usecase := NewMockusecase(ctrl)
			usecase.EXPECT().Search(gomock.Any(), req).Return(c.res, c.err)

			res, out, err := New(usecase).search(context.Background(), nil, req)
			require.NoError(t, err)

			if c.wantIsError {
				assert.True(t, res.IsError)
				assert.Equal(t, c.wantText, resultText(t, res))
				assert.Equal(t, dto.SearchResponse{}, out)

				return
			}

			assert.Nil(t, res)
			assert.Equal(t, c.wantRes, out)
		})
	}
}

func TestCreate(t *testing.T) {
	req := dto.CreateRequest{Text: "hello"}

	cases := []struct {
		name string
		res  dto.CreateResponse
		err  error

		wantIsError bool
		wantText    string
	}{
		{name: "success", res: dto.CreateResponse{ID: "1", Text: "hello"}},
		{
			name:        "read-only token",
			err:         x.ErrReadOnlyToken,
			wantIsError: true,
			wantText: "the access token is read-only: open the X Developer Console, set the app permissions " +
				"to Read and Write, then regenerate the Access Token and Secret — an existing token keeps " +
				"the permissions it was created with",
		},
		{
			name:        "duplicate",
			err:         x.ErrDuplicatePost,
			wantIsError: true,
			wantText:    "X rejected this as a duplicate of a recent post; change the text before trying again",
		},
		{
			name:        "text too long",
			err:         x.ErrTextTooLong,
			wantIsError: true,
			wantText:    "post text exceeds 280 characters",
		},
		{
			name:        "bad poll",
			err:         x.ErrInvalidPoll,
			wantIsError: true,
			wantText:    "a poll needs 2 to 4 options and a duration between 5 and 10080 minutes",
		},
		{
			name:        "too many media",
			err:         x.ErrTooManyMedia,
			wantIsError: true,
			wantText:    "a post carries at most 4 media items",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			usecase := NewMockusecase(ctrl)
			usecase.EXPECT().Create(gomock.Any(), req).Return(c.res, c.err)

			res, out, err := New(usecase).create(context.Background(), nil, req)
			require.NoError(t, err)

			if c.wantIsError {
				assert.True(t, res.IsError)
				assert.Equal(t, c.wantText, resultText(t, res))

				return
			}

			assert.Nil(t, res)
			assert.Equal(t, c.res, out)
		})
	}
}

func TestDelete(t *testing.T) {
	req := dto.PostIDRequest{ID: "1"}

	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		usecase := NewMockusecase(ctrl)
		usecase.EXPECT().Delete(gomock.Any(), req).Return(dto.ActionResponse{ID: "1", Status: "done"}, nil)

		res, out, err := New(usecase).delete(context.Background(), nil, req)

		require.NoError(t, err)
		assert.Nil(t, res)
		assert.Equal(t, "done", out.Status)
	})

	t.Run("not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		usecase := NewMockusecase(ctrl)
		usecase.EXPECT().Delete(gomock.Any(), req).Return(dto.ActionResponse{}, x.ErrPostNotFound)

		res, _, err := New(usecase).delete(context.Background(), nil, req)

		require.NoError(t, err)
		assert.True(t, res.IsError)
		assert.Equal(t, "post not found", resultText(t, res))
	})

	t.Run("forbidden", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		usecase := NewMockusecase(ctrl)
		usecase.EXPECT().Delete(gomock.Any(), req).Return(dto.ActionResponse{}, x.ErrForbidden)

		res, _, err := New(usecase).delete(context.Background(), nil, req)

		require.NoError(t, err)
		assert.Equal(t,
			"not authorized for this resource: it may be protected, deleted, or owned by somebody else",
			resultText(t, res))
	})
}

func TestLookup(t *testing.T) {
	req := dto.LookupRequest{IDs: []string{"1", "2"}}

	t.Run("passes partial results through", func(t *testing.T) {
		expected := dto.LookupResponse{
			Results: []dto.LookupResult{
				{ID: "1", Status: "success", Post: &dto.Post{ID: "1"}},
				{ID: "2", Status: "not found", Detail: "gone"},
			},
			Count: 1,
		}

		ctrl := gomock.NewController(t)
		usecase := NewMockusecase(ctrl)
		usecase.EXPECT().Lookup(gomock.Any(), req).Return(expected, nil)

		res, out, err := New(usecase).lookup(context.Background(), nil, req)

		require.NoError(t, err)
		assert.Nil(t, res)
		assert.Equal(t, expected, out)
	})

	t.Run("too many ids", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		usecase := NewMockusecase(ctrl)
		usecase.EXPECT().Lookup(gomock.Any(), req).Return(dto.LookupResponse{}, x.ErrTooManyIDs)

		res, _, err := New(usecase).lookup(context.Background(), nil, req)

		require.NoError(t, err)
		assert.Equal(t, "too many ids in one call; the limit is 100", resultText(t, res))
	})
}
