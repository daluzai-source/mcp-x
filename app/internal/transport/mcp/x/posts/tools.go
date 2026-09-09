package posts

import (
	"context"

	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"

	dto "github.com/Role1776/mcp-x/app/internal/dto/x/posts"
)

func (h *Handler) search(ctx context.Context, _ *mcpsdk.CallToolRequest, req dto.SearchRequest) (*mcpsdk.CallToolResult, dto.SearchResponse, error) {
	res, err := h.usecase.Search(ctx, req)
	if err != nil {
		return result(err), dto.SearchResponse{}, nil
	}

	return nil, res, nil
}

func (h *Handler) count(ctx context.Context, _ *mcpsdk.CallToolRequest, req dto.CountRequest) (*mcpsdk.CallToolResult, dto.CountResponse, error) {
	res, err := h.usecase.Count(ctx, req)
	if err != nil {
		return result(err), dto.CountResponse{}, nil
	}

	return nil, res, nil
}

func (h *Handler) lookup(ctx context.Context, _ *mcpsdk.CallToolRequest, req dto.LookupRequest) (*mcpsdk.CallToolResult, dto.LookupResponse, error) {
	res, err := h.usecase.Lookup(ctx, req)
	if err != nil {
		return result(err), dto.LookupResponse{}, nil
	}

	return nil, res, nil
}

func (h *Handler) byUser(ctx context.Context, _ *mcpsdk.CallToolRequest, req dto.ByUserRequest) (*mcpsdk.CallToolResult, dto.ByUserResponse, error) {
	res, err := h.usecase.ByUser(ctx, req)
	if err != nil {
		return result(err), dto.ByUserResponse{}, nil
	}

	return nil, res, nil
}

func (h *Handler) mentions(ctx context.Context, _ *mcpsdk.CallToolRequest, req dto.TimelineRequest) (*mcpsdk.CallToolResult, dto.PostsResponse, error) {
	res, err := h.usecase.Mentions(ctx, req)
	if err != nil {
		return result(err), dto.PostsResponse{}, nil
	}

	return nil, res, nil
}

func (h *Handler) home(ctx context.Context, _ *mcpsdk.CallToolRequest, req dto.TimelineRequest) (*mcpsdk.CallToolResult, dto.PostsResponse, error) {
	res, err := h.usecase.Home(ctx, req)
	if err != nil {
		return result(err), dto.PostsResponse{}, nil
	}

	return nil, res, nil
}

func (h *Handler) quotes(ctx context.Context, _ *mcpsdk.CallToolRequest, req dto.QuotesRequest) (*mcpsdk.CallToolResult, dto.PostsResponse, error) {
	res, err := h.usecase.Quotes(ctx, req)
	if err != nil {
		return result(err), dto.PostsResponse{}, nil
	}

	return nil, res, nil
}

func (h *Handler) postsLiked(ctx context.Context, _ *mcpsdk.CallToolRequest, req dto.TimelineRequest) (*mcpsdk.CallToolResult, dto.PostsResponse, error) {
	res, err := h.usecase.Liked(ctx, req)
	if err != nil {
		return result(err), dto.PostsResponse{}, nil
	}

	return nil, res, nil
}

func (h *Handler) bookmarksList(ctx context.Context, _ *mcpsdk.CallToolRequest, req dto.TimelineRequest) (*mcpsdk.CallToolResult, dto.PostsResponse, error) {
	res, err := h.usecase.Bookmarks(ctx, req)
	if err != nil {
		return result(err), dto.PostsResponse{}, nil
	}

	return nil, res, nil
}

func (h *Handler) likedBy(ctx context.Context, _ *mcpsdk.CallToolRequest, req dto.PostUsersRequest) (*mcpsdk.CallToolResult, dto.UsersResponse, error) {
	res, err := h.usecase.LikedBy(ctx, req)
	if err != nil {
		return result(err), dto.UsersResponse{}, nil
	}

	return nil, res, nil
}

func (h *Handler) repostedBy(ctx context.Context, _ *mcpsdk.CallToolRequest, req dto.PostUsersRequest) (*mcpsdk.CallToolResult, dto.UsersResponse, error) {
	res, err := h.usecase.RepostedBy(ctx, req)
	if err != nil {
		return result(err), dto.UsersResponse{}, nil
	}

	return nil, res, nil
}

func (h *Handler) create(ctx context.Context, _ *mcpsdk.CallToolRequest, req dto.CreateRequest) (*mcpsdk.CallToolResult, dto.CreateResponse, error) {
	res, err := h.usecase.Create(ctx, req)
	if err != nil {
		return result(err), dto.CreateResponse{}, nil
	}

	return nil, res, nil
}

func (h *Handler) delete(ctx context.Context, _ *mcpsdk.CallToolRequest, req dto.PostIDRequest) (*mcpsdk.CallToolResult, dto.ActionResponse, error) {
	return h.action(ctx, req, h.usecase.Delete)
}

func (h *Handler) like(ctx context.Context, _ *mcpsdk.CallToolRequest, req dto.PostIDRequest) (*mcpsdk.CallToolResult, dto.ActionResponse, error) {
	return h.action(ctx, req, h.usecase.Like)
}

func (h *Handler) unlike(ctx context.Context, _ *mcpsdk.CallToolRequest, req dto.PostIDRequest) (*mcpsdk.CallToolResult, dto.ActionResponse, error) {
	return h.action(ctx, req, h.usecase.Unlike)
}

func (h *Handler) repost(ctx context.Context, _ *mcpsdk.CallToolRequest, req dto.PostIDRequest) (*mcpsdk.CallToolResult, dto.ActionResponse, error) {
	return h.action(ctx, req, h.usecase.Repost)
}

func (h *Handler) unrepost(ctx context.Context, _ *mcpsdk.CallToolRequest, req dto.PostIDRequest) (*mcpsdk.CallToolResult, dto.ActionResponse, error) {
	return h.action(ctx, req, h.usecase.Unrepost)
}

func (h *Handler) action(
	ctx context.Context,
	req dto.PostIDRequest,
	call func(context.Context, dto.PostIDRequest) (dto.ActionResponse, error),
) (*mcpsdk.CallToolResult, dto.ActionResponse, error) {
	res, err := call(ctx, req)
	if err != nil {
		return result(err), dto.ActionResponse{}, nil
	}

	return nil, res, nil
}
