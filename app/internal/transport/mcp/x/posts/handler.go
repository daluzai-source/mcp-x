package posts

import (
	"context"

	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"

	dto "github.com/Role1776/mcp-x/app/internal/dto/x/posts"
)

const (
	toolSearch   = "x_posts_search"
	toolCount    = "x_posts_count"
	toolLookup   = "x_posts_lookup"
	toolByUser   = "x_posts_by_user"
	toolMentions = "x_posts_mentions"
	toolHome     = "x_posts_home"
	toolQuotes   = "x_posts_quotes"

	toolCreate = "x_post_create"
	toolDelete = "x_post_delete"

	toolLike       = "x_post_like"
	toolUnlike     = "x_post_unlike"
	toolLikedBy    = "x_post_liked_by"
	toolPostsLiked = "x_posts_liked"

	toolRepost     = "x_post_repost"
	toolUnrepost   = "x_post_unrepost"
	toolRepostedBy = "x_post_reposted_by"

	toolBookmarksList = "x_bookmarks_list"

	costNote = "Every post this returns is a billed read; ask for the smallest max_results that answers " +
		"the question and batch inputs into one call."
)

//go:generate mockgen -source=handler.go -destination=usecase_mock_test.go -package=posts
type usecase interface {
	Search(ctx context.Context, req dto.SearchRequest) (dto.SearchResponse, error)
	Count(ctx context.Context, req dto.CountRequest) (dto.CountResponse, error)
	Lookup(ctx context.Context, req dto.LookupRequest) (dto.LookupResponse, error)
	ByUser(ctx context.Context, req dto.ByUserRequest) (dto.ByUserResponse, error)
	Mentions(ctx context.Context, req dto.TimelineRequest) (dto.PostsResponse, error)
	Home(ctx context.Context, req dto.TimelineRequest) (dto.PostsResponse, error)
	Quotes(ctx context.Context, req dto.QuotesRequest) (dto.PostsResponse, error)
	Liked(ctx context.Context, req dto.TimelineRequest) (dto.PostsResponse, error)
	Bookmarks(ctx context.Context, req dto.TimelineRequest) (dto.PostsResponse, error)
	LikedBy(ctx context.Context, req dto.PostUsersRequest) (dto.UsersResponse, error)
	RepostedBy(ctx context.Context, req dto.PostUsersRequest) (dto.UsersResponse, error)
	Create(ctx context.Context, req dto.CreateRequest) (dto.CreateResponse, error)
	Delete(ctx context.Context, req dto.PostIDRequest) (dto.ActionResponse, error)
	Like(ctx context.Context, req dto.PostIDRequest) (dto.ActionResponse, error)
	Unlike(ctx context.Context, req dto.PostIDRequest) (dto.ActionResponse, error)
	Repost(ctx context.Context, req dto.PostIDRequest) (dto.ActionResponse, error)
	Unrepost(ctx context.Context, req dto.PostIDRequest) (dto.ActionResponse, error)
}

type Handler struct {
	usecase usecase
}

func New(usecase usecase) *Handler {
	return &Handler{usecase: usecase}
}

var destructiveHint = true

var (
	readOnly    = &mcpsdk.ToolAnnotations{ReadOnlyHint: true}
	writing     = &mcpsdk.ToolAnnotations{IdempotentHint: false}
	destructive = &mcpsdk.ToolAnnotations{DestructiveHint: &destructiveHint, IdempotentHint: false}
)

