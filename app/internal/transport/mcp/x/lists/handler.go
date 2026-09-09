package lists

import (
	"context"

	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"

	dto "github.com/Role1776/mcp-x/app/internal/dto/x/lists"
)

const (
	toolCreate = "x_list_create"
	toolUpdate = "x_list_update"
	toolDelete = "x_list_delete"
	toolGet    = "x_list_get"
	toolOwned  = "x_lists_owned"

	toolPosts        = "x_list_posts"
	toolMembers      = "x_list_members"
	toolMemberAdd    = "x_list_member_add"
	toolMemberRemove = "x_list_member_remove"

	toolFollow   = "x_list_follow"
	toolUnfollow = "x_list_unfollow"

	toolPin    = "x_list_pin"
	toolUnpin  = "x_list_unpin"
	toolPinned = "x_lists_pinned"

	selfNote = "Acts as the account whose API keys this server runs on."
)

var destructiveHint = true

var (
	readOnly    = &mcpsdk.ToolAnnotations{ReadOnlyHint: true}
	writing     = &mcpsdk.ToolAnnotations{IdempotentHint: false}
	destructive = &mcpsdk.ToolAnnotations{DestructiveHint: &destructiveHint, IdempotentHint: false}
)

//go:generate mockgen -source=handler.go -destination=usecase_mock_test.go -package=lists
type usecase interface {
	Create(ctx context.Context, req dto.CreateRequest) (dto.ActionResponse, error)
	Update(ctx context.Context, req dto.UpdateRequest) (dto.ActionResponse, error)
	Delete(ctx context.Context, req dto.ListIDRequest) (dto.ActionResponse, error)
	Get(ctx context.Context, req dto.ListIDRequest) (dto.ListResponse, error)
	Owned(ctx context.Context, req dto.OwnedRequest) (dto.ListsResponse, error)
	Posts(ctx context.Context, req dto.ListPageRequest) (dto.PostsResponse, error)
	Members(ctx context.Context, req dto.ListPageRequest) (dto.UsersResponse, error)
	MemberAdd(ctx context.Context, req dto.MemberRequest) (dto.ActionResponse, error)
	MemberRemove(ctx context.Context, req dto.MemberRequest) (dto.ActionResponse, error)
	Follow(ctx context.Context, req dto.ListIDRequest) (dto.ActionResponse, error)
	Unfollow(ctx context.Context, req dto.ListIDRequest) (dto.ActionResponse, error)
	Pin(ctx context.Context, req dto.ListIDRequest) (dto.ActionResponse, error)
	Unpin(ctx context.Context, req dto.ListIDRequest) (dto.ActionResponse, error)
	Pinned(ctx context.Context, req dto.SelfPageRequest) (dto.ListsResponse, error)
}

type Handler struct {
	usecase usecase
}

func New(usecase usecase) *Handler {
	return &Handler{usecase: usecase}
}

func (h *Handler) RegisterTools(s *mcpsdk.Server) {
	mcpsdk.AddTool(s, &mcpsdk.Tool{
		Name:        toolCreate,
		Title:       "Create an X list",
		Annotations: writing,
		Description: "Creates a list owned by the account whose API keys this server runs on. " +
			"The name is capped at 25 characters and the description at 100. A private list is visible " +
			"only to its owner.",
		OutputSchema: outputFor[dto.ActionResponse](),
	}, h.create)

	mcpsdk.AddTool(s, &mcpsdk.Tool{
		Name:         toolUpdate,
		Title:        "Update an X list",
		Annotations:  writing,
		Description:  "Changes the name, description or privacy of a list you own. Send only the fields you want changed.",
		OutputSchema: outputFor[dto.ActionResponse](),
	}, h.update)

	mcpsdk.AddTool(s, &mcpsdk.Tool{
		Name:        toolDelete,
		Title:       "Delete an X list",
		Annotations: destructive,
		Description: "Deletes a list you own. This is irreversible: the list, its membership and its " +
			"followers are gone and cannot be restored.",
		OutputSchema: outputFor[dto.ActionResponse](),
	}, h.delete)

	mcpsdk.AddTool(s, &mcpsdk.Tool{
		Name:         toolGet,
		Title:        "Get an X list",
		Annotations:  readOnly,
		Description:  "Returns one list by id, including its member and follower counts.",
		OutputSchema: outputFor[dto.ListResponse](),
	}, h.get)

	mcpsdk.AddTool(s, &mcpsdk.Tool{
		Name:         toolOwned,
		Title:        "X lists owned by a user",
		Annotations:  readOnly,
		Description:  "Returns the lists owned by the given username, or by the key owner when username is omitted.",
		OutputSchema: outputFor[dto.ListsResponse](),
	}, h.owned)

	mcpsdk.AddTool(s, &mcpsdk.Tool{
		Name:        toolPosts,
		Title:       "Posts in an X list",
		Annotations: readOnly,
		Description: "Returns recent posts from the members of a list. Every post is a billed read, so ask " +
			"for the smallest max_results that answers the question.",
		OutputSchema: outputFor[dto.PostsResponse](),
	}, h.posts)

	mcpsdk.AddTool(s, &mcpsdk.Tool{
		Name:         toolMembers,
		Title:        "Members of an X list",
		Annotations:  readOnly,
		Description:  "Returns the accounts that belong to a list, one page at a time.",
		OutputSchema: outputFor[dto.UsersResponse](),
	}, h.members)

	mcpsdk.AddTool(s, &mcpsdk.Tool{
		Name:         toolMemberAdd,
		Title:        "Add a member to an X list",
		Annotations:  writing,
		Description:  "Adds the given username to a list you own.",
		OutputSchema: outputFor[dto.ActionResponse](),
	}, h.memberAdd)

	mcpsdk.AddTool(s, &mcpsdk.Tool{
		Name:         toolMemberRemove,
		Title:        "Remove a member from an X list",
		Annotations:  destructive,
		Description:  "Removes the given username from a list you own.",
		OutputSchema: outputFor[dto.ActionResponse](),
	}, h.memberRemove)

	mcpsdk.AddTool(s, &mcpsdk.Tool{
		Name:         toolFollow,
		Title:        "Follow an X list",
		Annotations:  writing,
		Description:  "Follows a public list. " + selfNote,
		OutputSchema: outputFor[dto.ActionResponse](),
	}, h.follow)

	mcpsdk.AddTool(s, &mcpsdk.Tool{
		Name:         toolUnfollow,
		Title:        "Unfollow an X list",
		Annotations:  destructive,
		Description:  "Stops following a list. " + selfNote,
		OutputSchema: outputFor[dto.ActionResponse](),
	}, h.unfollow)

	mcpsdk.AddTool(s, &mcpsdk.Tool{
		Name:         toolPin,
		Title:        "Pin an X list",
		Annotations:  writing,
		Description:  "Pins a list to the top of your lists view. " + selfNote,
		OutputSchema: outputFor[dto.ActionResponse](),
	}, h.pin)

	mcpsdk.AddTool(s, &mcpsdk.Tool{
		Name:         toolUnpin,
		Title:        "Unpin an X list",
		Annotations:  writing,
		Description:  "Removes a list from your pinned lists. " + selfNote,
		OutputSchema: outputFor[dto.ActionResponse](),
	}, h.unpin)

	mcpsdk.AddTool(s, &mcpsdk.Tool{
		Name:         toolPinned,
		Title:        "X lists you pinned",
		Annotations:  readOnly,
		Description:  "Returns the pinned lists of the account whose API keys this server runs on.",
		OutputSchema: outputFor[dto.ListsResponse](),
	}, h.pinned)
}
