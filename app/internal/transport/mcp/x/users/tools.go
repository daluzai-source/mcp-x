package users

import (
	"context"

	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"

	dto "github.com/Role1776/mcp-x/app/internal/dto/x/users"
)

func (h *Handler) lookup(ctx context.Context, _ *mcpsdk.CallToolRequest, req dto.LookupRequest) (*mcpsdk.CallToolResult, dto.LookupResponse, error) {
	res, err := h.usecase.Lookup(ctx, req)
	if err != nil {
		return result(err), dto.LookupResponse{}, nil
	}

	return nil, res, nil
}

func (h *Handler) me(ctx context.Context, _ *mcpsdk.CallToolRequest, req dto.MeRequest) (*mcpsdk.CallToolResult, dto.MeResponse, error) {
	res, err := h.usecase.Me(ctx, req)
	if err != nil {
		return result(err), dto.MeResponse{}, nil
	}

	return nil, res, nil
}

func (h *Handler) followers(ctx context.Context, _ *mcpsdk.CallToolRequest, req dto.ListRequest) (*mcpsdk.CallToolResult, dto.UsersResponse, error) {
	return h.list(ctx, req, h.usecase.Followers)
}

func (h *Handler) following(ctx context.Context, _ *mcpsdk.CallToolRequest, req dto.ListRequest) (*mcpsdk.CallToolResult, dto.UsersResponse, error) {
	return h.list(ctx, req, h.usecase.Following)
}

func (h *Handler) follow(ctx context.Context, _ *mcpsdk.CallToolRequest, req dto.TargetRequest) (*mcpsdk.CallToolResult, dto.ActionResponse, error) {
	return h.action(ctx, req, h.usecase.Follow)
}

func (h *Handler) unfollow(ctx context.Context, _ *mcpsdk.CallToolRequest, req dto.TargetRequest) (*mcpsdk.CallToolResult, dto.ActionResponse, error) {
	return h.action(ctx, req, h.usecase.Unfollow)
}

func (h *Handler) mute(ctx context.Context, _ *mcpsdk.CallToolRequest, req dto.TargetRequest) (*mcpsdk.CallToolResult, dto.ActionResponse, error) {
	return h.action(ctx, req, h.usecase.Mute)
}

func (h *Handler) unmute(ctx context.Context, _ *mcpsdk.CallToolRequest, req dto.TargetRequest) (*mcpsdk.CallToolResult, dto.ActionResponse, error) {
	return h.action(ctx, req, h.usecase.Unmute)
}

func (h *Handler) blocked(ctx context.Context, _ *mcpsdk.CallToolRequest, req dto.SelfListRequest) (*mcpsdk.CallToolResult, dto.UsersResponse, error) {
	return h.selfList(ctx, req, h.usecase.Blocked)
}

func (h *Handler) muted(ctx context.Context, _ *mcpsdk.CallToolRequest, req dto.SelfListRequest) (*mcpsdk.CallToolResult, dto.UsersResponse, error) {
	return h.selfList(ctx, req, h.usecase.Muted)
}

func (h *Handler) list(
	ctx context.Context,
	req dto.ListRequest,
	call func(context.Context, dto.ListRequest) (dto.UsersResponse, error),
) (*mcpsdk.CallToolResult, dto.UsersResponse, error) {
	res, err := call(ctx, req)
	if err != nil {
		return result(err), dto.UsersResponse{}, nil
	}

	return nil, res, nil
}

func (h *Handler) selfList(
	ctx context.Context,
	req dto.SelfListRequest,
	call func(context.Context, dto.SelfListRequest) (dto.UsersResponse, error),
) (*mcpsdk.CallToolResult, dto.UsersResponse, error) {
	res, err := call(ctx, req)
	if err != nil {
		return result(err), dto.UsersResponse{}, nil
	}

	return nil, res, nil
}

func (h *Handler) action(
	ctx context.Context,
	req dto.TargetRequest,
	call func(context.Context, dto.TargetRequest) (dto.ActionResponse, error),
) (*mcpsdk.CallToolResult, dto.ActionResponse, error) {
	res, err := call(ctx, req)
	if err != nil {
		return result(err), dto.ActionResponse{}, nil
	}

	return nil, res, nil
}
