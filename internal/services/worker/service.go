package workerservice

import (
	"context"
	"log/slog"
	"net/url"
	"time"

	taskservice "github.com/Negat1v9/work-marketplace/internal/services/task"
	"github.com/Negat1v9/work-marketplace/internal/storage"
	mongoStore "github.com/Negat1v9/work-marketplace/internal/storage/mongo"
	filters "github.com/Negat1v9/work-marketplace/internal/storage/mongo/filter"
	tgbot "github.com/Negat1v9/work-marketplace/internal/tgBot"
	respondmodel "github.com/Negat1v9/work-marketplace/model/respond"
	taskmodel "github.com/Negat1v9/work-marketplace/model/taskModel"
	usermodel "github.com/Negat1v9/work-marketplace/model/userModel"
	httpresponse "github.com/Negat1v9/work-marketplace/pkg/httpResponse"
	"github.com/Negat1v9/work-marketplace/pkg/utils"
)

type WorkerServiceImpl struct {
	log      *slog.Logger
	tgClient tgbot.WebTgClient
	store    storage.Store
}

func NewServiceWorker(log *slog.Logger, tgClient tgbot.WebTgClient, store storage.Store) WorkerService {
	return &WorkerServiceImpl{
		log:      log,
		tgClient: tgClient,
		store:    store,
	}
}

func (s *WorkerServiceImpl) Create(ctx context.Context, userID string, data *usermodel.WorkerCreate) (*usermodel.User, error) {
	filter := filters.NewCmplxFilter().Add(filters.UserByID(userID))
	user, err := s.store.User().FindProj(ctx, filter.Filters(), usermodel.AuthWorker)
	switch {
	case err == mongoStore.ErrNoUser:
		return nil, httpresponse.NewError(404, err.Error())
	case err != nil:
		s.log.Error("create worker", slog.String("err", err.Error()))
		return nil, httpresponse.ServerError()
	case user.Role == usermodel.Worker:
		return nil, httpresponse.NewError(409, "alredy exist")

	}

	upd := usermodel.User{
		PhoneNumber: data.PhoneNumber,
		Role:        usermodel.Worker,
		WorkerInfo:  usermodel.NewWorkerInfo(0), // TODO: start star balance ?
		UpdatedAt:   time.Now(),
	}

	filter = filters.NewCmplxFilter().Add(filters.UserByID(userID))
	user, err = s.store.User().Edit(ctx, filter.Filters(), &upd)
	if err != nil {
		s.log.Error("create worker", slog.String("err", err.Error()))
		return nil, httpresponse.ServerError()
	}

	return user, nil
}

func (s *WorkerServiceImpl) IsWorker(ctx context.Context, userID string) (bool, error) {
	filter := filters.NewCmplxFilter().Add(filters.UserByID(userID))
	user, err := s.store.User().FindProj(ctx, filter.Filters(), usermodel.AuthWorker)
	switch {
	case err == mongoStore.ErrNoUser:
		return false, httpresponse.NewError(404, err.Error())
	case err != nil:
		s.log.Error("create worker", slog.String("err", err.Error()))
		return false, httpresponse.ServerError()
	}

	return user.Role == usermodel.Worker, nil
}

func (s *WorkerServiceImpl) WorkerPublicInfo(ctx context.Context, workerID string) (*usermodel.WorkerInfoWithTaskRes, error) {
	worker, err := s.store.User().FindProj(ctx, filters.New().Add(filters.UserByID(workerID)).Filters(), usermodel.WorkerPublic)
	switch {
	case err == mongoStore.ErrNoUser:
		return nil, httpresponse.NewError(404, err.Error())
	case err != nil:
		s.log.Error("worker public info", slog.String("err", err.Error()))
		return nil, httpresponse.ServerError()
	case worker.Role != usermodel.Worker:
		return nil, httpresponse.NewError(404, mongoStore.ErrNoUser.Error())
	case worker.WorkerInfo == nil:
		s.log.Error("worker info nil", slog.String("id", workerID))
		return nil, httpresponse.ServerError()
	}

	return &usermodel.WorkerInfoWithTaskRes{
		ID:          workerID,
		UserName:    worker.Username,
		Karma:       worker.WorkerInfo.Karma,
		Description: worker.WorkerInfo.Description,
	}, nil
}
func (s *WorkerServiceImpl) Worker(ctx context.Context, workerID string) (*usermodel.User, error) {
	filter := filters.NewCmplxFilter().Add(filters.UserByID(workerID))
	worker, err := s.store.User().Find(ctx, filter.Filters())

	switch {
	case err == mongoStore.ErrNoUser:
		return nil, httpresponse.NewError(404, err.Error())
	case err != nil:
		s.log.Error("worker", slog.String("err", err.Error()))
		return nil, httpresponse.ServerError()
	case worker.Role != usermodel.Worker: // if user not worker
		return nil, httpresponse.NewError(403, "user is not worker")
	}

	return worker, nil
}

