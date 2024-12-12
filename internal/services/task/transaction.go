package taskservice

import (
	"context"

	filters "github.com/Negat1v9/work-marketplace/internal/storage/mongo/filter"
	taskmodel "github.com/Negat1v9/work-marketplace/model/taskModel"
	tgbotmodel "github.com/Negat1v9/work-marketplace/model/tgBot"
	"go.mongodb.org/mongo-driver/mongo"
)

// Info: CreateTaskWithFilesTrx - creating a task if files are expected from the user
func (s *TaskServiceImpl) CreateTaskWithFilesTrx(ctx context.Context, tgUserID int64, task *taskmodel.Task) (string, error) {
	session, err := s.store.StartSession()
	if err != nil {
		return "", err
	}

	err = session.StartTransaction()
	if err != nil {
		return "", err
	}
	var taskID string
	err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		taskID, err = s.store.Task().Create(sc, task)
		if err != nil {
			session.AbortTransaction(ctx)
			return err
		}

		tgCmd := &tgbotmodel.UserCommand{
			ID:             tgUserID,
			ExpectedAction: tgbotmodel.WaitingForFiles,
			TaskID:         taskID,
		}

		err = s.store.TgCmd().Create(sc, tgCmd)
		if err != nil {
			session.AbortTransaction(ctx)
			return err
		}

		err = s.tgClient.WaitFiles(ctx, tgUserID)
		if err != nil {
			session.AbortTransaction(ctx)
			return err
		}

		return session.CommitTransaction(ctx)
	})

	return taskID, err
}

func (s *TaskServiceImpl) selectWorkerOnTaskTrx(ctx context.Context, tgWorkerID int64, upd *taskmodel.Task, tgCmd *tgbotmodel.UserCommand) (*taskmodel.Task, error) {
	session, err := s.store.StartSession()
	if err != nil {
		return nil, err
	}

	err = session.StartTransaction()
	if err != nil {
		return nil, err
	}

	var updatedTask *taskmodel.Task
	err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		updatedTask, err = s.store.Task().Update(sc, filters.New().Add(filters.TaskByID(upd.ID)).Filters(), upd)
		if err != nil {
			session.AbortTransaction(ctx)
			return err
		}
		err = s.store.TgCmd().Create(sc, tgCmd)
		if err != nil {
			session.AbortTransaction(ctx)
			return err
		}

		err = s.tgClient.SelectWorker(ctx, tgWorkerID, updatedTask)
		if err != nil {
			session.AbortTransaction(ctx)
			return err
		}

		return session.CommitTransaction(ctx)
	})

	return updatedTask, err
}
