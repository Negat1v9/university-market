package workerservice

import (
	"context"
	"log/slog"

	filters "github.com/Negat1v9/work-marketplace/internal/storage/mongo/filter"
	respondmodel "github.com/Negat1v9/work-marketplace/model/respond"
	taskmodel "github.com/Negat1v9/work-marketplace/model/taskModel"
	usermodel "github.com/Negat1v9/work-marketplace/model/userModel"
	"go.mongodb.org/mongo-driver/mongo"
)

func (s *WorkerServiceImpl) respondOnTaskTrx(ctx context.Context, createrTaskTgID int64, worker *usermodel.User, task *taskmodel.Task, respond *respondmodel.Respond) error {
	session, err := s.store.StartSession()
	if err != nil {
		return err
	}

	err = session.StartTransaction()
	if err != nil {
		return err
	}

	err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		_, err = s.store.User().Edit(sc, filters.New().Add(filters.UserByID(worker.ID)).Filters(), worker)
		if err != nil {
			session.AbortTransaction(ctx)
			return err
		}

		updatedTask, err := s.store.Task().Update(sc, filters.New().Add(filters.TaskByID(task.ID)).Filters(), task)
		if err != nil {
			session.AbortTransaction(ctx)
			return err
		}

		_, err = s.store.Respond().Create(ctx, respond)
		if err != nil {
			session.AbortTransaction(ctx)
			return err
		}

		s.log.Debug("", slog.String("taskID", task.ID), slog.String("worker ID", worker.ID))

		err = s.tgClient.SendRespond(ctx, createrTaskTgID, worker.ID, updatedTask)
		if err != nil {
			session.AbortTransaction(ctx)
			return err
		}
		return session.CommitTransaction(ctx)
	})

	return err
}
