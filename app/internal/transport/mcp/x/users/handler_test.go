package users

import (
	"context"
	"errors"
	"testing"

	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/Role1776/mcp-x/app/internal/domain/x"
	dto "github.com/Role1776/mcp-x/app/internal/dto/x/users"
)

func TestLookup(t *testing.T) {
	req := dto.LookupRequest{Usernames: []string{"jack"}}

	okResponse := dto.LookupResponse{
		Users:    []dto.User{{ID: "1", Username: "jack", URL: "https://x.com/jack"}},
		Count:    1,
		NotFound: []string{"jill"},
	}

	cases := []struct {
		name string
		res  dto.LookupResponse
		err  error

		wantIsError bool
		wantText    string
		wantRes     dto.LookupResponse
	}{
		{name: "success", res: okResponse, wantRes: okResponse},
		{name: "success with empty results", res: dto.LookupResponse{}, wantRes: dto.LookupResponse{}},
		{
			name:        "too many users",
			err:         x.ErrTooManyUsers,
			wantIsError: true,
			wantText:    "too many users in one call",
		},
		{
			name:        "invalid username",
			err:         x.ErrInvalidUsername,
			wantIsError: true,
			wantText:    "invalid username: 1 to 15 letters, digits or underscores",
		},
		{
			name:        "invalid request",
			err:         x.ErrInvalidRequest,
			wantIsError: true,
			wantText:    "invalid request",
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
			usecase.EXPECT().Lookup(gomock.Any(), req).Return(c.res, c.err)

			res, out, err := New(usecase).lookup(context.Background(), nil, req)
			require.NoError(t, err)

			if c.wantIsError {
				assert.True(t, res.IsError)
				assert.Equal(t, c.wantText, resultText(t, res))
				assert.Equal(t, dto.LookupResponse{}, out)

				return
			}

			assert.Nil(t, res)
			assert.Equal(t, c.wantRes, out)
		})
	}
}

func TestMe(t *testing.T) {
	req := dto.MeRequest{}

	okResponse := dto.MeResponse{User: dto.User{ID: "1", Username: "jack"}}

	cases := []struct {
		name string
		res  dto.MeResponse
		err  error

		wantIsError bool
		wantText    string
	}{
		{name: "success", res: okResponse},
		{
			name:        "unauthorized",
			err:         x.ErrUnauthorized,
			wantIsError: true,
			wantText:    "X rejected the credentials; check the four keys in the server configuration",
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
			usecase.EXPECT().Me(gomock.Any(), req).Return(c.res, c.err)

			res, out, err := New(usecase).me(context.Background(), nil, req)
			require.NoError(t, err)

			if c.wantIsError {
				assert.True(t, res.IsError)
				assert.Equal(t, c.wantText, resultText(t, res))
				assert.Equal(t, dto.MeResponse{}, out)

				return
			}

			assert.Nil(t, res)
			assert.Equal(t, c.res, out)
		})
	}
}

func TestUserListTools(t *testing.T) {
	req := dto.ListRequest{Username: "jack", MaxResults: 10}

	okResponse := dto.UsersResponse{
		Users:     []dto.User{{ID: "1", Username: "jill"}},
		Count:     1,
		NextToken: "next",
	}

	tools := []struct {
		name   string
		expect func(*MockusecaseMockRecorder, dto.UsersResponse, error)
		call   func(*Handler, context.Context, dto.ListRequest) (*mcpsdk.CallToolResult, dto.UsersResponse, error)
	}{
		{
			name: "followers",
			expect: func(rec *MockusecaseMockRecorder, res dto.UsersResponse, err error) {
				rec.Followers(gomock.Any(), req).Return(res, err)
			},
			call: func(h *Handler, ctx context.Context, req dto.ListRequest) (*mcpsdk.CallToolResult, dto.UsersResponse, error) {
				return h.followers(ctx, nil, req)
			},
		},
		{
			name: "following",
			expect: func(rec *MockusecaseMockRecorder, res dto.UsersResponse, err error) {
				rec.Following(gomock.Any(), req).Return(res, err)
			},
			call: func(h *Handler, ctx context.Context, req dto.ListRequest) (*mcpsdk.CallToolResult, dto.UsersResponse, error) {
				return h.following(ctx, nil, req)
			},
		},
	}

	outcomes := []struct {
		name string
		res  dto.UsersResponse
		err  error

		wantIsError bool
		wantText    string
	}{
		{name: "success", res: okResponse},
		{name: "success with an empty page", res: dto.UsersResponse{}},
		{
			name:        "user not found",
			err:         x.ErrUserNotFound,
			wantIsError: true,
			wantText:    "user not found",
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
					assert.Equal(t, dto.UsersResponse{}, out)

					return
				}

				assert.Nil(t, res)
				assert.Equal(t, c.res, out)
			})
		}
	}
}

