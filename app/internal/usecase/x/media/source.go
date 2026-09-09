package media

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/Role1776/mcp-x/app/internal/domain/x"
	dto "github.com/Role1776/mcp-x/app/internal/dto/x/media"
)

type source struct {
	bytes []byte
	mime  string
}

func (u *UseCase) readSource(req dto.UploadRequest) (source, error) {
	const op = "usecase.media.readSource"

	if (req.Path == "") == (req.Base64 == "") {
		return source{}, fmt.Errorf("%s: %w", op, x.ErrInvalidRequest)
	}

	if req.Base64 != "" {
		if req.MIME == "" {
			return source{}, fmt.Errorf("%s: %w", op, x.ErrInvalidRequest)
		}

		decoded, err := base64.StdEncoding.DecodeString(req.Base64)
		if err != nil {
			return source{}, fmt.Errorf("%s: %w", op, x.ErrInvalidRequest)
		}

		return source{bytes: decoded, mime: req.MIME}, nil
	}

	path, err := u.resolvePath(req.Path)
	if err != nil {
		return source{}, fmt.Errorf("%s: %w", op, err)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return source{}, fmt.Errorf("%s: %w", op, x.ErrMediaNotFound)
	}

	mime := req.MIME
	if mime == "" {
		mime = detectMIME(content)
	}

	return source{bytes: content, mime: mime}, nil
}

// resolvePath expands symlinks before comparing against the root: a link that
// lives inside the root but points outside it would otherwise let an MCP client
// read any file on the host.
func (u *UseCase) resolvePath(raw string) (string, error) {
	root, err := filepath.EvalSymlinks(u.cfg.Root)
	if err != nil {
		return "", x.ErrMediaPathOutsideRoot
	}

	path, err := filepath.Abs(raw)
	if err != nil {
		return "", x.ErrMediaPathOutsideRoot
	}

	path, err = filepath.EvalSymlinks(path)
	if err != nil {
		return "", x.ErrMediaNotFound
	}

	rel, err := filepath.Rel(root, path)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", x.ErrMediaPathOutsideRoot
	}

	return path, nil
}

func detectMIME(content []byte) string {
	mime := http.DetectContentType(content)
	if index := strings.IndexByte(mime, ';'); index >= 0 {
		mime = mime[:index]
	}

	return strings.TrimSpace(mime)
}
