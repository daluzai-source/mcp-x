package lists

import (
	"context"
	"errors"
	"testing"

	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/Role1776/mcp-x/app/internal/domain/x"
	dto "github.com/Role1776/mcp-x/app/internal/dto/x/lists"
)

func TestCreate(t *testing.T) {
	req := dto.CreateRequest{Name: "reading", Description: "books"}

	cases := []struct {
		name string
		res  dto.ActionResponse
		err  error

		wantIsError bool
		wantText    string
	}{
		{name: "success", res: dto.ActionResponse{ID: "1", Status: "done"}},
		{
			name:        "name too long",
			err:         x.ErrListNameTooLong,
			wantIsError: true,
			wantText:    "list name exceeds 25 characters",
		},
		{
			name:        "description too long",
			err:         x.ErrListDescriptionTooLong,
			wantIsError: true,
			wantText:    "list description exceeds 100 characters",
		},
		{
			name:        "read-only token",
			err:         x.ErrReadOnlyToken,
			wantIsError: true,
			wantText: "the access token is read-only: open the X Developer Console, set the app permissions " +
				"to Read and Write, then regenerate the Access Token and Secret — an existing token keeps " +
				"the permissions it was created with",
		},
		{
			name:        "invalid request",
			err:         x.ErrInvalidRequest,
			wantIsError: true,
			wantText:    "invalid request",
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
			usecase.EXPECT().Create(gomock.Any(), req).Return(c.res, c.err)

			res, out, err := New(usecase).create(context.Background(), nil, req)
			require.NoError(t, err)

			if c.wantIsError {
				assert.True(t, res.IsError)
				assert.Equal(t, c.wantText, resultText(t, res))
				assert.Equal(t, dto.ActionResponse{}, out)

				return
			}

			assert.Nil(t, res)
			assert.Equal(t, c.res, out)
		})
	}
}

func TestUpdate(t *testing.T) {
	name := "reading"
	req := dto.UpdateRequest{ID: "1", Name: &name}

	cases := []struct {
		name string
		res  dto.ActionResponse
		err  error

		wantIsError bool
		wantText    string
	}{
		{name: "success", res: dto.ActionResponse{ID: "1", Status: "done"}},
		{
			name:        "list not found",
			err:         x.ErrListNotFound,
			wantIsError: true,
			wantText:    "list not found",
		},
		{
			name:        "invalid id",
			err:         x.ErrInvalidListID,
			wantIsError: true,
			wantText:    "invalid id",
		},
		{
			name:        "name too long",
			err:         x.ErrListNameTooLong,
			wantIsError: true,
			wantText:    "list name exceeds 25 characters",
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
			usecase.EXPECT().Update(gomock.Any(), req).Return(c.res, c.err)

			res, out, err := New(usecase).update(context.Background(), nil, req)
			require.NoError(t, err)

			if c.wantIsError {
				assert.True(t, res.IsError)
				assert.Equal(t, c.wantText, resultText(t, res))
				assert.Equal(t, dto.ActionResponse{}, out)

				return
			}

			assert.Nil(t, res)
			assert.Equal(t, c.res, out)
		})
	}
}

func TestGet(t *testing.T) {
	req := dto.ListIDRequest{ID: "1"}

	okResponse := dto.ListResponse{List: dto.List{ID: "1", Name: "reading", URL: "https://x.com/i/lists/1"}}

	cases := []struct {
		name string
		res  dto.ListResponse
		err  error

		wantIsError bool
		wantText    string
	}{
		{name: "success", res: okResponse},
		{
			name:        "list not found",
			err:         x.ErrListNotFound,
			wantIsError: true,
			wantText:    "list not found",
		},
		{
			name:        "forbidden",
			err:         x.ErrForbidden,
			wantIsError: true,
			wantText:    "not authorized for this resource: it may be protected, deleted, or owned by somebody else",
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
			usecase.EXPECT().Get(gomock.Any(), req).Return(c.res, c.err)

			res, out, err := New(usecase).get(context.Background(), nil, req)
			require.NoError(t, err)

			if c.wantIsError {
				assert.True(t, res.IsError)
				assert.Equal(t, c.wantText, resultText(t, res))
				assert.Equal(t, dto.ListResponse{}, out)

				return
			}

			assert.Nil(t, res)
			assert.Equal(t, c.res, out)
		})
	}
}

