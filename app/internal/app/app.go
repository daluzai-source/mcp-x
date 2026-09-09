package app

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	xclient "github.com/Role1776/mcp-x/app/internal/adapter/x/client"
	listsAdapter "github.com/Role1776/mcp-x/app/internal/adapter/x/lists"
	mediaAdapter "github.com/Role1776/mcp-x/app/internal/adapter/x/media"
	postsAdapter "github.com/Role1776/mcp-x/app/internal/adapter/x/posts"
	usersAdapter "github.com/Role1776/mcp-x/app/internal/adapter/x/users"
	"github.com/Role1776/mcp-x/app/internal/config"
	"github.com/Role1776/mcp-x/app/internal/domain/x"
	"github.com/Role1776/mcp-x/app/internal/pkg/mcpserver"
	"github.com/Role1776/mcp-x/app/internal/pkg/server"
	mcpRouter "github.com/Role1776/mcp-x/app/internal/transport/mcp/router"
	mcpLists "github.com/Role1776/mcp-x/app/internal/transport/mcp/x/lists"
	mcpMedia "github.com/Role1776/mcp-x/app/internal/transport/mcp/x/media"
	mcpPosts "github.com/Role1776/mcp-x/app/internal/transport/mcp/x/posts"
	mcpUsers "github.com/Role1776/mcp-x/app/internal/transport/mcp/x/users"
	listsUseCase "github.com/Role1776/mcp-x/app/internal/usecase/x/lists"
	mediaUseCase "github.com/Role1776/mcp-x/app/internal/usecase/x/media"
	postsUseCase "github.com/Role1776/mcp-x/app/internal/usecase/x/posts"
	usersUseCase "github.com/Role1776/mcp-x/app/internal/usecase/x/users"
)

const (
	shutdownTimeout  = 10 * time.Second
	bootstrapTimeout = 15 * time.Second
)

type stopper interface {
	Stop(ctx context.Context) error
}

type App struct {
	Config *config.Config
	Logger *slog.Logger
	srv    stopper
}

func New(cfg *config.Config, logger *slog.Logger) *App {
	return &App{
		Config: cfg,
		Logger: logger,
	}
}

func (a *App) Run() error {
	const op = "app.Run"
	client, err := xclient.InitClient(a.Config.XClient)
	if err != nil {
		return fmt.Errorf("%s: failed to create x client: %w", op, err)
	}

	bootstrapCtx, cancelBootstrap := context.WithTimeout(context.Background(), bootstrapTimeout)
	defer cancelBootstrap()

	selfID, err := xclient.Bootstrap(bootstrapCtx, client)
	if err != nil {
		return fmt.Errorf("%s: %w", op, bootstrapError(err))
	}

	a.Logger.Info("authenticated with the X API", slog.String("op", op), slog.String("user_id", selfID.String()))

	posts := postsAdapter.New(client, selfID)
	users := usersAdapter.New(client, selfID)
	lists := listsAdapter.New(client, selfID)
	media := mediaAdapter.New(client)

	postsUC := postsUseCase.New(posts, users, a.Logger, a.Config.Posts)
	usersUC := usersUseCase.New(users, a.Logger, a.Config.Users)
	listsUC := listsUseCase.New(lists, users, a.Logger, a.Config.Lists)
	mediaUC := mediaUseCase.New(media, a.Logger, a.Config.Media)

	mcpSrv := mcpserver.NewServer(a.Config.MCP)

	mcpRouter.New(
		mcpSrv.MCP(),
		mcpPosts.New(postsUC),
		mcpUsers.New(usersUC),
		mcpLists.New(listsUC),
		mcpMedia.New(mediaUC),
	).RegisterAllTools()

	run, err := a.buildMCPServer(mcpSrv)
	if err != nil {
		return err
	}

	notifyCtx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	errCh := make(chan error, 1)

	go func() {
		if err := run(notifyCtx); err != nil && !errors.Is(err, http.ErrServerClosed) && !errors.Is(err, io.EOF) {
			errCh <- fmt.Errorf("%s: failed to run mcp server: %w", op, err)
		}
		close(errCh)

	}()

	a.Logger.Info("MCP server has started",
		slog.String("op", op),
		slog.String("transport", a.Config.MCP.Transport),
	)

	select {
	case <-notifyCtx.Done():
		break
	case err := <-errCh:
		if err != nil {
			return err
		}
	}

	ctxShutdown, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	a.Logger.Info("shutting down", slog.String("op", op))

	if err := a.stop(ctxShutdown); err != nil {
		return err
	}

	return nil
}

func (a *App) buildMCPServer(mcpSrv *mcpserver.Server) (func(context.Context) error, error) {
	const op = "app.buildMCPServer"

	switch a.Config.MCP.Transport {
	case mcpserver.TransportStdio:
		a.srv = mcpSrv

		return mcpSrv.RunStdio, nil

	case mcpserver.TransportHTTP:
		if a.Config.Server.Port == "" {
			return nil, fmt.Errorf("%s: http transport: server port is not configured", op)
		}

		httpSrv := server.NewServer(a.Config.Server, mcpSrv.HTTPHandler())
		a.srv = httpSrv

		return httpSrv.Run, nil

	default:
		return nil, fmt.Errorf("%s: unknown mcp transport: %s", op, a.Config.MCP.Transport)
	}
}

func (a *App) stop(ctx context.Context) error {
	const op = "app.stop"

	if err := a.srv.Stop(ctx); err != nil {
		return fmt.Errorf("%s: failed to shutdown server: %w", op, err)
	}

	return nil
}

func bootstrapError(err error) error {
	switch {
	case errors.Is(err, x.ErrUnauthorized):
		return fmt.Errorf("X rejected the API keys: check X_API_KEY, X_API_KEY_SECRET, "+
			"X_ACCESS_TOKEN and X_ACCESS_TOKEN_SECRET: %w", err)
	case errors.Is(err, x.ErrReadOnlyToken):
		return fmt.Errorf("the access token is read-only: set the app permissions to Read and Write "+
			"in the X Developer Console, then regenerate the Access Token and Secret: %w", err)
	case errors.Is(err, context.DeadlineExceeded):
		return fmt.Errorf("the X API did not answer within %s: %w", bootstrapTimeout, err)
	default:
		return fmt.Errorf("failed to verify the X API credentials: %w", err)
	}
}
