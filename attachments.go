package pocketsmith

import (
	"context"
	"fmt"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
)

// Attachments represents a slice of Attachment.
type Attachments []Attachment

// CreateAttachmentOptions defines the options for creating an attachment for a
// user.
type CreateAttachmentOptions struct {
	Title    string `json:"title"`
	FileName string `json:"file_name"`
	FileData string `json:"file_data"`
}

// CreateAttachmentForUserOptions ...
type CreateAttachmentForUserOptions struct {
	UserID int `json:"-" validate:"required"`

	CreateAttachmentOptions
}

// CreateAttachmentForUser, using the given user id, creates an attachment for a user.
// https://developers.pocketsmith.com/reference#post_users-id-attachments.
func (c *Client) CreateAttachmentForUser(
	ctx context.Context,
	options *CreateAttachmentForUserOptions,
) (attachment *Attachment, err error) {

	// setup tracing.
	newCtx, span := otel.Tracer(c.tracerName).Start(ctx, "CreateAttachmentForUser")
	defer span.End()

	// validate options.
	if err := c.validator.StructCtx(newCtx, options); err != nil {
		span.SetStatus(codes.Error, fmt.Sprintf("failed to validate options: %v", err))
		span.RecordError(err)
		return nil, err
	}

	// create attachment.
	_, err = c.sender(newCtx, senderRequest{
		method: http.MethodPost,
		path:   fmt.Sprintf("/users/%v/attachments", options.UserID),
		body:   options,
	}, &attachment)
	if err != nil {
		span.SetStatus(codes.Error, fmt.Sprintf("failed to create attachment: %v", err))
		span.RecordError(err)
		return nil, err
	}
	return attachment, nil
}

// CreateAttachment, using the token attached to the client,
// creates an attachment for the authed user.
func (c *Client) CreateAttachment(
	ctx context.Context,
	options *CreateAttachmentOptions,
) (*Attachment, error) {

	// setup tracing.
	newCtx, span := otel.Tracer(c.tracerName).Start(ctx, "CreateAttachment")
	defer span.End()

	userID, err := c.authedUserID(newCtx)
	if err != nil {
		span.SetStatus(codes.Error, fmt.Sprintf("failed to get authed user: %v", err))
		span.RecordError(err)
		return nil, err
	}

	// create attachment for user.
	attachment, err := c.CreateAttachmentForUser(
		newCtx,
		&CreateAttachmentForUserOptions{UserID: userID, CreateAttachmentOptions: *options},
	)
	if err != nil {
		span.SetStatus(codes.Error, fmt.Sprintf("failed to create attachment: %v", err))
		span.RecordError(err)
		return nil, err
	}
	return attachment, nil
}

// DeleteAttachmentOptions ...
type DeleteAttachmentOptions struct {
	AttachmentID int `validate:"required"`
}

// DeleteAttachment, using the given attachment id, deletes an attachment.
// https://developers.pocketsmith.com/reference#delete_attachments-id.
func (c *Client) DeleteAttachment(
	ctx context.Context,
	options *DeleteAttachmentOptions,
) (err error) {

	// setup tracing.
	newCtx, span := otel.Tracer(c.tracerName).Start(ctx, "DeleteAttachment")
	defer span.End()

	// validate options.
	if err := c.validator.StructCtx(newCtx, options); err != nil {
		span.SetStatus(codes.Error, fmt.Sprintf("failed to validate options: %v", err))
		span.RecordError(err)
		return err
	}

	// delete attachment.
	_, err = c.sender(newCtx, senderRequest{
		method: http.MethodDelete,
		path:   fmt.Sprintf("/attachments/%v", options.AttachmentID),
	}, nil)
	if err != nil {
		span.SetStatus(codes.Error, fmt.Sprintf("failed to delete attachment: %v", err))
		span.RecordError(err)
		return err
	}
	return nil
}

// ListAttachmentsOptions defines the options for listing attachments for a user.
// Unassigned: 0 = all attachments (default), 1 = unassigned only.
type ListAttachmentsOptions struct {
	Unassigned int `json:"unassigned,omitempty"`
}

// ListAttachmentsForUsersOptions ...
type ListAttachmentsForUserOptions struct {
	UserID int `json:"-" validate:"required"`

	ListAttachmentsOptions
}

// ListAttachmentsForUser, using the given user id, lists the attachments for a user.
// https://developers.pocketsmith.com/reference#get_users-id-attachments.
func (c *Client) ListAttachmentsForUser(
	ctx context.Context,
	options *ListAttachmentsForUserOptions,
) (attachments Attachments, err error) {

	// setup tracing.
	newCtx, span := otel.Tracer(c.tracerName).Start(ctx, "ListAttachmentsForUser")
	defer span.End()

	// validate options.
	if err := c.validator.StructCtx(newCtx, options); err != nil {
		span.SetStatus(codes.Error, fmt.Sprintf("failed to validate options: %v", err))
		span.RecordError(err)
		return nil, err
	}

	// list attachments — pass the unassigned filter as a query param, not a
	// request body (GET requests do not have a body; the API would ignore it).
	queryMap := map[string]string{}
	if options.Unassigned != 0 {
		queryMap["unassigned"] = fmt.Sprintf("%v", options.Unassigned)
	}
	_, err = c.sender(newCtx, senderRequest{
		method:  http.MethodGet,
		path:    fmt.Sprintf("/users/%v/attachments", options.UserID),
		queries: setupQueries(&queryMap),
	}, &attachments)
	if err != nil {
		span.SetStatus(codes.Error, fmt.Sprintf("failed to list attachments: %v", err))
		span.RecordError(err)
		return nil, err
	}
	return attachments, nil
}

// ListAttachments, using the token attached to the client, lists
// the attachment for the authed user.
func (c *Client) ListAttachments(
	ctx context.Context,
	options *ListAttachmentsOptions,
) (Attachments, error) {

	// setup tracing.
	newCtx, span := otel.Tracer(c.tracerName).Start(ctx, "ListAttachmentsForAuthedUser")
	defer span.End()

	userID, err := c.authedUserID(newCtx)
	if err != nil {
		span.SetStatus(codes.Error, fmt.Sprintf("failed to get authed user: %v", err))
		span.RecordError(err)
		return nil, err
	}

	// list attachments.
	attachments, err := c.ListAttachmentsForUser(
		newCtx,
		&ListAttachmentsForUserOptions{UserID: userID, ListAttachmentsOptions: *options},
	)
	if err != nil {
		span.SetStatus(codes.Error, fmt.Sprintf("failed to list attachments: %v", err))
		span.RecordError(err)
		return nil, err
	}
	return attachments, nil
}

// AssignAttachmentToTransactionOptions defines the options for assigning an
// attachment to a transaction.
type AssignAttachmentToTransactionOptions struct {
	TransactionID int32 `json:"-"             validate:"required"`
	AttachmentID  int   `json:"attachment_id" validate:"required"`
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

	// validate options.
	if err := c.validator.StructCtx(newCtx, options); err != nil {
		span.SetStatus(codes.Error, fmt.Sprintf("failed to validate options: %v", err))
		span.RecordError(err)
		return nil, err
	}

	// assign attachment to transaction.
	_, err = c.sender(newCtx, senderRequest{
		method: http.MethodPost,
		path:   fmt.Sprintf("/transactions/%v/attachments", options.TransactionID),
		body:   options,
	}, &attachment)
	if err != nil {
		span.SetStatus(
			codes.Error,
			fmt.Sprintf("failed to assign attachment to transaction: %v", err),
		)
		span.RecordError(err)
		return nil, err
	}
	return attachment, nil
}
