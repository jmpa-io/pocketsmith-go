package pocketsmith

import "time"

// Attachment defines a PocketSmith attachment.
// https://developers.pocketsmith.com/reference#get_attachments-id
type Attachment struct {
	ID              int    `json:"id"`
	Title           string `json:"title"`
	FileName        string `json:"file_name"`
	Type            string `json:"type"`
	ContentType     string `json:"content_type"`
	ContentTypeMeta struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Extension   string `json:"extension"`
	} `json:"content_type_meta"`
	OriginalURL string `json:"original_url"`
	Variants    struct {
		ThumbURL string `json:"thumb_url"`
		LargeURL string `json:"large_url"`
	} `json:"variants"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
