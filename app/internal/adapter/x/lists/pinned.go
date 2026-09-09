package lists

import (
	"context"
	"fmt"

	"github.com/michimani/gotwi/list/pinnedlist"
	"github.com/michimani/gotwi/list/pinnedlist/types"

	"github.com/Role1776/mcp-x/app/internal/adapter/x/apierr"
	"github.com/Role1776/mcp-x/app/internal/domain/x"
	domainlists "github.com/Role1776/mcp-x/app/internal/domain/x/lists"
)

func (a *Adapter) Pin(ctx context.Context, id domainlists.ListID) error {
	const op = "adapter.lists.Pin"

	if _, err := pinnedlist.Create(ctx, a.client, &types.CreateInput{ID: a.selfID, ListID: id.String()}); err != nil {
		return fmt.Errorf("%s: %w", op, apierr.NotFoundAs(err, x.ErrListNotFound))
	}

	return nil
}

func (a *Adapter) Unpin(ctx context.Context, id domainlists.ListID) error {
	const op = "adapter.lists.Unpin"

	if _, err := pinnedlist.Delete(ctx, a.client, &types.DeleteInput{ID: a.selfID, ListID: id.String()}); err != nil {
		return fmt.Errorf("%s: %w", op, apierr.NotFoundAs(err, x.ErrListNotFound))
	}

	return nil
}

func (a *Adapter) Pinned(ctx context.Context) (domainlists.Lists, error) {
	const op = "adapter.lists.Pinned"

	out, err := pinnedlist.List(ctx, a.client, &types.ListInput{ID: a.selfID, ListFields: listFields})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, apierr.ToDomain(err))
	}

	return toLists(out.Data), nil
}
