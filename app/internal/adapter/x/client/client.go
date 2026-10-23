package client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/michimani/gotwi"
	"github.com/michimani/gotwi/user/userlookup"
	"github.com/michimani/gotwi/user/userlookup/types"

	"github.com/Role1776/mcp-x/app/internal/adapter/x/apierr"
	"github.com/Role1776/mcp-x/app/internal/domain/x"
	"github.com/Role1776/mcp-x/app/internal/domain/x/users"
)

type Config struct {
	APIKey              string `env:"X_API_KEY" validate:"required"`
	APIKeySecret        string `env:"X_API_KEY_SECRET" validate:"required"`
	AccessToken         string `env:"X_ACCESS_TOKEN" validate:"required"`
	AccessTokenSecret   string `env:"X_ACCESS_TOKEN_SECRET" validate:"required"`
	MaxIdleConnsPerHost int    `env:"MAX_IDLE_CONNS_PER_HOST" envDefault:"100" validate:"gt=0"`
	Debug               bool   `env:"X_DEBUG" envDefault:"false"`
}

func InitClient(cfg *Config) (*gotwi.Client, error) {
	const op = "client.InitClient"

	httpClient := &http.Client{
		Transport: &http.Transport{
			Proxy:               http.ProxyFromEnvironment,
			MaxIdleConnsPerHost: cfg.MaxIdleConnsPerHost,
		},
	}

	c, err := gotwi.NewClient(&gotwi.NewClientInput{
		HTTPClient:           httpClient,
		AuthenticationMethod: gotwi.AuthenMethodOAuth1UserContext,
		APIKey:               cfg.APIKey,
		APIKeySecret:         cfg.APIKeySecret,
		OAuthToken:           cfg.AccessToken,
		OAuthTokenSecret:     cfg.AccessTokenSecret,
		Debug:                cfg.Debug,
	})

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return c, nil
}

// Bootstrap resolves the key owner once at startup: every write endpoint is
// POST /2/users/:id/..., the id never changes while the process lives, and a
// lookup per write would be another billed request. A failure here also means
// the keys are unusable, so the server must not start.
func Bootstrap(ctx context.Context, c gotwi.IClient) (users.UserID, error) {
	const op = "client.Bootstrap"

	out, err := userlookup.GetMe(ctx, c, &types.GetMeInput{})
	if err != nil {
		return users.UserID{}, fmt.Errorf("%s: %w", op, apierr.ToDomain(err))
	}

	if out.Data.ID == nil {
		return users.UserID{}, fmt.Errorf("%s: %w", op, x.ErrUnauthorized)
	}

	id, err := users.NewUserID(gotwi.StringValue(out.Data.ID))
	if err != nil {
		return users.UserID{}, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}
