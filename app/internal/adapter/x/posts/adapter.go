package posts

import (
	"github.com/michimani/gotwi"

	"github.com/Role1776/mcp-x/app/internal/domain/x/users"
)

type Adapter struct {
	client gotwi.IClient
	selfID string
}

func New(client gotwi.IClient, selfID users.UserID) *Adapter {
	return &Adapter{client: client, selfID: selfID.String()}
}
