package media

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Role1776/mcp-x/app/internal/domain/x"
)

func newCategory(t *testing.T, raw string) Category {
	t.Helper()

	category, err := NewCategory(raw)
	require.NoError(t, err)

	return category
}

func TestNewCategory(t *testing.T) {
	cases := []struct {
		name        string
		raw         string
		wantIsVideo bool
		wantErrIs   error
	}{
		{name: "an image", raw: CategoryImage},
		{name: "a gif", raw: CategoryGIF},
		{name: "a video", raw: CategoryVideo, wantIsVideo: true},
		{name: "empty string is rejected", raw: "", wantErrIs: x.ErrUnsupportedMedia},
		{name: "a bare mime type is rejected", raw: "image/jpeg", wantErrIs: x.ErrUnsupportedMedia},
		{name: "a short form is rejected", raw: "image", wantErrIs: x.ErrUnsupportedMedia},
		{name: "the enum is case sensitive", raw: "TWEET_IMAGE", wantErrIs: x.ErrUnsupportedMedia},
		{name: "an unknown category is rejected", raw: "tweet_audio", wantErrIs: x.ErrUnsupportedMedia},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			category, err := NewCategory(c.raw)

			if c.wantErrIs != nil {
				require.ErrorIs(t, err, c.wantErrIs)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, c.raw, category.String())
			assert.Equal(t, c.wantIsVideo, category.IsVideo())
		})
	}
}

func TestCategoryMaxBytes(t *testing.T) {
	image := newCategory(t, CategoryImage).MaxBytes()
	gif := newCategory(t, CategoryGIF).MaxBytes()
	video := newCategory(t, CategoryVideo).MaxBytes()

	assert.Positive(t, image)
	assert.Less(t, image, gif)
	assert.Less(t, gif, video)
	assert.Zero(t, Category{}.MaxBytes())
}

func TestCategoryAllowsMIME(t *testing.T) {
	cases := []struct {
		name     string
		category string
		mime     string
		want     bool
	}{
		{name: "jpeg belongs to images", category: CategoryImage, mime: "image/jpeg", want: true},
		{name: "png belongs to images", category: CategoryImage, mime: "image/png", want: true},
		{name: "webp belongs to images", category: CategoryImage, mime: "image/webp", want: true},
		{name: "a gif is not an image upload", category: CategoryImage, mime: "image/gif", want: false},
		{name: "a video is not an image upload", category: CategoryImage, mime: "video/mp4", want: false},
		{name: "gif belongs to the gif category", category: CategoryGIF, mime: "image/gif", want: true},
		{name: "jpeg does not belong to the gif category", category: CategoryGIF, mime: "image/jpeg", want: false},
		{name: "mp4 belongs to videos", category: CategoryVideo, mime: "video/mp4", want: true},
		{name: "an image does not belong to videos", category: CategoryVideo, mime: "image/jpeg", want: false},
		{name: "an empty mime is never allowed", category: CategoryImage, mime: "", want: false},
		{name: "an unknown mime is never allowed", category: CategoryImage, mime: "application/pdf", want: false},
		{name: "matching ignores case", category: CategoryImage, mime: "IMAGE/JPEG", want: true},
		{name: "matching ignores mime parameters", category: CategoryImage, mime: "image/jpeg; charset=binary", want: true},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			assert.Equal(t, c.want, newCategory(t, c.category).AllowsMIME(c.mime))
		})
	}
}

func TestZeroCategoryAllowsNothing(t *testing.T) {
	assert.False(t, Category{}.AllowsMIME("image/jpeg"))
	assert.False(t, Category{}.IsVideo())
}
