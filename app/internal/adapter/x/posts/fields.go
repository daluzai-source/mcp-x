package posts

import "github.com/michimani/gotwi/fields"

var (
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
		fields.TweetFieldReferencedTweets,
		fields.TweetFieldNoteTweet,
	}

	userFields = fields.UserFieldList{
		fields.UserFieldID,
		fields.UserFieldName,
		fields.UserFieldUsername,
		fields.UserFieldVerified,
		fields.UserFieldProtected,
	}

	mediaFields = fields.MediaFieldList{
		fields.MediaFieldMediaKey,
		fields.MediaFieldType,
		fields.MediaFieldUrl,
		fields.MediaFieldPreviewImageUrl,
		fields.MediaFieldAltText,
	}

	expansions = fields.ExpansionList{
		fields.ExpansionAuthorID,
		fields.ExpansionAttachmentsMediaKeys,
	}

	listUserFields = fields.UserFieldList{
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
)