func (s *WorkerServiceImpl) Update(ctx context.Context, workerID string, data *usermodel.WorkerInfo) (*usermodel.User, error) {
	filter := filters.NewCmplxFilter().Add(filters.UserByID(workerID))
	worker, err := s.store.User().Find(ctx, filter.Filters())
	switch {
	case err == mongoStore.ErrNoUser:
		return nil, httpresponse.NewError(404, err.Error())
	case err != nil:
		s.log.Error("update worker", slog.String("err", err.Error()))
		return nil, httpresponse.ServerError()
	}

	updWorker := worker.WorkerInfo.FindUpdates(data)

	upd := usermodel.User{
		WorkerInfo: updWorker,
		UpdatedAt:  time.Now(),
	}

	filter = filters.NewCmplxFilter().Add(filters.UserByID(workerID))
	afterWorker, err := s.store.User().Edit(ctx, filter.Filters(), &upd)
	if err != nil {
		s.log.Error("update worker", slog.String("err", err.Error()))
		return nil, httpresponse.ServerError()
	}
	return afterWorker, nil
}

func (s *WorkerServiceImpl) AvailableTasks(ctx context.Context, v url.Values) ([]taskmodel.Task, error) {
	filter := taskservice.FindFilterTasks(v)
	filter.Add(filters.TaskByStatus(taskmodel.WaitingExecution))

	limit, err := utils.ConvertStringToInt64(v.Get("limit"))
	if err != nil || limit > 20 {
		limit = 20 // max limit
	}
	skip, err := utils.ConvertStringToInt64(v.Get("skip"))
	if err != nil {
		skip = 0
	}

	tasks, err := s.store.Task().FindMany(ctx, filter.Filters(), taskmodel.ManyTasks, limit, skip)
	switch {
	case err == mongoStore.ErrNoTask:
		return nil, httpresponse.NewError(404, err.Error())
	case err != nil:
		s.log.Error("available tasks worker", slog.String("err", err.Error()))
		return nil, httpresponse.ServerError()
	}

	return tasks, nil
}

// Info: TaskInfo - returns all information about the task and information about the cost of
// responding it. How much does it cost to respond for a given worker
func (s *WorkerServiceImpl) TaskInfo(ctx context.Context, workerID string, taskID string) (*taskmodel.InfoTaskRes, error) {
	filter := filters.NewCmplxFilter().
		Add(filters.TaskByID(taskID))

	task, err := s.store.Task().Find(ctx, filter.Filters())
	switch {
	case err == mongoStore.ErrNoTask:
		return nil, httpresponse.NewError(404, err.Error())
	case err != nil:
		s.log.Error("task info worker", slog.String("err", err.Error()))
		return nil, httpresponse.ServerError()
	}

	filter = filters.NewCmplxFilter().Add(filters.UserByID(workerID))

	workerInfo, err := s.store.User().FindProj(ctx, filter.Filters(), usermodel.OnlyWorkerInfo)
	if err != nil {
		s.log.Error("task info worker", slog.String("err", err.Error()))
		return nil, httpresponse.ServerError()
	}

	respondPrice := taskmodel.CalculateRespondStarPrice(task.Meta, workerInfo.WorkerInfo)
	if respondPrice == 0 {
		s.log.Error("task info worker", slog.String("err", "no info worker or task meta task ID - "+taskID))
		return nil, httpresponse.ServerError()
	}

	return &taskmodel.InfoTaskRes{Task: task, QuantityFiles: len(task.FilesID), RespondPrice: respondPrice}, nil
}

