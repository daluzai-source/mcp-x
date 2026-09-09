package media

type UploadRequest struct {
	Path      string `json:"path,omitempty" jsonschema:"path to a local file; must sit inside the directory the server was configured with"`
	Base64    string `json:"base64,omitempty" jsonschema:"file contents encoded as base64; use instead of path"`
	MIME      string `json:"mime,omitempty" jsonschema:"media type such as image/jpeg or video/mp4; required together with base64"`
	Category  string `json:"category" jsonschema:"tweet_image (max 5 MB), tweet_gif (max 15 MB) or tweet_video (max 512 MB)" validate:"required,oneof=tweet_image tweet_gif tweet_video"`
	TimeoutMs int64  `json:"timeout_ms,omitempty" jsonschema:"timeout for the whole upload including video processing"`
}

type UploadResponse struct {
	MediaID         string `json:"media_id"`
	ExpiresAt       string `json:"expires_at"`
	ProcessingState string `json:"processing_state,omitempty"`
	Status          string `json:"status"`
	SizeBytes       int    `json:"size_bytes"`
	Chunks          int    `json:"chunks"`
}
