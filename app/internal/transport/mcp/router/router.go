package router

import (
	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"

	mcpLists "github.com/Role1776/mcp-x/app/internal/transport/mcp/x/lists"
	mcpMedia "github.com/Role1776/mcp-x/app/internal/transport/mcp/x/media"
	mcpPosts "github.com/Role1776/mcp-x/app/internal/transport/mcp/x/posts"
	mcpUsers "github.com/Role1776/mcp-x/app/internal/transport/mcp/x/users"
)

type Router struct {
	mcp   *mcpsdk.Server
	posts *mcpPosts.Handler
	users *mcpUsers.Handler
	lists *mcpLists.Handler
	media *mcpMedia.Handler
}

func New(
	mcpServer *mcpsdk.Server,
	posts *mcpPosts.Handler,
	users *mcpUsers.Handler,
	lists *mcpLists.Handler,
	media *mcpMedia.Handler,
) *Router {
	return &Router{
		mcp:   mcpServer,
		posts: posts,
		users: users,
		lists: lists,
		media: media,
	}
}

func (r *Router) RegisterAllTools() {
	r.posts.RegisterTools(r.mcp)
	r.users.RegisterTools(r.mcp)
	r.lists.RegisterTools(r.mcp)
	r.media.RegisterTools(r.mcp)
}
