package media

import (
	"github.com/Role1776/mcp-x/app/internal/domain/x"
)

type Upload struct {
	bytes    []byte
	mime     string
	category Category
}

type UploadProps struct {
	Bytes    []byte
	MIME     string
	Category string
}

// NewUpload takes ownership of props.Bytes rather than copying it: the payload is
// as large as the category allows, and every caller hands over a buffer it has just
// read and does not reuse. The caller must not touch props.Bytes afterwards.
func NewUpload(props UploadProps) (Upload, error) {
	category, err := NewCategory(props.Category)
	if err != nil {
		return Upload{}, err
	}

	if len(props.Bytes) == 0 {
		return Upload{}, x.ErrMediaEmpty
	}

	if !category.AllowsMIME(props.MIME) {
		return Upload{}, x.ErrUnsupportedMedia
	}

	if int64(len(props.Bytes)) > category.MaxBytes() {
		return Upload{}, x.ErrMediaTooLarge
	}

	return Upload{bytes: props.Bytes, mime: props.MIME, category: category}, nil
}

func (u Upload) MIME() string {
	return u.mime
}

func (u Upload) Category() Category {
	return u.category
}

func (u Upload) TotalBytes() int {
	return len(u.bytes)
}

// Chunks splits the payload into pieces of at most size bytes, or returns the
// whole payload as a single piece when size is not positive. The returned slices
// alias the upload's payload and must not be modified; copy a chunk if the caller
// needs to own it.
func (u Upload) Chunks(size int) [][]byte {
	if size <= 0 {
		return [][]byte{u.bytes}
	}

	chunks := make([][]byte, 0, (len(u.bytes)+size-1)/size)
	for start := 0; start < len(u.bytes); start += size {
		end := start + size
		if end > len(u.bytes) {
			end = len(u.bytes)
		}

		chunks = append(chunks, u.bytes[start:end])
	}

	return chunks
}
