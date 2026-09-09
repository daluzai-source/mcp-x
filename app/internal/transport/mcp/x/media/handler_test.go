package media

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/Role1776/mcp-x/app/internal/domain/x"
	dto "github.com/Role1776/mcp-x/app/internal/dto/x/media"
)

func TestUpload(t *testing.T) {
	req := dto.UploadRequest{Path: "/media/cat.jpg", Category: "tweet_image"}

	okResponse := dto.UploadResponse{
		MediaID:         "1146654567674912769",
		ExpiresAt:       "2026-01-01T00:00:00Z",
		ProcessingState: "succeeded",
		Status:          "ready",
		SizeBytes:       1024,
		Chunks:          1,
	}

	cases := []struct {
		name string
		res  dto.UploadResponse
		err  error

		wantIsError bool
		wantText    string
	}{
		{name: "success", res: okResponse},
		{
			name: "success while the video is still processing",
			res: dto.UploadResponse{
				MediaID: "1146654567674912769", ProcessingState: "in_progress", Status: "pending",
			},
		},
		{
			name:        "the file is outside the configured root",
			err:         x.ErrMediaPathOutsideRoot,
			wantIsError: true,
			wantText:    "the path is outside the directory this server is allowed to read",
		},
		{
			name:        "the file does not exist",
			err:         x.ErrMediaNotFound,
			wantIsError: true,
			wantText:    "media file not found",
		},
		{
			name:        "the file is empty",
			err:         x.ErrMediaEmpty,
			wantIsError: true,
			wantText:    "the media file is empty",
		},
		{
			name:        "the file is too large for its category",
			err:         x.ErrMediaTooLarge,
			wantIsError: true,
			wantText:    "media exceeds the size limit for its category: 5 MB for images, 15 MB for gifs, 512 MB for video",
		},
		{
			name:        "unsupported media type",
			err:         x.ErrUnsupportedMedia,
			wantIsError: true,
			wantText:    "unsupported media type for this category",
		},
		{
			name:        "processing failed on the X side",
			err:         x.ErrMediaProcessingFailed,
			wantIsError: true,
			wantText:    "X failed to process the uploaded media; re-encode the file and upload it again",
		},
		{
			name:        "neither path nor base64",
			err:         x.ErrInvalidRequest,
			wantIsError: true,
			wantText:    "invalid request",
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
			usecase.EXPECT().Upload(gomock.Any(), req).Return(c.res, c.err)

			res, out, err := New(usecase).upload(context.Background(), nil, req)
			require.NoError(t, err)

			if c.wantIsError {
				assert.True(t, res.IsError)
				assert.Equal(t, c.wantText, resultText(t, res))
				assert.Equal(t, dto.UploadResponse{}, out)

				return
			}

			assert.Nil(t, res)
			assert.Equal(t, c.res, out)
		})
	}
}
