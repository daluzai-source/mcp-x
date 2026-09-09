package media

import (
	"strings"

	"github.com/Role1776/mcp-x/app/internal/domain/x"
)

const (
	CategoryImage = "tweet_image"
	CategoryGIF   = "tweet_gif"
	CategoryVideo = "tweet_video"

	maxImageBytes = 5 << 20
	maxGIFBytes   = 15 << 20
	maxVideoBytes = 512 << 20
)

var allowedMIME = map[string]map[string]struct{}{
	CategoryImage: {
		"image/jpeg":  {},
		"image/pjpeg": {},
		"image/png":   {},
		"image/webp":  {},
		"image/bmp":   {},
		"image/tiff":  {},
	},
	CategoryGIF: {
		"image/gif": {},
	},
	CategoryVideo: {
		"video/mp4":       {},
		"video/webm":      {},
		"video/mp2t":      {},
		"video/quicktime": {},
	},
}

type Category struct {
	value string
}

func NewCategory(raw string) (Category, error) {
	if _, ok := allowedMIME[raw]; !ok {
		return Category{}, x.ErrUnsupportedMedia
	}

	return Category{value: raw}, nil
}

func (c Category) String() string {
	return c.value
}

func (c Category) IsVideo() bool {
	return c.value == CategoryVideo
}

func (c Category) MaxBytes() int64 {
	switch c.value {
	case CategoryImage:
		return maxImageBytes
	case CategoryGIF:
		return maxGIFBytes
	case CategoryVideo:
		return maxVideoBytes
	default:
		return 0
	}
}

func (c Category) AllowsMIME(mime string) bool {
	allowed, ok := allowedMIME[c.value]
	if !ok {
		return false
	}

	_, ok = allowed[normalizeMIME(mime)]

	return ok
}

func normalizeMIME(mime string) string {
	value := mime
	if index := strings.IndexByte(value, ';'); index >= 0 {
		value = value[:index]
	}

	return strings.ToLower(strings.TrimSpace(value))
}
