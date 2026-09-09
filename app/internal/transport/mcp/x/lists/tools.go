package lists

import (
	"context"

	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"

	dto "github.com/Role1776/mcp-x/app/internal/dto/x/lists"
)

func (h *Handler) create(ctx context.Context, _ *mcpsdk.CallToolRequest, req dto.CreateRequest) (*mcpsdk.CallToolResult, dto.ActionResponse, error) {
	res, err := h.usecase.Create(ctx, req)
	if err != nil {
		return result(err), dto.ActionResponse{}, nil
	}

	return nil, res, nil
}

func (h *Handler) update(ctx context.Context, _ *mcpsdk.CallToolRequest, req dto.UpdateRequest) (*mcpsdk.CallToolResult, dto.ActionResponse, error) {
	res, err := h.usecase.Update(ctx, req)
	if err != nil {
		return result(err), dto.ActionResponse{}, nil
	}

	return nil, res, nil
}

func (h *Handler) get(ctx context.Context, _ *mcpsdk.CallToolRequest, req dto.ListIDRequest) (*mcpsdk.CallToolResult, dto.ListResponse, error) {
	res, err := h.usecase.Get(ctx, req)
	if err != nil {
		return result(err), dto.ListResponse{}, nil
	}

	return nil, res, nil
}

func (h *Handler) owned(ctx context.Context, _ *mcpsdk.CallToolRequest, req dto.OwnedRequest) (*mcpsdk.CallToolResult, dto.ListsResponse, error) {
	res, err := h.usecase.Owned(ctx, req)
	if err != nil {
		return result(err), dto.ListsResponse{}, nil
	}

	return nil, res, nil
}

func (h *Handler) posts(ctx context.Context, _ *mcpsdk.CallToolRequest, req dto.ListPageRequest) (*mcpsdk.CallToolResult, dto.PostsResponse, error) {
	res, err := h.usecase.Posts(ctx, req)
	if err != nil {
		return result(err), dto.PostsResponse{}, nil
	}

	return nil, res, nil
}

func (h *Handler) members(ctx context.Context, _ *mcpsdk.CallToolRequest, req dto.ListPageRequest) (*mcpsdk.CallToolResult, dto.UsersResponse, error) {
	res, err := h.usecase.Members(ctx, req)
	if err != nil {
		return result(err), dto.UsersResponse{}, nil
	}

	return nil, res, nil
}

func (h *Handler) memberAdd(ctx context.Context, _ *mcpsdk.CallToolRequest, req dto.MemberRequest) (*mcpsdk.CallToolResult, dto.ActionResponse, error) {
	return h.member(ctx, req, h.usecase.MemberAdd)
}

func (h *Handler) memberRemove(ctx context.Context, _ *mcpsdk.CallToolRequest, req dto.MemberRequest) (*mcpsdk.CallToolResult, dto.ActionResponse, error) {
	return h.member(ctx, req, h.usecase.MemberRemove)
}

func (h *Handler) delete(ctx context.Context, _ *mcpsdk.CallToolRequest, req dto.ListIDRequest) (*mcpsdk.CallToolResult, dto.ActionResponse, error) {
	return h.action(ctx, req, h.usecase.Delete)
}

func (h *Handler) follow(ctx context.Context, _ *mcpsdk.CallToolRequest, req dto.ListIDRequest) (*mcpsdk.CallToolResult, dto.ActionResponse, error) {
	return h.action(ctx, req, h.usecase.Follow)
}

func (h *Handler) unfollow(ctx context.Context, _ *mcpsdk.CallToolRequest, req dto.ListIDRequest) (*mcpsdk.CallToolResult, dto.ActionResponse, error) {
	return h.action(ctx, req, h.usecase.Unfollow)
}

func (h *Handler) pin(ctx context.Context, _ *mcpsdk.CallToolRequest, req dto.ListIDRequest) (*mcpsdk.CallToolResult, dto.ActionResponse, error) {
	return h.action(ctx, req, h.usecase.Pin)
}

func (h *Handler) unpin(ctx context.Context, _ *mcpsdk.CallToolRequest, req dto.ListIDRequest) (*mcpsdk.CallToolResult, dto.ActionResponse, error) {
	return h.action(ctx, req, h.usecase.Unpin)
}

func (h *Handler) pinned(ctx context.Context, _ *mcpsdk.CallToolRequest, req dto.SelfPageRequest) (*mcpsdk.CallToolResult, dto.ListsResponse, error) {
	return h.selfLists(ctx, req, h.usecase.Pinned)
}

func (h *Handler) action(
	ctx context.Context,
	req dto.ListIDRequest,
	call func(context.Context, dto.ListIDRequest) (dto.ActionResponse, error),
) (*mcpsdk.CallToolResult, dto.ActionResponse, error) {
	res, err := call(ctx, req)
	if err != nil {
		return result(err), dto.ActionResponse{}, nil
	}

	return nil, res, nil
}

func (h *Handler) member(
	ctx context.Context,
	req dto.MemberRequest,
	call func(context.Context, dto.MemberRequest) (dto.ActionResponse, error),
) (*mcpsdk.CallToolResult, dto.ActionResponse, error) {
	res, err := call(ctx, req)
	if err != nil {
		return result(err), dto.ActionResponse{}, nil
	}

	return nil, res, nil
}

func (h *Handler) selfLists(
	ctx context.Context,
	req dto.SelfPageRequest,
	call func(context.Context, dto.SelfPageRequest) (dto.ListsResponse, error),
) (*mcpsdk.CallToolResult, dto.ListsResponse, error) {
	res, err := call(ctx, req)
	if err != nil {
		return result(err), dto.ListsResponse{}, nil
	}

	return nil, res, nil
}