func TestOwned(t *testing.T) {
	req := dto.OwnedRequest{Username: "jack"}

	okResponse := dto.ListsResponse{Lists: []dto.List{{ID: "1", Name: "reading"}}, Count: 1, NextToken: "next"}

	cases := []struct {
		name string
		res  dto.ListsResponse
		err  error

		wantIsError bool
		wantText    string
	}{
		{name: "success", res: okResponse},
		{name: "success with an empty page", res: dto.ListsResponse{}},
		{
			name:        "user not found",
			err:         x.ErrUserNotFound,
			wantIsError: true,
			wantText:    "user not found",
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
			usecase.EXPECT().Owned(gomock.Any(), req).Return(c.res, c.err)

			res, out, err := New(usecase).owned(context.Background(), nil, req)
			require.NoError(t, err)

			if c.wantIsError {
				assert.True(t, res.IsError)
				assert.Equal(t, c.wantText, resultText(t, res))
				assert.Equal(t, dto.ListsResponse{}, out)

				return
			}

			assert.Nil(t, res)
			assert.Equal(t, c.res, out)
		})
	}
}

func TestPosts(t *testing.T) {
	req := dto.ListPageRequest{ID: "1", MaxResults: 10}

	okResponse := dto.PostsResponse{Posts: []dto.Post{{ID: "2", Text: "hello"}}, Count: 1}

	cases := []struct {
		name string
		res  dto.PostsResponse
		err  error

		wantIsError bool
		wantText    string
	}{
		{name: "success", res: okResponse},
		{
			name:        "list not found",
			err:         x.ErrListNotFound,
			wantIsError: true,
			wantText:    "list not found",
		},
		{
			name:        "rate limited",
			err:         x.ErrRateLimited,
			wantIsError: true,
			wantText:    "rate limit exceeded; the window resets in about 15 minutes, so do not retry immediately",
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
			usecase.EXPECT().Posts(gomock.Any(), req).Return(c.res, c.err)

			res, out, err := New(usecase).posts(context.Background(), nil, req)
			require.NoError(t, err)

			if c.wantIsError {
				assert.True(t, res.IsError)
				assert.Equal(t, c.wantText, resultText(t, res))
				assert.Equal(t, dto.PostsResponse{}, out)

				return
			}

			assert.Nil(t, res)
			assert.Equal(t, c.res, out)
		})
	}
}

func TestMembers(t *testing.T) {
	req := dto.ListPageRequest{ID: "1"}

	okResponse := dto.UsersResponse{Users: []dto.User{{ID: "2", Username: "jack"}}, Count: 1}

	cases := []struct {
		name string
		res  dto.UsersResponse
		err  error

		wantIsError bool
		wantText    string
	}{
		{name: "success", res: okResponse},
		{
			name:        "list not found",
			err:         x.ErrListNotFound,
			wantIsError: true,
			wantText:    "list not found",
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
			usecase.EXPECT().Members(gomock.Any(), req).Return(c.res, c.err)

			res, out, err := New(usecase).members(context.Background(), nil, req)
			require.NoError(t, err)

			if c.wantIsError {
				assert.True(t, res.IsError)
				assert.Equal(t, c.wantText, resultText(t, res))
				assert.Equal(t, dto.UsersResponse{}, out)

				return
			}

			assert.Nil(t, res)
			assert.Equal(t, c.res, out)
		})
	}
}

