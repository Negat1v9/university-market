package workerservice

import (
	"context"
	"net/url"

	respondmodel "github.com/Negat1v9/work-marketplace/model/respond"
	taskmodel "github.com/Negat1v9/work-marketplace/model/taskModel"
	usermodel "github.com/Negat1v9/work-marketplace/model/userModel"
)

type WorkerService interface {
	Create(ctx context.Context, userID string, data *usermodel.WorkerCreate) (*usermodel.User, error)
	IsWorker(ctx context.Context, userID string) (bool, error)
	WorkerPublicInfo(ctx context.Context, workerID string) (*usermodel.WorkerInfoWithTaskRes, error)
	Worker(ctx context.Context, workerID string) (*usermodel.User, error)
	Update(ctx context.Context, workerID string, data *usermodel.WorkerInfo) (*usermodel.User, error)
	AvailableTasks(ctx context.Context, v url.Values) ([]taskmodel.Task, error)
	TaskInfo(ctx context.Context, workerID string, taskID string) (*taskmodel.InfoTaskRes, error)
	SendTaskFiles(ctx context.Context, workerID, taskID string) error
	RespondOnTask(ctx context.Context, workerID, taskID string) error
	TasksResponded(ctx context.Context, workerID string, v url.Values) ([]taskmodel.Task, error)
	Responds(ctx context.Context, workerID string, v url.Values) ([]respondmodel.Respond, error)
}