func (h *Handler) RegisterTools(s *mcpsdk.Server) {
	mcpsdk.AddTool(s, &mcpsdk.Tool{
		Name:        toolSearch,
		Title:       "Search X posts",
		Annotations: readOnly,
		Description: "Searches recent X posts for one or more queries, run in parallel. " +
			"Recent search only reaches back 7 days — anything older is not available through this tool.\n\n" +
			"Query operators go inside the query string: a space means AND, 'OR' is explicit, a leading '-' " +
			"negates, and parentheses group. 'from:user' and 'to:user' filter by author and recipient, " +
			"'@user' and '#tag' match mentions and hashtags, double quotes force an exact phrase, " +
			"'lang:en' filters by language, 'is:retweet', 'is:reply', 'is:quote' and 'is:verified' filter by " +
			"post kind, 'has:media', 'has:images', 'has:videos' and 'has:links' filter by attachment, " +
			"'url:example.com' matches linked domains and 'conversation_id:123' pulls one thread.\n\n" +
			"Retweets are excluded by default; set include_retweets to keep them. " + costNote,
		OutputSchema: outputFor[dto.SearchResponse](),
	}, h.search)

	mcpsdk.AddTool(s, &mcpsdk.Tool{
		Name:        toolCount,
		Title:       "Count X posts",
		Annotations: readOnly,
		Description: "Counts how many posts match each query, bucketed by minute, hour or day, without " +
			"returning the posts themselves. This does not consume the post-read budget, so use it first " +
			"to size a topic before paying for " + toolSearch + ". Same operators and same 7-day window.",
		OutputSchema: outputFor[dto.CountResponse](),
	}, h.count)

	mcpsdk.AddTool(s, &mcpsdk.Tool{
		Name:        toolLookup,
		Title:       "Fetch X posts by id",
		Annotations: readOnly,
		Description: "Fetches posts by id, up to 100 in a single call — batch every id you need instead of " +
			"calling this once per post. Ids that are deleted, private or otherwise unavailable come back " +
			"with status 'not found' next to the ones that resolved, so a partial answer is normal. " + costNote,
		OutputSchema: outputFor[dto.LookupResponse](),
	}, h.lookup)

	mcpsdk.AddTool(s, &mcpsdk.Tool{
		Name:        toolByUser,
		Title:       "X posts by user",
		Annotations: readOnly,
		Description: "Returns recent posts from each of the given usernames, fetched in parallel. " +
			"A username that does not exist comes back with status 'not found' rather than failing the " +
			"whole call. Use 'exclude' to drop replies or retweets. " + costNote,
		OutputSchema: outputFor[dto.ByUserResponse](),
	}, h.byUser)

	mcpsdk.AddTool(s, &mcpsdk.Tool{
		Name:         toolMentions,
		Title:        "X mentions",
		Annotations:  readOnly,
		Description:  "Returns posts that mention the account whose API keys this server runs on. " + costNote,
		OutputSchema: outputFor[dto.PostsResponse](),
	}, h.mentions)

	mcpsdk.AddTool(s, &mcpsdk.Tool{
		Name:        toolHome,
		Title:       "X home timeline",
		Annotations: readOnly,
		Description: "Returns the home timeline of the account whose API keys this server runs on: posts " +
			"from the accounts it follows, newest first. " + costNote,
		OutputSchema: outputFor[dto.PostsResponse](),
	}, h.home)

	mcpsdk.AddTool(s, &mcpsdk.Tool{
		Name:         toolQuotes,
		Title:        "X quotes of a post",
		Annotations:  readOnly,
		Description:  "Returns posts that quote the given post. " + costNote,
		OutputSchema: outputFor[dto.PostsResponse](),
	}, h.quotes)

	mcpsdk.AddTool(s, &mcpsdk.Tool{
		Name:         toolPostsLiked,
		Title:        "X posts I liked",
		Annotations:  readOnly,
		Description:  "Returns posts liked by the account whose API keys this server runs on. " + costNote,
		OutputSchema: outputFor[dto.PostsResponse](),
	}, h.postsLiked)

	mcpsdk.AddTool(s, &mcpsdk.Tool{
		Name:         toolBookmarksList,
		Title:        "X bookmarks",
		Annotations:  readOnly,
		Description:  "Returns the bookmarks of the account whose API keys this server runs on. " + costNote,
		OutputSchema: outputFor[dto.PostsResponse](),
	}, h.bookmarksList)

	mcpsdk.AddTool(s, &mcpsdk.Tool{
		Name:         toolLikedBy,
		Title:        "Who liked an X post",
		Annotations:  readOnly,
		Description:  "Returns the users who liked the given post.",
		OutputSchema: outputFor[dto.UsersResponse](),
	}, h.likedBy)

	mcpsdk.AddTool(s, &mcpsdk.Tool{
		Name:         toolRepostedBy,
		Title:        "Who reposted an X post",
		Annotations:  readOnly,
		Description:  "Returns the users who reposted the given post.",
		OutputSchema: outputFor[dto.UsersResponse](),
	}, h.repostedBy)

	mcpsdk.AddTool(s, &mcpsdk.Tool{
		Name:        toolCreate,
		Title:       "Publish an X post",
		Annotations: writing,
		Description: "Publishes a post from the account whose API keys this server runs on. This is public " +
			"and billed per post, with a higher price when the text contains a link — do not call it to " +
			"test whether it works. Text is capped at 280 characters. Set reply_to_id to reply, quote_id " +
			"to quote, media_ids (1 to 4, from " + "x_media_upload" + ") to attach media, or poll for a poll; " +
			"media and poll cannot be combined. Publishing the same text twice in a row is rejected by X " +
			"as a duplicate.",
		OutputSchema: outputFor[dto.CreateResponse](),
	}, h.create)

	mcpsdk.AddTool(s, &mcpsdk.Tool{
		Name:        toolDelete,
		Title:       "Delete an X post",
		Annotations: destructive,
		Description: "Deletes one of your own posts. This is irreversible: the post cannot be restored and " +
			"its replies lose their parent. Confirm the id with " + toolLookup + " before calling.",
		OutputSchema: outputFor[dto.ActionResponse](),
	}, h.delete)

	mcpsdk.AddTool(s, &mcpsdk.Tool{
		Name:         toolLike,
		Title:        "Like an X post",
		Annotations:  writing,
		Description:  "Likes a post as the account whose API keys this server runs on. The like is public.",
		OutputSchema: outputFor[dto.ActionResponse](),
	}, h.like)

	mcpsdk.AddTool(s, &mcpsdk.Tool{
		Name:         toolUnlike,
		Title:        "Remove a like from an X post",
		Annotations:  destructive,
		Description:  "Removes your like from a post.",
		OutputSchema: outputFor[dto.ActionResponse](),
	}, h.unlike)

	mcpsdk.AddTool(s, &mcpsdk.Tool{
		Name:         toolRepost,
		Title:        "Repost an X post",
		Annotations:  writing,
		Description:  "Reposts a post to your followers. The repost is public.",
		OutputSchema: outputFor[dto.ActionResponse](),
	}, h.repost)

	mcpsdk.AddTool(s, &mcpsdk.Tool{
		Name:         toolUnrepost,
		Title:        "Undo a repost",
		Annotations:  destructive,
		Description:  "Removes your repost of a post.",
		OutputSchema: outputFor[dto.ActionResponse](),
	}, h.unrepost)

}