func TestMemberTools(t *testing.T) {
	req := dto.MemberRequest{ID: "1", Username: "jack"}

	okResponse := dto.ActionResponse{ID: "1", Status: "done"}

	tools := []struct {
		name   string
		expect func(*MockusecaseMockRecorder, dto.ActionResponse, error)
		call   func(*Handler, context.Context, dto.MemberRequest) (*mcpsdk.CallToolResult, dto.ActionResponse, error)
	}{
		{
			name: "member add",
			expect: func(rec *MockusecaseMockRecorder, res dto.ActionResponse, err error) {
				rec.MemberAdd(gomock.Any(), req).Return(res, err)
			},
			call: func(h *Handler, ctx context.Context, req dto.MemberRequest) (*mcpsdk.CallToolResult, dto.ActionResponse, error) {
				return h.memberAdd(ctx, nil, req)
			},
		},
		{
			name: "member remove",
			expect: func(rec *MockusecaseMockRecorder, res dto.ActionResponse, err error) {
				rec.MemberRemove(gomock.Any(), req).Return(res, err)
			},
			call: func(h *Handler, ctx context.Context, req dto.MemberRequest) (*mcpsdk.CallToolResult, dto.ActionResponse, error) {
				return h.memberRemove(ctx, nil, req)
			},
		},
	}

	outcomes := []struct {
		name string
		res  dto.ActionResponse
		err  error

		wantIsError bool
		wantText    string
	}{
		{name: "success", res: okResponse},
		{
			name:        "user not found",
			err:         x.ErrUserNotFound,
			wantIsError: true,
			wantText:    "user not found",
		},
		{
			name:        "list not found",
			err:         x.ErrListNotFound,
			wantIsError: true,
			wantText:    "list not found",
		},
		{
			name:        "invalid username",
			err:         x.ErrInvalidUsername,
			wantIsError: true,
			wantText:    "invalid username: 1 to 15 letters, digits or underscores",
		},
		{
			name:        "unknown error",
			err:         errors.New("boom"),
			wantIsError: true,
			wantText:    "internal server error",
		},
	}

	for _, tool := range tools {
		for _, c := range outcomes {
			t.Run(tool.name+": "+c.name, func(t *testing.T) {
				ctrl := gomock.NewController(t)
				usecase := NewMockusecase(ctrl)
				tool.expect(usecase.EXPECT(), c.res, c.err)

				res, out, err := tool.call(New(usecase), context.Background(), req)
				require.NoError(t, err)

				if c.wantIsError {
					assert.True(t, res.IsError)
					assert.Equal(t, c.wantText, resultText(t, res))
					assert.Equal(t, dto.ActionResponse{}, out)

					return
				}

				assert.Nil(t, res)
				assert.Equal(t, c.res, out)
			})
		}
	}
}

