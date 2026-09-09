package users

import (
	"github.com/michimani/gotwi"
	"github.com/michimani/gotwi/fields"

	domainusers "github.com/Role1776/mcp-x/app/internal/domain/x/users"
)

type Adapter struct {
	client gotwi.IClient
	selfID string
}

func New(client gotwi.IClient, selfID domainusers.UserID) *Adapter {
	return &Adapter{client: client, selfID: selfID.String()}
}

var userFields = fields.UserFieldList{
	fields.UserFieldID,
	fields.UserFieldName,
	fields.UserFieldUsername,
	fields.UserFieldDescription,
	fields.UserFieldVerified,
	fields.UserFieldProtected,
	fields.UserFieldPublicMetrics,
	fields.UserFieldProfileImageUrl,
	fields.UserFieldCreatedAt,
}
