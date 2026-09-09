package lists

import (
	"context"
	"fmt"

	"github.com/michimani/gotwi/list/listfollow"
	"github.com/michimani/gotwi/list/listfollow/types"

	"github.com/Role1776/mcp-x/app/internal/adapter/x/apierr"
	"github.com/Role1776/mcp-x/app/internal/domain/x"
	domainlists "github.com/Role1776/mcp-x/app/internal/domain/x/lists"
)

func (a *Adapter) Follow(ctx context.Context, id domainlists.ListID) error {
	const op = "adapter.lists.Follow"

	if _, err := listfollow.Create(ctx, a.client, &types.CreateInput{ID: a.selfID, ListID: id.String()}); err != nil {
		return fmt.Errorf("%s: %w", op, apierr.NotFoundAs(err, x.ErrListNotFound))
	}

	return nil
}

func (a *Adapter) Unfollow(ctx context.Context, id domainlists.ListID) error {
	const op = "adapter.lists.Unfollow"

	if _, err := listfollow.Delete(ctx, a.client, &types.DeleteInput{ID: a.selfID, ListID: id.String()}); err != nil {
		return fmt.Errorf("%s: %w", op, apierr.NotFoundAs(err, x.ErrListNotFound))
	}

	return nil
}
