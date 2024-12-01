package commentservice

import (
	"context"
	"net/url"

	commentmodel "github.com/Negat1v9/work-marketplace/model/comment"
)

type CommentService interface {
	Create(ctx context.Context, userID string, comment *commentmodel.Comment) (string, error)
	UserComments(ctx context.Context, userID string, v url.Values) ([]commentmodel.Comment, error)
	WorkerComments(ctx context.Context, workerID string, v url.Values) ([]commentmodel.Comment, error)
}
