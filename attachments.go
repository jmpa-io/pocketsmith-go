package pocketsmith

import (
	"context"
	"fmt"
	"net/http"

	"go.opentelemetry.io/otel"
)

// CreateAttachmentOptions defines the options for creating an attachment for a
// user.
type CreateAttachmentOptions struct {
	Title    string `json:"title"`
	FileName string `json:"file_name"`
	FileData string `json:"file_data"`
}

// CreateAttachment, using the given user id, creates an attachment for a user.
// https://developers.pocketsmith.com/reference#post_users-id-attachments.
func (c *Client) CreateAttachment(
	ctx context.Context,
	userId int,
	options *CreateAttachmentOptions,
) (attachment *Attachment, err error) {

	// setup tracing.
	newCtx, span := otel.Tracer(c.tracerName).Start(ctx, "CreateAttachment")
	defer span.End()

	// create attachment.
	_, err = c.sender(newCtx, senderRequest{
		method: http.MethodPost,
		path:   fmt.Sprintf("/users/%v/attachments", userId),
		body:   options,
	}, &attachment)
	return attachment, err
}

// CreateAttachmentForAuthedUser, using the token attached to the client,
// creates an attachment for the authed user.
func (c *Client) CreateAttachmentForAuthedUser(
	ctx context.Context,
	options *CreateAttachmentOptions,
) (*Attachment, error) {

	// setup tracing.
	newCtx, span := otel.Tracer(c.tracerName).Start(ctx, "CreateAttachmentForAuthedUser")
	defer span.End()

	// create attachment for user.
	return c.CreateAttachment(newCtx, c.authedUser.ID, options)
}

// DeleteAttachment, using the given attachment id, deletes an attachment.
// https://developers.pocketsmith.com/reference#delete_attachments-id.
func (c *Client) DeleteAttachment(ctx context.Context, attachmentId int) (err error) {

	// setup tracing.
	newCtx, span := otel.Tracer(c.tracerName).Start(ctx, "DeleteAttachment")
	defer span.End()

	// delete attachment.
	_, err = c.sender(newCtx, senderRequest{
		method: http.MethodDelete,
		path:   fmt.Sprintf("/attachments/%v", attachmentId),
	}, nil)
	return err
}

// ListAttachmentsOptions defines the options for listing attachments for a user.
type ListAttachmentsOptions struct {
	Unassigned int `json:"unassigned"`
}

// ListAttachments, using the given user id, lists the attachments for a user.
// https://developers.pocketsmith.com/reference#get_users-id-attachments.
func (c *Client) ListAttachments(
	ctx context.Context,
	userId int,
	options *ListAttachmentsOptions,
) (attachments []Attachment, err error) {

	// setup tracing.
	newCtx, span := otel.Tracer(c.tracerName).Start(ctx, "ListAttachments")
	defer span.End()

	// list attachments.
	_, err = c.sender(newCtx, senderRequest{
		method: http.MethodGet,
		path:   fmt.Sprintf("/users/%v/attachments", userId),
		body:   options,
	}, &attachments)
	return attachments, err
}

// ListAttachmentsForAuthedUser, using the token attached to the client, lists
// the attachment for the authed user.
func (c *Client) ListAttachmentsForAuthedUser(
	ctx context.Context,
	options *ListAttachmentsOptions,
) ([]Attachment, error) {

	// setup tracing.
	newCtx, span := otel.Tracer(c.tracerName).Start(ctx, "ListAttachmentsForAuthedUser")
	defer span.End()

	// list attachments.
	return c.ListAttachments(newCtx, c.authedUser.ID, options)
}

// AssignAttachmentToTransactionOptions defines the options for assigning an
// attachment to a transaction.
type AssignAttachmentToTransactionOptions struct {
	TransactionID int32 `json:"-"`
	AttachmentID  int   `json:"attachment_id"`
}

// AssignAttachmentToTransaction assigns an attachment to a transaction.
// https://developers.pocketsmith.com/reference#post_transactions-id-attachments.
func (c *Client) AssignAttachmentToTransaction(
	ctx context.Context,
	options *AssignAttachmentToTransactionOptions,
) (attachment *Attachment, err error) {

	// setup tracing.
	newCtx, span := otel.Tracer(c.tracerName).Start(ctx, "AssignAttachmentToTransaction")
	defer span.End()

	// assign attachment to transaction.
	_, err = c.sender(newCtx, senderRequest{
		method: http.MethodPost,
		path:   fmt.Sprintf("/transactions/%v/attachments", options.TransactionID),
		body:   options,
	}, &attachment)
	return attachment, err
}
