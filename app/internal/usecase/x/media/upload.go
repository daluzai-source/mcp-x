package media

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/Role1776/mcp-x/app/internal/domain/x"
	"github.com/Role1776/mcp-x/app/internal/domain/x/media"
	dto "github.com/Role1776/mcp-x/app/internal/dto/x/media"
	"github.com/Role1776/mcp-x/app/internal/pkg/validator"
)

func (u *UseCase) Upload(ctx context.Context, req dto.UploadRequest) (dto.UploadResponse, error) {
	const op = "usecase.media.Upload"

	if err := validator.Validate(req); err != nil {
		return dto.UploadResponse{}, fmt.Errorf("%s: %w", op, x.ErrInvalidRequest)
	}

	src, err := u.readSource(req)
	if err != nil {
		return dto.UploadResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	upload, err := media.NewUpload(media.UploadProps{
		Bytes:    src.bytes,
		MIME:     src.mime,
		Category: req.Category,
	})
	if err != nil {
		return dto.UploadResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	ctx, cancel := u.withTimeout(ctx, req.TimeoutMs)
	defer cancel()

	id, err := u.media.Initialize(ctx, upload)
	if err != nil {
		return dto.UploadResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	chunks := upload.Chunks(u.cfg.ChunkBytes)
	for index, chunk := range chunks {
		if err := u.media.Append(ctx, id, index, chunk); err != nil {
			return dto.UploadResponse{}, fmt.Errorf("%s: %w", op, err)
		}
	}

	result, err := u.media.Finalize(ctx, id)
	if err != nil {
		return dto.UploadResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	if result.IsPending() {
		result, err = u.waitProcessed(ctx, result)
		if err != nil {
			return dto.UploadResponse{}, fmt.Errorf("%s: %w", op, err)
		}
	}

	return u.response(result, upload.TotalBytes(), len(chunks)), nil
}

func (u *UseCase) waitProcessed(ctx context.Context, current media.Media) (media.Media, error) {
	const op = "usecase.media.waitProcessed"

	interval := time.Duration(u.cfg.PollIntervalMs) * time.Millisecond

	for current.IsPending() {
		wait := max(current.CheckAfter(), interval)

		timer := time.NewTimer(wait)
		select {
		case <-ctx.Done():
			timer.Stop()

			if errors.Is(ctx.Err(), context.DeadlineExceeded) {
				return current, nil
			}

			return current, fmt.Errorf("%s: %w", op, ctx.Err())
		case <-timer.C:
		}

		next, err := u.media.Status(ctx, current.ID())
		if err != nil {
			u.logger.Error("media status failed", slog.String("op", op), slog.Any("err", err))

			return current, err
		}

		current = next
	}

	if current.IsFailed() {
		return current, x.ErrMediaProcessingFailed
	}

	return current, nil
}

func (u *UseCase) response(result media.Media, size, chunks int) dto.UploadResponse {
	view := result.View()

	status := statusReady
	if result.IsPending() {
		status = statusPending
	}

	response := dto.UploadResponse{
		MediaID:         view.ID,
		ProcessingState: view.ProcessingState,
		Status:          status,
		SizeBytes:       size,
		Chunks:          chunks,
	}

	if !view.ExpiresAt.IsZero() {
		response.ExpiresAt = view.ExpiresAt.Format(time.RFC3339)
	}

	return response
}
