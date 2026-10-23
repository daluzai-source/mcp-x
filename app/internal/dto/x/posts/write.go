package posts

type Poll struct {
	Options         []string `json:"options" jsonschema:"2 to 4 poll options" validate:"required,min=2,max=4,dive,required"`
	DurationMinutes int      `json:"duration_minutes" jsonschema:"how long the poll stays open, 5..10080 minutes" validate:"required,gte=5,lte=10080"`
}

type CreateRequest struct {
	Text          string   `json:"text,omitempty" jsonschema:"post text, at most 25000 characters (Premium long post); required unless media_ids is set"`
	ReplyToID     string   `json:"reply_to_id,omitempty" jsonschema:"id of the post this one replies to"`
	QuoteID       string   `json:"quote_id,omitempty" jsonschema:"id of the post to quote"`
	MediaIDs      []string `json:"media_ids,omitempty" jsonschema:"1 to 4 media ids returned by x_media_upload" validate:"omitempty,max=4,dive,required"`
	Poll          *Poll    `json:"poll,omitempty" jsonschema:"attach a poll; cannot be combined with media_ids"`
	ReplySettings string   `json:"reply_settings,omitempty" jsonschema:"who may reply: following, mentionedUsers, subscribers or verified" validate:"omitempty,oneof=following mentionedUsers subscribers verified"`
}

type CreateResponse struct {
	ID   string `json:"id"`
	Text string `json:"text"`
	URL  string `json:"url"`
}

type PostIDRequest struct {
	ID string `json:"id" jsonschema:"post id" validate:"required"`
}

type ActionResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}
