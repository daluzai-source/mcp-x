package media

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Role1776/mcp-x/app/internal/domain/x"
)

func newUpload(t *testing.T, payload []byte) Upload {
	t.Helper()

	upload, err := NewUpload(UploadProps{Bytes: payload, MIME: "image/jpeg", Category: CategoryImage})
	require.NoError(t, err)

	return upload
}

func TestNewUpload(t *testing.T) {
	imageLimit := newCategory(t, CategoryImage).MaxBytes()

	cases := []struct {
		name      string
		props     UploadProps
		wantErrIs error
	}{
		{name: "a small image", props: UploadProps{Bytes: []byte("data"), MIME: "image/jpeg", Category: CategoryImage}},
		{name: "a gif", props: UploadProps{Bytes: []byte("data"), MIME: "image/gif", Category: CategoryGIF}},
		{name: "a video", props: UploadProps{Bytes: []byte("data"), MIME: "video/mp4", Category: CategoryVideo}},
		{
			name:  "a payload exactly at the category limit",
			props: UploadProps{Bytes: make([]byte, imageLimit), MIME: "image/jpeg", Category: CategoryImage},
		},
		{
			name:      "a payload one byte above the category limit",
			props:     UploadProps{Bytes: make([]byte, imageLimit+1), MIME: "image/jpeg", Category: CategoryImage},
			wantErrIs: x.ErrMediaTooLarge,
		},
		{
			name:      "a nil payload",
			props:     UploadProps{Bytes: nil, MIME: "image/jpeg", Category: CategoryImage},
			wantErrIs: x.ErrMediaEmpty,
		},
		{
			name:      "an empty payload",
			props:     UploadProps{Bytes: []byte{}, MIME: "image/jpeg", Category: CategoryImage},
			wantErrIs: x.ErrMediaEmpty,
		},
		{
			name:      "an unknown category",
			props:     UploadProps{Bytes: []byte("data"), MIME: "image/jpeg", Category: "tweet_audio"},
			wantErrIs: x.ErrUnsupportedMedia,
		},
		{
			name:      "a missing category",
			props:     UploadProps{Bytes: []byte("data"), MIME: "image/jpeg", Category: ""},
			wantErrIs: x.ErrUnsupportedMedia,
		},
		{
			name:      "a mime that belongs to another category",
			props:     UploadProps{Bytes: []byte("data"), MIME: "image/gif", Category: CategoryImage},
			wantErrIs: x.ErrUnsupportedMedia,
		},
		{
			name:      "a missing mime",
			props:     UploadProps{Bytes: []byte("data"), MIME: "", Category: CategoryImage},
			wantErrIs: x.ErrUnsupportedMedia,
		},
		{
			name:      "a mime nobody supports",
			props:     UploadProps{Bytes: []byte("data"), MIME: "application/pdf", Category: CategoryImage},
			wantErrIs: x.ErrUnsupportedMedia,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			upload, err := NewUpload(c.props)

			if c.wantErrIs != nil {
				require.ErrorIs(t, err, c.wantErrIs)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, len(c.props.Bytes), upload.TotalBytes())
			assert.Equal(t, c.props.Category, upload.Category().String())
		})
	}
}

func TestUploadGetters(t *testing.T) {
	upload, err := NewUpload(UploadProps{Bytes: []byte("data"), MIME: "image/jpeg", Category: CategoryImage})
	require.NoError(t, err)

	assert.Equal(t, "image/jpeg", upload.MIME())
	assert.Equal(t, CategoryImage, upload.Category().String())
	assert.Equal(t, 4, upload.TotalBytes())
}

func TestUploadTakesOwnershipOfProps(t *testing.T) {
	payload := []byte("data")

	upload, err := NewUpload(UploadProps{Bytes: payload, MIME: "image/jpeg", Category: CategoryImage})
	require.NoError(t, err)

	payload[0] = 'X'

	assert.Equal(t, []byte("Xata"), bytes.Join(upload.Chunks(64), nil))
}

func TestUploadChunks(t *testing.T) {
	cases := []struct {
		name      string
		payload   []byte
		size      int
		wantCount int
	}{
		{name: "a payload smaller than one chunk", payload: []byte("abc"), size: 10, wantCount: 1},
		{name: "a payload exactly one chunk", payload: []byte("abcde"), size: 5, wantCount: 1},
		{name: "a payload that divides evenly", payload: []byte("abcdef"), size: 3, wantCount: 2},
		{name: "a payload with a short last chunk", payload: []byte("abcdefg"), size: 3, wantCount: 3},
		{name: "one byte per chunk", payload: []byte("abcd"), size: 1, wantCount: 4},
		{name: "a non-positive size does not split and does not hang", payload: []byte("abcd"), size: 0, wantCount: 1},
		{name: "a negative size does not split and does not hang", payload: []byte("abcd"), size: -1, wantCount: 1},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			chunks := newUpload(t, c.payload).Chunks(c.size)

			require.Len(t, chunks, c.wantCount)
			assert.Equal(t, c.payload, bytes.Join(chunks, nil))
		})
	}
}

func TestUploadChunksAliasTheUpload(t *testing.T) {
	upload := newUpload(t, []byte("abcdef"))

	first := upload.Chunks(3)
	first[0][0] = 'X'

	assert.Equal(t, []byte("Xbcdef"), bytes.Join(upload.Chunks(3), nil))
}
