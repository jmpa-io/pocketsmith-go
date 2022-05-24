package pocketsmith

import (
	"fmt"
	"net/http"
	"time"
)

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

// CreateAttachmentOptions defines the options for creating an attachment for a
// user.
type CreateAttachmentOptions struct {
	Title    string `json:"title"`
	FileName string `json:"file_name"`
	FileData string `json:"file_data"`
}

// CreateAttachment, using the given user id, creates an attachment for a user.
// https://developers.pocketsmith.com/reference#post_users-id-attachments
func (c *Client) CreateAttachment(userId int, options CreateAttachmentOptions) (attachment *Attachment, err error) {
	cr := clientRequest{
		method: http.MethodPost,
		path:   fmt.Sprintf("/users/%v/attachments", userId),
		data:   options,
	}
	_, err = c.sender(cr, &attachment)
	return attachment, err
}

// CreateAttachmentForAuthedUser, using the token attached to the client,
// creates an attachment for the authed user.
func (c *Client) CreateAttachmentForAuthedUser(options CreateAttachmentOptions) (*Attachment, error) {
	return c.CreateAttachment(c.user.ID, options)
}

// DeleteAttachment, using the given attachment id, deletes an attachment.
// https://developers.pocketsmith.com/reference#delete_attachments-id
func (c *Client) DeleteAttachment(attachmentId int) error {
	cr := clientRequest{
		method: http.MethodDelete,
		path:   fmt.Sprintf("/attachments/%v", attachmentId),
	}
	_, err := c.sender(cr, nil)
	return err
}

// ListAttachmentsOptions defines the options for listing attachments for a user.
type ListAttachmentsOptions struct {
	Unassigned int `json:"unassigned"`
}

// ListAttachments, using the given user id, lists the attachments for a user.
// https://developers.pocketsmith.com/reference#get_users-id-attachments
func (c *Client) ListAttachments(userId int, options ListAttachmentsOptions) ([]Attachment, error) {
	cr := clientRequest{
		method: http.MethodGet,
		path:   fmt.Sprintf("/users/%v/attachments", userId),
		data:   options,
	}
	var attachments []Attachment
	_, err := c.sender(cr, &attachments)
	return attachments, err
}

// ListAttachmentsForAuthedUser, using the token attached to the client, lists
// the attachment for the authed user.
func (c *Client) ListAttachmentsForAuthedUser(options ListAttachmentsOptions) ([]Attachment, error) {
	return c.ListAttachments(c.user.ID, options)
}

// AssignAttachmentToTransactionOptions defines the options for assigning an
// attachment to a transaction.
type AssignAttachmentToTransactionOptions struct {
	TransactionID int32 `json:"-"`
	AttachmentID  int   `json:"attachment_id"`
}

// AssignAttachmentToTransaction assigns an attachment to a transaction.
// https://developers.pocketsmith.com/reference#post_transactions-id-attachments
func (c *Client) AssignAttachmentToTransaction(options AssignAttachmentToTransactionOptions) (*Attachment, error) {
	cr := clientRequest{
		method: http.MethodPost,
		path:   fmt.Sprintf("/transactions/%v/attachments", options.TransactionID),
		data:   options,
	}
	var attachment *Attachment
	_, err := c.sender(cr, &attachment)
	return attachment, err
}