func TestListIDTools(t *testing.T) {
	req := dto.ListIDRequest{ID: "1"}

	okResponse := dto.ActionResponse{ID: "1", Status: "done"}

	tools := []struct {
		name   string
		expect func(*MockusecaseMockRecorder, dto.ActionResponse, error)
		call   func(*Handler, context.Context, dto.ListIDRequest) (*mcpsdk.CallToolResult, dto.ActionResponse, error)
	}{
		{
			name: "delete",
			expect: func(rec *MockusecaseMockRecorder, res dto.ActionResponse, err error) {
				rec.Delete(gomock.Any(), req).Return(res, err)
			},
			call: func(h *Handler, ctx context.Context, req dto.ListIDRequest) (*mcpsdk.CallToolResult, dto.ActionResponse, error) {
				return h.delete(ctx, nil, req)
			},
		},
		{
			name: "follow",
			expect: func(rec *MockusecaseMockRecorder, res dto.ActionResponse, err error) {
				rec.Follow(gomock.Any(), req).Return(res, err)
			},
			call: func(h *Handler, ctx context.Context, req dto.ListIDRequest) (*mcpsdk.CallToolResult, dto.ActionResponse, error) {
				return h.follow(ctx, nil, req)
			},
		},
		{
			name: "unfollow",
			expect: func(rec *MockusecaseMockRecorder, res dto.ActionResponse, err error) {
				rec.Unfollow(gomock.Any(), req).Return(res, err)
			},
			call: func(h *Handler, ctx context.Context, req dto.ListIDRequest) (*mcpsdk.CallToolResult, dto.ActionResponse, error) {
				return h.unfollow(ctx, nil, req)
			},
		},
		{
			name: "pin",
			expect: func(rec *MockusecaseMockRecorder, res dto.ActionResponse, err error) {
				rec.Pin(gomock.Any(), req).Return(res, err)
			},
			call: func(h *Handler, ctx context.Context, req dto.ListIDRequest) (*mcpsdk.CallToolResult, dto.ActionResponse, error) {
				return h.pin(ctx, nil, req)
			},
		},
		{
			name: "unpin",
			expect: func(rec *MockusecaseMockRecorder, res dto.ActionResponse, err error) {
				rec.Unpin(gomock.Any(), req).Return(res, err)
			},
			call: func(h *Handler, ctx context.Context, req dto.ListIDRequest) (*mcpsdk.CallToolResult, dto.ActionResponse, error) {
				return h.unpin(ctx, nil, req)
			},
		},
	}

	outcomes := []struct {
		name string
		res  dto.ActionResponse
		err  error

		wantIsError bool
		wantText    string
	}{
		{name: "success", res: okResponse},
		{
			name:        "list not found",
			err:         x.ErrListNotFound,
			wantIsError: true,
			wantText:    "list not found",
		},
		{
			name:        "invalid id",
			err:         x.ErrInvalidListID,
			wantIsError: true,
			wantText:    "invalid id",
		},
		{
			name:        "forbidden",
			err:         x.ErrForbidden,
			wantIsError: true,
			wantText:    "not authorized for this resource: it may be protected, deleted, or owned by somebody else",
		},
		{
			name:        "unknown error",
			err:         errors.New("boom"),
			wantIsError: true,
			wantText:    "internal server error",
		},
	}

	for _, tool := range tools {
		for _, c := range outcomes {
			t.Run(tool.name+": "+c.name, func(t *testing.T) {
				ctrl := gomock.NewController(t)
				usecase := NewMockusecase(ctrl)
				tool.expect(usecase.EXPECT(), c.res, c.err)

				res, out, err := tool.call(New(usecase), context.Background(), req)
				require.NoError(t, err)

				if c.wantIsError {
					assert.True(t, res.IsError)
					assert.Equal(t, c.wantText, resultText(t, res))
					assert.Equal(t, dto.ActionResponse{}, out)

					return
				}

				assert.Nil(t, res)
				assert.Equal(t, c.res, out)
			})
		}
	}
}

func TestPinned(t *testing.T) {
	req := dto.SelfPageRequest{MaxResults: 10}

	okResponse := dto.ListsResponse{Lists: []dto.List{{ID: "1", Name: "reading"}}, Count: 1}

	cases := []struct {
		name string
		res  dto.ListsResponse
		err  error

		wantIsError bool
		wantText    string
	}{
		{name: "success", res: okResponse},
		{name: "success with no pinned lists", res: dto.ListsResponse{}},
		{
			name:        "unauthorized",
			err:         x.ErrUnauthorized,
			wantIsError: true,
			wantText:    "X rejected the credentials; check the four keys in the server configuration",
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
			usecase.EXPECT().Pinned(gomock.Any(), req).Return(c.res, c.err)

			res, out, err := New(usecase).pinned(context.Background(), nil, req)
			require.NoError(t, err)

			if c.wantIsError {
				assert.True(t, res.IsError)
				assert.Equal(t, c.wantText, resultText(t, res))
				assert.Equal(t, dto.ListsResponse{}, out)

				return
			}

			assert.Nil(t, res)
			assert.Equal(t, c.res, out)
		})
	}
}
