package media

import (
	"context"

	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"

	dto "github.com/Role1776/mcp-x/app/internal/dto/x/media"
)

const toolUpload = "x_media_upload"

//go:generate mockgen -source=handler.go -destination=usecase_mock_test.go -package=media
type usecase interface {
	Upload(ctx context.Context, req dto.UploadRequest) (dto.UploadResponse, error)
}

type Handler struct {
	usecase usecase
}

func New(usecase usecase) *Handler {
	return &Handler{usecase: usecase}
}

func (h *Handler) RegisterTools(s *mcpsdk.Server) {
	mcpsdk.AddTool(s, &mcpsdk.Tool{
		Name:        toolUpload,
		Title:       "Upload media to X",
		Annotations: &mcpsdk.ToolAnnotations{IdempotentHint: false},
		Description: "Uploads an image, gif or video and returns a media_id to pass to x_post_create. " +
			"Give either 'path', pointing at a local file inside the directory this server was configured " +
			"to read, or 'base64' together with 'mime'. Size limits follow the category: 5 MB for " +
			"tweet_image, 15 MB for tweet_gif, 512 MB for tweet_video. Video is transcoded by X after the " +
			"upload; when the response comes back with status 'pending' the id is valid but not yet usable, " +
			"so retry the post shortly. A media_id expires after a while and is meant for a single post.",
		OutputSchema: outputFor[dto.UploadResponse](),
	}, h.upload)
}

func (h *Handler) upload(ctx context.Context, _ *mcpsdk.CallToolRequest, req dto.UploadRequest) (*mcpsdk.CallToolResult, dto.UploadResponse, error) {
	res, err := h.usecase.Upload(ctx, req)
	if err != nil {
		return result(err), dto.UploadResponse{}, nil
	}

	return nil, res, nil
}
