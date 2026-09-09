package media

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/michimani/gotwi/media/upload"
	"github.com/michimani/gotwi/media/upload/types"
	"github.com/michimani/gotwi/resources"

	"github.com/Role1776/mcp-x/app/internal/adapter/x/apierr"
	domainmedia "github.com/Role1776/mcp-x/app/internal/domain/x/media"
)

func (a *Adapter) Initialize(ctx context.Context, u domainmedia.Upload) (domainmedia.MediaID, error) {
	const op = "adapter.media.Initialize"

	out, err := upload.Initialize(ctx, a.client, &types.InitializeInput{
		MediaCategory: types.MediaCategory(u.Category().String()),
		MediaType:     types.MediaType(u.MIME()),
		TotalBytes:    u.TotalBytes(),
	})
	if err != nil {
		return domainmedia.MediaID{}, fmt.Errorf("%s: %w", op, apierr.ToDomain(err))
	}

	id, err := domainmedia.NewMediaID(out.Data.MediaID)
	if err != nil {
		return domainmedia.MediaID{}, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (a *Adapter) Append(ctx context.Context, id domainmedia.MediaID, index int, chunk []byte) error {
	const op = "adapter.media.Append"

	if _, err := upload.Append(ctx, a.client, &types.AppendInput{
		MediaID:      id.String(),
		Media:        bytes.NewReader(chunk),
		SegmentIndex: index,
	}); err != nil {
		return fmt.Errorf("%s: %w", op, apierr.ToDomain(err))
	}

	return nil
}

func (a *Adapter) Finalize(ctx context.Context, id domainmedia.MediaID) (domainmedia.Media, error) {
	const op = "adapter.media.Finalize"

	out, err := upload.Finalize(ctx, a.client, &types.FinalizeInput{MediaID: id.String()})
	if err != nil {
		return domainmedia.Media{}, fmt.Errorf("%s: %w", op, apierr.ToDomain(err))
	}

	result, err := toMedia(out.Data)
	if err != nil {
		return domainmedia.Media{}, fmt.Errorf("%s: %w", op, err)
	}

	return result, nil
}

func toMedia(data resources.UploadedMedia) (domainmedia.Media, error) {
	return domainmedia.NewMedia(domainmedia.MediaProps{
		ID:              data.MediaID,
		ExpiresAt:       time.Now().Add(time.Duration(data.ExpiresAfterSecs) * time.Second),
		ProcessingState: string(data.ProcessingInfo.State),
		ProgressPercent: data.ProcessingInfo.ProgressPercent,
		CheckAfterSecs:  data.ProcessingInfo.CheckAfterSecs,
	})
}
