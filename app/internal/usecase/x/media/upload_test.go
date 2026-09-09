package media

import (
	"bytes"
	"context"
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Role1776/mcp-x/app/internal/domain/x"
	domainmedia "github.com/Role1776/mcp-x/app/internal/domain/x/media"
	dto "github.com/Role1776/mcp-x/app/internal/dto/x/media"
)

type appendCall struct {
	index int
	size  int
}

type fakeAdapter struct {
	appends  []appendCall
	finalize domainmedia.Media
	statuses []domainmedia.Media
	initErr  error
}

func (f *fakeAdapter) Initialize(_ context.Context, _ domainmedia.Upload) (domainmedia.MediaID, error) {
	if f.initErr != nil {
		return domainmedia.MediaID{}, f.initErr
	}

	return domainmedia.NewMediaID("1800000000000000001")
}

func (f *fakeAdapter) Append(_ context.Context, _ domainmedia.MediaID, index int, chunk []byte) error {
	f.appends = append(f.appends, appendCall{index: index, size: len(chunk)})

	return nil
}

func (f *fakeAdapter) Finalize(_ context.Context, _ domainmedia.MediaID) (domainmedia.Media, error) {
	return f.finalize, nil
}

func (f *fakeAdapter) Status(_ context.Context, _ domainmedia.MediaID) (domainmedia.Media, error) {
	next := f.statuses[0]
	if len(f.statuses) > 1 {
		f.statuses = f.statuses[1:]
	}

	return next, nil
}

func newMedia(t *testing.T, state string) domainmedia.Media {
	t.Helper()

	result, err := domainmedia.NewMedia(domainmedia.MediaProps{ID: "1800000000000000001", ProcessingState: state})
	require.NoError(t, err)

	return result
}

func request(payload []byte, category string) dto.UploadRequest {
	return dto.UploadRequest{
		Base64:   base64.StdEncoding.EncodeToString(payload),
		MIME:     "image/png",
		Category: category,
	}
}

func TestUploadChunking(t *testing.T) {
	adapter := &fakeAdapter{finalize: newMedia(t, "")}

	usecase := newUseCase(t, t.TempDir())
	usecase.media = adapter
	usecase.cfg.ChunkBytes = 4

	res, err := usecase.Upload(context.Background(), request(bytes.Repeat([]byte{1}, 10), "tweet_image"))
	require.NoError(t, err)

	assert.Equal(t, []appendCall{{index: 0, size: 4}, {index: 1, size: 4}, {index: 2, size: 2}}, adapter.appends)
	assert.Equal(t, 3, res.Chunks)
	assert.Equal(t, 10, res.SizeBytes)
	assert.Equal(t, "ready", res.Status)
	assert.Equal(t, "1800000000000000001", res.MediaID)
}

func TestUploadPollsUntilProcessed(t *testing.T) {
	adapter := &fakeAdapter{
		finalize: newMedia(t, domainmedia.StateInProgress),
		statuses: []domainmedia.Media{newMedia(t, domainmedia.StateInProgress), newMedia(t, domainmedia.StateSucceeded)},
	}

	usecase := newUseCase(t, t.TempDir())
	usecase.media = adapter

	res, err := usecase.Upload(context.Background(), request([]byte("payload"), "tweet_image"))
	require.NoError(t, err)

	assert.Equal(t, "ready", res.Status)
	assert.Equal(t, domainmedia.StateSucceeded, res.ProcessingState)
}

func TestUploadReportsProcessingFailure(t *testing.T) {
	adapter := &fakeAdapter{
		finalize: newMedia(t, domainmedia.StatePending),
		statuses: []domainmedia.Media{newMedia(t, domainmedia.StateFailed)},
	}

	usecase := newUseCase(t, t.TempDir())
	usecase.media = adapter

	_, err := usecase.Upload(context.Background(), request([]byte("payload"), "tweet_image"))
	require.ErrorIs(t, err, x.ErrMediaProcessingFailed)
}

func TestUploadRejectsBeforeCallingTheAPI(t *testing.T) {
	cases := []struct {
		name    string
		req     dto.UploadRequest
		wantErr error
	}{
		{
			name:    "too large for its category",
			req:     request(bytes.Repeat([]byte{1}, 5<<20+1), "tweet_image"),
			wantErr: x.ErrMediaTooLarge,
		},
		{
			name: "mime does not match the category",
			req: dto.UploadRequest{
				Base64:   base64.StdEncoding.EncodeToString([]byte("payload")),
				MIME:     "video/mp4",
				Category: "tweet_image",
			},
			wantErr: x.ErrUnsupportedMedia,
		},
		{
			name:    "empty payload",
			req:     request(nil, "tweet_image"),
			wantErr: x.ErrInvalidRequest,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			adapter := &fakeAdapter{finalize: newMedia(t, "")}

			usecase := newUseCase(t, t.TempDir())
			usecase.media = adapter

			_, err := usecase.Upload(context.Background(), c.req)

			require.ErrorIs(t, err, c.wantErr)
			assert.Empty(t, adapter.appends)
		})
	}
}
