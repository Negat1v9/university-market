package commentservice

import (
	"context"

	filters "github.com/Negat1v9/work-marketplace/internal/storage/mongo/filter"
	commentmodel "github.com/Negat1v9/work-marketplace/model/comment"
	taskmodel "github.com/Negat1v9/work-marketplace/model/taskModel"
	"go.mongodb.org/mongo-driver/mongo"
)

func (s *CommentServiceImpl) createCommentTrx(ctx context.Context, comment *commentmodel.Comment, taskID string) (string, error) {
	session, err := s.store.StartSession()
	if err != nil {
		return "", err
	}

	err = session.StartTransaction()
	if err != nil {
		return "", err
	}
	var commentID string
	err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {

		commentID, err = s.store.Comment().Create(ctx, comment)
		if err != nil {
			session.AbortTransaction(ctx)
			return err
		}

		updTask := &taskmodel.Task{
			IsComment: true,
			CommentID: commentID,
		}

		_, err = s.store.Task().Update(sc, filters.New().Add(filters.TaskByID(taskID)).Filters(), updTask)
		if err != nil {
			session.AbortTransaction(ctx)
			return err
		}

		return session.CommitTransaction(ctx)
	})

	return commentID, err
}
