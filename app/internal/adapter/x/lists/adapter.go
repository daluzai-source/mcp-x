package lists

import (
	"github.com/michimani/gotwi"
	"github.com/michimani/gotwi/fields"

	"github.com/Role1776/mcp-x/app/internal/domain/x/users"
)

type Adapter struct {
	client gotwi.IClient
	selfID string
}

func New(client gotwi.IClient, selfID users.UserID) *Adapter {
	return &Adapter{client: client, selfID: selfID.String()}
}

var (
	listFields = fields.ListFieldList{
		fields.ListFieldCreatedAt,
		fields.ListFieldFollowerCount,
		fields.ListFieldMemberCount,
		fields.ListFieldPrivate,
		fields.ListFieldDescription,
		fields.ListFieldOwnerID,
	}

	userFields = fields.UserFieldList{
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

	tweetFields = fields.TweetFieldList{
		fields.TweetFieldID,
		fields.TweetFieldText,
		fields.TweetFieldAuthorID,
		fields.TweetFieldCreatedAt,
		fields.TweetFieldLang,
		fields.TweetFieldConversationID,
		fields.TweetFieldPublicMetrics,
		fields.TweetFieldEntities,
		fields.TweetFieldAttachments,
		fields.TweetFieldNoteTweet,
	}

	expansions = fields.ExpansionList{fields.ExpansionAuthorID}
)