func TestUserSelfListTools(t *testing.T) {
	req := dto.SelfListRequest{MaxResults: 10}

	okResponse := dto.UsersResponse{Users: []dto.User{{ID: "1", Username: "jill"}}, Count: 1}

	tools := []struct {
		name   string
		expect func(*MockusecaseMockRecorder, dto.UsersResponse, error)
		call   func(*Handler, context.Context, dto.SelfListRequest) (*mcpsdk.CallToolResult, dto.UsersResponse, error)
	}{
		{
			name: "blocked",
			expect: func(rec *MockusecaseMockRecorder, res dto.UsersResponse, err error) {
				rec.Blocked(gomock.Any(), req).Return(res, err)
			},
			call: func(h *Handler, ctx context.Context, req dto.SelfListRequest) (*mcpsdk.CallToolResult, dto.UsersResponse, error) {
				return h.blocked(ctx, nil, req)
			},
		},
		{
			name: "muted",
			expect: func(rec *MockusecaseMockRecorder, res dto.UsersResponse, err error) {
				rec.Muted(gomock.Any(), req).Return(res, err)
			},
			call: func(h *Handler, ctx context.Context, req dto.SelfListRequest) (*mcpsdk.CallToolResult, dto.UsersResponse, error) {
				return h.muted(ctx, nil, req)
			},
		},
	}

	outcomes := []struct {
		name string
		res  dto.UsersResponse
		err  error

		wantIsError bool
		wantText    string
	}{
		{name: "success", res: okResponse},
		{
			name:        "read-only token",
			err:         x.ErrReadOnlyToken,
			wantIsError: true,
			wantText: "the access token is read-only: open the X Developer Console, set the app permissions " +
				"to Read and Write, then regenerate the Access Token and Secret — an existing token keeps " +
				"the permissions it was created with",
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
					assert.Equal(t, dto.UsersResponse{}, out)

					return
				}

				assert.Nil(t, res)
				assert.Equal(t, c.res, out)
			})
		}
	}
}

func TestUserTargetTools(t *testing.T) {
	req := dto.TargetRequest{Username: "jack"}

	okResponse := dto.ActionResponse{Username: "jack", UserID: "1", Status: "done"}

	tools := []struct {
		name   string
		expect func(*MockusecaseMockRecorder, dto.ActionResponse, error)
		call   func(*Handler, context.Context, dto.TargetRequest) (*mcpsdk.CallToolResult, dto.ActionResponse, error)
	}{
		{
			name: "follow",
			expect: func(rec *MockusecaseMockRecorder, res dto.ActionResponse, err error) {
				rec.Follow(gomock.Any(), req).Return(res, err)
			},
			call: func(h *Handler, ctx context.Context, req dto.TargetRequest) (*mcpsdk.CallToolResult, dto.ActionResponse, error) {
				return h.follow(ctx, nil, req)
			},
		},
		{
			name: "unfollow",
			expect: func(rec *MockusecaseMockRecorder, res dto.ActionResponse, err error) {
				rec.Unfollow(gomock.Any(), req).Return(res, err)
			},
			call: func(h *Handler, ctx context.Context, req dto.TargetRequest) (*mcpsdk.CallToolResult, dto.ActionResponse, error) {
				return h.unfollow(ctx, nil, req)
			},
		},
		{
			name: "mute",
			expect: func(rec *MockusecaseMockRecorder, res dto.ActionResponse, err error) {
				rec.Mute(gomock.Any(), req).Return(res, err)
			},
			call: func(h *Handler, ctx context.Context, req dto.TargetRequest) (*mcpsdk.CallToolResult, dto.ActionResponse, error) {
				return h.mute(ctx, nil, req)
			},
		},
		{
			name: "unmute",
			expect: func(rec *MockusecaseMockRecorder, res dto.ActionResponse, err error) {
				rec.Unmute(gomock.Any(), req).Return(res, err)
			},
			call: func(h *Handler, ctx context.Context, req dto.TargetRequest) (*mcpsdk.CallToolResult, dto.ActionResponse, error) {
				return h.unmute(ctx, nil, req)
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
