package taskservice

import (
	"context"
	"net/url"

	taskmodel "github.com/Negat1v9/work-marketplace/model/taskModel"
)

type TaskService interface {
	Create(ctx context.Context, userID string, meta *taskmodel.TaskMeta) (string, error)
	FindOne(ctx context.Context, userID, taskID string) (*taskmodel.InfoTaskRes, error)
	UpdateTaskMeta(ctx context.Context, taskID, userID string, data *taskmodel.UpdateTaskMeta) (*taskmodel.InfoTaskRes, error)
	FindUserTasks(ctx context.Context, userID string, v url.Values) ([]taskmodel.Task, error)
	SelectWorker(ctx context.Context, taskID, userID, workerID string) (*taskmodel.InfoTaskRes, error)
	CompleteTask(ctx context.Context, taskID, userID string) (*taskmodel.InfoTaskRes, error)
	DeleteTask(ctx context.Context, taskID, userID string) error
	AttachFiles(ctx context.Context, taskID string, fileID string) error
	PublishTask(ctx context.Context, taskID string) error
}
