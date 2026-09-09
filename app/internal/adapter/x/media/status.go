package media

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/michimani/gotwi/resources"

	"github.com/Role1776/mcp-x/app/internal/adapter/x/apierr"
	domainmedia "github.com/Role1776/mcp-x/app/internal/domain/x/media"
)

const statusEndpoint = "https://api.x.com/2/media/upload"

type statusParams struct {
	accessToken string
	mediaID     string
}

func (p *statusParams) SetAccessToken(token string) {
	p.accessToken = token
}

func (p *statusParams) AccessToken() string {
	return p.accessToken
}

func (p *statusParams) ResolveEndpoint(endpointBase string) string {
	return endpointBase + "?command=STATUS&media_id=" + p.mediaID
}

func (p *statusParams) Body() (io.Reader, error) {
	return nil, nil
}

func (p *statusParams) ParameterMap() map[string]string {
	return map[string]string{"command": "STATUS", "media_id": p.mediaID}
}

type statusResponse struct {
	Data   resources.UploadedMedia  `json:"data"`
	Errors []resources.PartialError `json:"errors"`
}

func (r *statusResponse) HasPartialError() bool {
	return len(r.Errors) > 0
}

func (a *Adapter) Status(ctx context.Context, id domainmedia.MediaID) (domainmedia.Media, error) {
	const op = "adapter.media.Status"

	out := &statusResponse{}
	if err := a.client.CallAPI(ctx, statusEndpoint, http.MethodGet, &statusParams{mediaID: id.String()}, out); err != nil {
		return domainmedia.Media{}, fmt.Errorf("%s: %w", op, apierr.ToDomain(err))
	}

	if out.Data.MediaID == "" {
		out.Data.MediaID = id.String()
	}

	result, err := toMedia(out.Data)
	if err != nil {
		return domainmedia.Media{}, fmt.Errorf("%s: %w", op, err)
	}

	return result, nil
}
