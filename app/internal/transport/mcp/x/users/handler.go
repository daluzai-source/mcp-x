package users

import (
	"context"

	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"

	dto "github.com/Role1776/mcp-x/app/internal/dto/x/users"
)

const (
	toolLookup = "x_users_lookup"
	toolMe     = "x_users_me"

	toolFollowers = "x_user_followers"
	toolFollowing = "x_user_following"
	toolFollow    = "x_user_follow"
	toolUnfollow  = "x_user_unfollow"

	toolBlocked = "x_users_blocked"

	toolMute   = "x_user_mute"
	toolUnmute = "x_user_unmute"
	toolMuted  = "x_users_muted"

	selfNote = "Acts as the account whose API keys this server runs on."
)

var destructiveHint = true

var (
	readOnly    = &mcpsdk.ToolAnnotations{ReadOnlyHint: true}
	writing     = &mcpsdk.ToolAnnotations{IdempotentHint: false}
	destructive = &mcpsdk.ToolAnnotations{DestructiveHint: &destructiveHint, IdempotentHint: false}
)

//go:generate mockgen -source=handler.go -destination=usecase_mock_test.go -package=users
type usecase interface {
	Lookup(ctx context.Context, req dto.LookupRequest) (dto.LookupResponse, error)
	Me(ctx context.Context, req dto.MeRequest) (dto.MeResponse, error)
	Followers(ctx context.Context, req dto.ListRequest) (dto.UsersResponse, error)
	Following(ctx context.Context, req dto.ListRequest) (dto.UsersResponse, error)
	Follow(ctx context.Context, req dto.TargetRequest) (dto.ActionResponse, error)
	Unfollow(ctx context.Context, req dto.TargetRequest) (dto.ActionResponse, error)
	Blocked(ctx context.Context, req dto.SelfListRequest) (dto.UsersResponse, error)
	Mute(ctx context.Context, req dto.TargetRequest) (dto.ActionResponse, error)
	Unmute(ctx context.Context, req dto.TargetRequest) (dto.ActionResponse, error)
	Muted(ctx context.Context, req dto.SelfListRequest) (dto.UsersResponse, error)
}

type Handler struct {
	usecase usecase
}

func New(usecase usecase) *Handler {
	return &Handler{usecase: usecase}
}

func (h *Handler) RegisterTools(s *mcpsdk.Server) {
	mcpsdk.AddTool(s, &mcpsdk.Tool{
		Name:        toolLookup,
		Title:       "Look up X users",
		Annotations: readOnly,
		Description: "Returns profiles for up to 100 usernames or up to 100 numeric ids in a single call — " +
			"pass either usernames or ids, not both, and batch them instead of calling once per user. " +
			"Usernames that do not resolve come back in 'not_found'.",
		OutputSchema: outputFor[dto.LookupResponse](),
	}, h.lookup)

	mcpsdk.AddTool(s, &mcpsdk.Tool{
		Name:         toolMe,
		Title:        "The X account behind these keys",
		Annotations:  readOnly,
		Description:  "Returns the profile of the account whose API keys this server runs on.",
		OutputSchema: outputFor[dto.MeResponse](),
	}, h.me)

	mcpsdk.AddTool(s, &mcpsdk.Tool{
		Name:        toolFollowers,
		Title:       "X followers of a user",
		Annotations: readOnly,
		Description: "Returns the followers of the given username, one page at a time. " +
			"Pass next_token back as pagination_token to walk further.",
		OutputSchema: outputFor[dto.UsersResponse](),
	}, h.followers)

	mcpsdk.AddTool(s, &mcpsdk.Tool{
		Name:        toolFollowing,
		Title:       "Accounts an X user follows",
		Annotations: readOnly,
		Description: "Returns the accounts the given username follows, one page at a time. " +
			"Pass next_token back as pagination_token to walk further.",
		OutputSchema: outputFor[dto.UsersResponse](),
	}, h.following)

	mcpsdk.AddTool(s, &mcpsdk.Tool{
		Name:         toolFollow,
		Title:        "Follow an X user",
		Annotations:  writing,
		Description:  "Follows the given username. The follow is public. " + selfNote,
		OutputSchema: outputFor[dto.ActionResponse](),
	}, h.follow)

	mcpsdk.AddTool(s, &mcpsdk.Tool{
		Name:         toolUnfollow,
		Title:        "Unfollow an X user",
		Annotations:  destructive,
		Description:  "Stops following the given username. " + selfNote,
		OutputSchema: outputFor[dto.ActionResponse](),
	}, h.unfollow)

	mcpsdk.AddTool(s, &mcpsdk.Tool{
		Name:         toolBlocked,
		Title:        "Accounts you block on X",
		Annotations:  readOnly,
		Description:  "Returns the accounts blocked by " + "the account whose API keys this server runs on.",
		OutputSchema: outputFor[dto.UsersResponse](),
	}, h.blocked)

	mcpsdk.AddTool(s, &mcpsdk.Tool{
		Name:        toolMute,
		Title:       "Mute an X user",
		Annotations: writing,
		Description: "Mutes the given username: their posts disappear from your timeline, but they are not " +
			"told and can still see yours. " + selfNote,
		OutputSchema: outputFor[dto.ActionResponse](),
	}, h.mute)

	mcpsdk.AddTool(s, &mcpsdk.Tool{
		Name:         toolUnmute,
		Title:        "Unmute an X user",
		Annotations:  writing,
		Description:  "Removes a mute. " + selfNote,
		OutputSchema: outputFor[dto.ActionResponse](),
	}, h.unmute)

	mcpsdk.AddTool(s, &mcpsdk.Tool{
		Name:         toolMuted,
		Title:        "Accounts you mute on X",
		Annotations:  readOnly,
		Description:  "Returns the accounts muted by the account whose API keys this server runs on.",
		OutputSchema: outputFor[dto.UsersResponse](),
	}, h.muted)
}
