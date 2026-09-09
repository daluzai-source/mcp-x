package lists

import (
	"context"
	"fmt"

	"github.com/michimani/gotwi"
	"github.com/michimani/gotwi/list/managelist"
	"github.com/michimani/gotwi/list/managelist/types"

	"github.com/Role1776/mcp-x/app/internal/adapter/x/apierr"
	"github.com/Role1776/mcp-x/app/internal/domain/x"
	"github.com/Role1776/mcp-x/app/internal/domain/x/lists"
)

func (a *Adapter) Create(ctx context.Context, draft lists.Draft) (string, string, error) {
	const op = "adapter.lists.Create"

	in := &types.CreateInput{Name: draft.Name(), Private: gotwi.Bool(draft.Private())}
	if description := draft.Description(); description != "" {
		in.Description = gotwi.String(description)
	}

	out, err := managelist.Create(ctx, a.client, in)
	if err != nil {
		return "", "", fmt.Errorf("%s: %w", op, apierr.ToDomain(err))
	}

	return out.Data.ID, out.Data.Name, nil
}

func (a *Adapter) Update(ctx context.Context, id lists.ListID, patch lists.Patch) error {
	const op = "adapter.lists.Update"

	out, err := managelist.Update(ctx, a.client, &types.UpdateInput{
		ID:          id.String(),
		Name:        patch.Name(),
		Description: patch.Description(),
		Private:     patch.Private(),
	})
	if err != nil {
		return fmt.Errorf("%s: %w", op, apierr.NotFoundAs(err, x.ErrListNotFound))
	}

	if !out.Data.Updated {
		return fmt.Errorf("%s: %w", op, x.ErrForbidden)
	}

	return nil
}

func (a *Adapter) Delete(ctx context.Context, id lists.ListID) error {
	const op = "adapter.lists.Delete"

	out, err := managelist.Delete(ctx, a.client, &types.DeleteInput{ID: id.String()})
	if err != nil {
		return fmt.Errorf("%s: %w", op, apierr.NotFoundAs(err, x.ErrListNotFound))
	}

	if !out.Data.Deleted {
		return fmt.Errorf("%s: %w", op, x.ErrForbidden)
	}

	return nil
}
