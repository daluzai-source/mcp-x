package media

import (
	"encoding/base64"
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Role1776/mcp-x/app/internal/domain/x"
	dto "github.com/Role1776/mcp-x/app/internal/dto/x/media"
)

func newUseCase(t *testing.T, root string) *UseCase {
	t.Helper()

	return New(nil, slog.New(slog.DiscardHandler), &Config{
		Root:             root,
		ChunkBytes:       4 << 20,
		PollIntervalMs:   1,
		DefaultTimeoutMs: 1000,
		MinTimeoutMs:     1,
		MaxTimeoutMs:     10000,
	})
}

func TestReadSource(t *testing.T) {
	root := t.TempDir()
	inside := filepath.Join(root, "photo.png")
	require.NoError(t, os.WriteFile(inside, []byte("\x89PNG\r\n\x1a\npayload"), 0o600))

	outside := filepath.Join(t.TempDir(), "secret.png")
	require.NoError(t, os.WriteFile(outside, []byte("secret"), 0o600))

	usecase := newUseCase(t, root)

	t.Run("reads a file inside the root", func(t *testing.T) {
		src, err := usecase.readSource(dto.UploadRequest{Path: inside, Category: "tweet_image"})

		require.NoError(t, err)
		assert.Equal(t, "image/png", src.mime)
		assert.NotEmpty(t, src.bytes)
	})

	t.Run("refuses a path outside the root", func(t *testing.T) {
		_, err := usecase.readSource(dto.UploadRequest{Path: outside, Category: "tweet_image"})

		require.ErrorIs(t, err, x.ErrMediaPathOutsideRoot)
	})

	t.Run("refuses traversal out of the root", func(t *testing.T) {
		_, err := usecase.readSource(dto.UploadRequest{
			Path:     filepath.Join(root, "..", filepath.Base(filepath.Dir(outside)), "secret.png"),
			Category: "tweet_image",
		})

		require.ErrorIs(t, err, x.ErrMediaPathOutsideRoot)
	})

	t.Run("refuses a symlink that escapes the root", func(t *testing.T) {
		link := filepath.Join(root, "escape.png")
		require.NoError(t, os.Symlink(outside, link))

		_, err := usecase.readSource(dto.UploadRequest{Path: link, Category: "tweet_image"})

		require.ErrorIs(t, err, x.ErrMediaPathOutsideRoot)
	})

	t.Run("reports a missing file", func(t *testing.T) {
		_, err := usecase.readSource(dto.UploadRequest{
			Path:     filepath.Join(root, "nope.png"),
			Category: "tweet_image",
		})

		require.ErrorIs(t, err, x.ErrMediaNotFound)
	})

	t.Run("decodes base64", func(t *testing.T) {
		src, err := usecase.readSource(dto.UploadRequest{
			Base64:   base64.StdEncoding.EncodeToString([]byte("payload")),
			MIME:     "image/jpeg",
			Category: "tweet_image",
		})

		require.NoError(t, err)
		assert.Equal(t, []byte("payload"), src.bytes)
		assert.Equal(t, "image/jpeg", src.mime)
	})

	t.Run("rejects base64 without a mime", func(t *testing.T) {
		_, err := usecase.readSource(dto.UploadRequest{
			Base64:   base64.StdEncoding.EncodeToString([]byte("payload")),
			Category: "tweet_image",
		})

		require.ErrorIs(t, err, x.ErrInvalidRequest)
	})

	t.Run("rejects malformed base64", func(t *testing.T) {
		_, err := usecase.readSource(dto.UploadRequest{
			Base64:   "!!!not base64!!!",
			MIME:     "image/jpeg",
			Category: "tweet_image",
		})

		require.ErrorIs(t, err, x.ErrInvalidRequest)
	})

	t.Run("rejects path and base64 together", func(t *testing.T) {
		_, err := usecase.readSource(dto.UploadRequest{
			Path:     inside,
			Base64:   "aGk=",
			MIME:     "image/jpeg",
			Category: "tweet_image",
		})

		require.ErrorIs(t, err, x.ErrInvalidRequest)
	})

	t.Run("rejects neither path nor base64", func(t *testing.T) {
		_, err := usecase.readSource(dto.UploadRequest{Category: "tweet_image"})

		require.ErrorIs(t, err, x.ErrInvalidRequest)
	})
}
