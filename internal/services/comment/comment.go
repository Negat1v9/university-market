package commentservice

import (
	"context"
	"log/slog"
	"net/url"

	"github.com/Negat1v9/work-marketplace/internal/storage"
	mongoStore "github.com/Negat1v9/work-marketplace/internal/storage/mongo"
	filters "github.com/Negat1v9/work-marketplace/internal/storage/mongo/filter"
	commentmodel "github.com/Negat1v9/work-marketplace/model/comment"
	taskmodel "github.com/Negat1v9/work-marketplace/model/taskModel"
	httpresponse "github.com/Negat1v9/work-marketplace/pkg/httpResponse"
	"github.com/Negat1v9/work-marketplace/pkg/utils"
)

type CommentServiceImpl struct {
	log   *slog.Logger
	store storage.Store
}

func NewServiceComment(log *slog.Logger, store storage.Store) CommentService {
	return &CommentServiceImpl{
		log:   log,
		store: store,
	}
}

func (s *CommentServiceImpl) Create(ctx context.Context, userID string, comment *commentmodel.Comment) (string, error) {
	if err := beforeCreate(userID, comment); err != nil {
		return "", err
	}

	filter := filters.New().Add(filters.TaskByID(comment.TaskID)).Add(filters.TaskByCreator(userID))
	task, err := s.store.Task().FindProj(ctx, filter.Filters(), taskmodel.ProjOnCreateComment)
	switch {
	case err == mongoStore.ErrNoTask:
		return "", httpresponse.NewError(404, err.Error())
	case err != nil:
		s.log.Error("create comment", slog.String("err", err.Error()))
		return "", httpresponse.ServerError()
	case task.IsComment:
		return "", httpresponse.NewError(409, "alredy commented")
	case task.AssignedTo != comment.WorkerID:
		return "", httpresponse.NewError(409, "task.assigned_to is not "+comment.WorkerID)
	case task.Status != taskmodel.Completed:
		return "", httpresponse.NewError(406, "task.status is not "+string(taskmodel.Completed))
	}

	commentID, err := s.createCommentTrx(ctx, comment, task.ID)
	if err != nil {
		s.log.Error("create comment", slog.String("err", err.Error()))
		return "", httpresponse.ServerError()
	}

	return commentID, nil
}

func (s *CommentServiceImpl) UserComments(ctx context.Context, userID string, v url.Values) ([]commentmodel.Comment, error) {
	limit, err := utils.ConvertStringToInt64(v.Get("limit"))
	if err != nil || limit > 20 {
		limit = 20 // max limit
	}
	skip, err := utils.ConvertStringToInt64(v.Get("skip"))
	if err != nil {
		skip = 0
	}

	filter := filters.New().Add(filters.CommentByCreator(userID))
	comments, err := s.store.Comment().FindMany(ctx, filter.Filters(), limit, skip)
	switch {
	case err == mongoStore.ErrNoComment:
		return nil, httpresponse.NewError(404, err.Error())
	case err != nil:
		s.log.Error("user comments", slog.String("err", err.Error()))
		return nil, httpresponse.ServerError()
	}

	return comments, nil
}

func (s *CommentServiceImpl) WorkerComments(ctx context.Context, workerID string, v url.Values) ([]commentmodel.Comment, error) {
	limit, err := utils.ConvertStringToInt64(v.Get("limit"))
	if err != nil || limit > 20 {
		limit = 20 // max limit
	}
	skip, err := utils.ConvertStringToInt64(v.Get("skip"))
	if err != nil {
		skip = 0
	}

	filter := filters.New().Add(filters.CommentByWorker(workerID))
	comments, err := s.store.Comment().FindMany(ctx, filter.Filters(), limit, skip)
	switch {
	case err == mongoStore.ErrNoComment:
		return nil, httpresponse.NewError(404, err.Error())
	case err != nil:
		s.log.Error("user comments", slog.String("err", err.Error()))
		return nil, httpresponse.ServerError()
	}

	return comments, nil
}