func (s *WorkerServiceImpl) RespondOnTask(ctx context.Context, workerID, taskID string) error {
	filter := filters.New().
		Add(filters.TaskByID(taskID)).
		Add(filters.TaskByStatus(taskmodel.WaitingExecution))

	task, err := s.store.Task().FindProj(ctx, filter.Filters(), taskmodel.ProjOnRespond)
	switch {
	case err == mongoStore.ErrNoTask:
		return httpresponse.NewError(404, err.Error())
	case err != nil:
		s.log.Error("respond on task", slog.String("err", err.Error()))
		return httpresponse.ServerError()
	}

	if checkWorkerAlredyRespond(workerID, task.Responds) {
		return httpresponse.NewError(409, "have already responded")
	}

	// receiving telegram user ID to send him a message
	createrTask, err := s.store.User().FindProj(
		ctx,
		filters.New().Add(filters.UserByID(task.CreatedBy)).Filters(),
		usermodel.OnlyTgID)
	if err != nil {
		s.log.Error("respond on task", slog.String("err", err.Error()))
		return httpresponse.ServerError()
	}

	filter = filters.New().Add(filters.UserByID(workerID))
	worker, err := s.store.User().FindProj(ctx, filter.Filters(), usermodel.OnlyWorkerInfo)
	if err != nil {
		s.log.Error("respond on task", slog.String("err", err.Error()))
		return httpresponse.ServerError()
	}

	respondPrice := taskmodel.CalculateRespondStarPrice(task.Meta, worker.WorkerInfo)
	if respondPrice == 0 {
		s.log.Error("respond on task", slog.String("err", "no info worker or task meta task ID - "+taskID))
		return httpresponse.ServerError()
	}
	if respondPrice > worker.WorkerInfo.StarsBalance {
		return httpresponse.NewError(406, "not enough funds")
	}

	worker.WorkerInfo.StarsBalance -= respondPrice
	updTask := &taskmodel.Task{
		ID:       task.ID,
		Responds: append(task.Responds, workerID),
	}

	respond := respondmodel.New(taskID, workerID, task.Meta.TaskType, respondPrice)

	err = s.respondOnTaskTrx(ctx, createrTask.TelegramID, worker, updTask, respond)
	if err != nil {
		s.log.Error("respond on task", slog.String("err", err.Error()))
		return httpresponse.ServerError()
	}

	return nil
}

func (s *WorkerServiceImpl) TasksResponded(ctx context.Context, workerID string, v url.Values) ([]taskmodel.Task, error) {
	filter := filters.New().Add(filters.TaskByAssigned(workerID))
	limit, err := utils.ConvertStringToInt64(v.Get("limit"))
	if err != nil || limit > 20 {
		limit = 20 // max limit
	}
	skip, err := utils.ConvertStringToInt64(v.Get("skip"))
	if err != nil {
		skip = 0
	}

	tasks, err := s.store.Task().FindMany(ctx, filter.Filters(), taskmodel.ManyTasks, limit, skip)
	switch {
	case err == mongoStore.ErrNoTask:
		return nil, httpresponse.NewError(404, err.Error())
	case err != nil:
		s.log.Error("tasks responded", slog.String("err", err.Error()))
		return nil, httpresponse.ServerError()
	}

	return tasks, nil
}
func (s *WorkerServiceImpl) Responds(ctx context.Context, workerID string, v url.Values) ([]respondmodel.Respond, error) {
	filter := filters.New().Add(filters.RespondByWorkerID(workerID))
	limit, err := utils.ConvertStringToInt64(v.Get("limit"))
	if err != nil || limit > 20 {
		limit = 20 // max limit
	}
	skip, err := utils.ConvertStringToInt64(v.Get("skip"))
	if err != nil {
		skip = 0
	}
	responds, err := s.store.Respond().FindMany(ctx, filter.Filters(), limit, skip)
	switch {
	case err == mongoStore.ErrNoRespond:
		return nil, httpresponse.NewError(404, err.Error())
	case err != nil:
		s.log.Error("find responds", slog.String("err", err.Error()))
		return nil, err
	}

	return responds, nil
}
