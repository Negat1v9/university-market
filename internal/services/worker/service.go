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
	tgvalidation "github.com/Negat1v9/work-marketplace/pkg/tgValidation"
	"github.com/Negat1v9/work-marketplace/pkg/utils"
)

type WorkerServiceImpl struct {
	botToken string
	log      *slog.Logger
	tgClient tgbot.WebTgClient
	store    storage.Store
}

func NewServiceWorker(log *slog.Logger, tgClient tgbot.WebTgClient, store storage.Store, botToken string) WorkerService {
	return &WorkerServiceImpl{
		botToken: botToken,
		log:      log,
		tgClient: tgClient,
		store:    store,
	}
}

func (s *WorkerServiceImpl) Create(ctx context.Context, userID string, data *usermodel.WorkerCreate) (*usermodel.User, error) {
	user, err := s.store.User().FindProj(
		ctx,
		filters.New().Add(filters.UserByID(userID)).Filters(),
		usermodel.ProjCreateWorker)
	switch {
	case err == mongoStore.ErrNoUser:
		return nil, httpresponse.NewError(404, err.Error())
	case err != nil:
		s.log.Error("create worker", slog.String("err", err.Error()))
		return nil, httpresponse.ServerError()
	case user.Role == usermodel.Worker:
		return nil, httpresponse.NewError(409, "alredy exist")
	}

	err = tgvalidation.ValidateInitData(data.Response, s.botToken)
	if err != nil {
		return nil, httpresponse.NewError(401, "forbiden")
	}
	responseReqContact, err := tgvalidation.ParsePhoneNumber(data.Response)
	if err != nil {
		return nil, httpresponse.NewError(400, err.Error())
	}
	if responseReqContact.UserID != user.TelegramID {
		return nil, httpresponse.NewError(401, "forbiden")
	}
	fullName := ""
	if user.WorkerInfo != nil {
		fullName = user.WorkerInfo.FullName
	}

	upd := usermodel.User{
		PhoneNumber: responseReqContact.PhoneNumber,
		Role:        usermodel.Worker,
		WorkerInfo:  usermodel.NewWorkerInfo(fullName),
		UpdatedAt:   time.Now(),
	}

	user, err = s.store.User().Edit(
		ctx,
		filters.New().Add(filters.UserByID(userID)).Filters(),
		&upd)
	if err != nil {
		s.log.Error("create worker", slog.String("err", err.Error()))
		return nil, httpresponse.ServerError()
	}

	return user, nil
}

func (s *WorkerServiceImpl) IsWorker(ctx context.Context, userID string) (bool, error) {
	user, err := s.store.User().FindProj(
		ctx,
		filters.New().Add(filters.UserByID(userID)).Filters(),
		usermodel.AuthWorker)
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
	worker, err := s.store.User().FindProj(
		ctx,
		filters.New().Add(filters.UserByID(workerID)).Filters(),
		usermodel.WorkerPublic)
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

	rating, err := s.store.Comment().CountWorkerLikesDislikes(ctx, workerID)
	if err != nil {
		s.log.Error("worker public info", slog.String("err", err.Error()))
	}

	return &usermodel.WorkerInfoWithTaskRes{
		ID:          workerID,
		FullName:    worker.WorkerInfo.FullName,
		Rating:      rating,
		Education:   worker.WorkerInfo.Education,
		Experience:  worker.WorkerInfo.Experience,
		Description: worker.WorkerInfo.Description,
	}, nil
}
func (s *WorkerServiceImpl) Worker(ctx context.Context, workerID string) (*usermodel.User, error) {
	worker, err := s.store.User().Find(
		ctx,
		filters.New().Add(filters.UserByID(workerID)).Filters())

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
	if err := data.ValidateFields(); err != nil {
		return nil, httpresponse.NewError(406, err.Error())
	}
	worker, err := s.store.User().Find(
		ctx,
		filters.New().Add(filters.UserByID(workerID)).Filters())
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

	afterWorker, err := s.store.User().Edit(
		ctx,
		filters.New().Add(filters.UserByID(workerID)).Filters(),
		&upd)
	if err != nil {
		s.log.Error("update worker", slog.String("err", err.Error()))
		return nil, httpresponse.ServerError()
	}
	return afterWorker, nil
}

func (s *WorkerServiceImpl) AvailableTasks(ctx context.Context, workerID string, v url.Values) ([]taskmodel.Task, error) {
	v.Del("status")
	filter := taskservice.FindFilterTasks(v)
	filter.Add(filters.TaskByStatus(taskmodel.WaitingExecution))
	filter.Add(filters.TaskByNotCreator(workerID))

	limit, err := utils.ConvertStringToInt64(v.Get("limit"))
	if err != nil || limit > 20 {
		limit = 20 // max limit
	}
	skip, err := utils.ConvertStringToInt64(v.Get("skip"))
	if err != nil {
		skip = 0
	}

	tasks, err := s.store.Task().FindMany(
		ctx,
		filter.Filters(),
		taskmodel.ManyTasks,
		limit,
		skip)
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
	task, err := s.store.Task().Find(
		ctx,
		filters.New().Add(filters.TaskByID(taskID)).Filters())
	switch {
	case err == mongoStore.ErrNoTask:
		return nil, httpresponse.NewError(404, err.Error())
	case err != nil:
		s.log.Error("task info worker", slog.String("err", err.Error()))
		return nil, httpresponse.ServerError()
	case task.Status == taskmodel.Deleted:
		return nil, httpresponse.NewError(404, mongoStore.ErrNoTask.Error())
	}

	workerInfo, err := s.store.User().FindProj(
		ctx,
		filters.New().Add(filters.UserByID(workerID)).Filters(),
		usermodel.ProjOnlyBalance)
	if err != nil {
		s.log.Error("task info worker", slog.String("err", err.Error()))
		return nil, httpresponse.ServerError()
	}

	respondPrice := taskmodel.CalculateRespondStarPrice(task.Meta, workerInfo.Balance)
	if respondPrice == 0 {
		s.log.Error("task info worker", slog.String("err", "no info worker or task meta task ID - "+taskID))
		return nil, httpresponse.ServerError()
	}

	return &taskmodel.InfoTaskRes{Task: task, QuantityFiles: len(task.FilesID), RespondPrice: respondPrice}, nil
}

func (s *WorkerServiceImpl) SendTaskFiles(ctx context.Context, workerID, taskID string) error {
	task, err := s.store.Task().FindProj(
		ctx,
		filters.New().Add(filters.TaskByID(taskID)).Filters(),
		taskmodel.ProjOnSendFiles)
	switch {
	case err == mongoStore.ErrNoTask:
		return httpresponse.NewError(404, err.Error())
	case err != nil:
		s.log.Error("send task files", slog.String("err", err.Error()))
		return httpresponse.ServerError()
	case task.Status == taskmodel.Deleted: // task is deleted
		return httpresponse.NewError(404, mongoStore.ErrNoTask.Error())
	case task.CreatedBy == workerID: // worker is creater
		return httpresponse.NewError(409, "you are the creator of task")
	case task.Status == taskmodel.Pending: // task not ready
		return httpresponse.NewError(404, mongoStore.ErrNoTask.Error())
	case len(task.FilesID) == 0: // no files to send
		return httpresponse.NewError(404, "no files")
	}

	worker, err := s.store.User().FindProj(
		ctx,
		filters.New().Add(filters.UserByID(workerID)).Filters(),
		usermodel.OnlyTgID)
	if err != nil {
		s.log.Error("send task files", slog.String("err", err.Error()))
		return httpresponse.ServerError()
	}

	if err = checkFilesAlredySend(workerID, task.FilesSend); err != nil {
		return httpresponse.NewError(406, err.Error())
	}

	updTask := taskmodel.Task{
		ID:        task.ID,
		FilesSend: append(task.FilesSend, workerID),
	}

	err = s.sendTaskFilesTrx(ctx, worker.TelegramID, &updTask, task.FilesID)
	if err != nil {
		s.log.Error("send task files", slog.String("err", err.Error()))
		return httpresponse.ServerError()
	}
	return nil
}

func (s *WorkerServiceImpl) RespondOnTask(ctx context.Context, workerID, taskID string) error {
	task, err := s.store.Task().FindProj(
		ctx,
		filters.New().Add(filters.TaskByID(taskID)).Add(filters.TaskByStatus(taskmodel.WaitingExecution)).Filters(),
		taskmodel.ProjOnRespond)
	switch {
	case err == mongoStore.ErrNoTask:
		return httpresponse.NewError(404, err.Error())
	case err != nil:
		s.log.Error("respond on task", slog.String("err", err.Error()))
		return httpresponse.ServerError()
	case task.Status == taskmodel.Deleted: // task is deleted
		return httpresponse.NewError(404, mongoStore.ErrNoTask.Error())
	case task.CreatedBy == workerID:
		return httpresponse.NewError(409, "you are the creator of task")
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

	worker, err := s.store.User().FindProj(
		ctx,
		filters.New().Add(filters.UserByID(workerID)).Filters(),
		usermodel.ProjOnlyBalance)
	if err != nil {
		s.log.Error("respond on task", slog.String("err", err.Error()))
		return httpresponse.ServerError()
	}

	respondPrice := taskmodel.CalculateRespondStarPrice(task.Meta, worker.Balance)
	if respondPrice == 0 {
		s.log.Error("respond on task", slog.String("err", "no info worker or task meta task ID - "+taskID))
		return httpresponse.ServerError()
	}
	if respondPrice > worker.Balance.StarsBalance {
		return httpresponse.NewError(406, "not enough funds")
	}

	worker.Balance.StarsBalance -= respondPrice

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
	limit, err := utils.ConvertStringToInt64(v.Get("limit"))
	if err != nil || limit > 20 {
		limit = 20 // max limit
	}
	skip, err := utils.ConvertStringToInt64(v.Get("skip"))
	if err != nil {
		skip = 0
	}

	tasks, err := s.store.Task().FindMany(
		ctx,
		filters.New().Add(filters.TaskByAssigned(workerID)).Filters(),
		taskmodel.ManyTasks,
		limit,
		skip)
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
	limit, err := utils.ConvertStringToInt64(v.Get("limit"))
	if err != nil || limit > 20 {
		limit = 20 // max limit
	}
	skip, err := utils.ConvertStringToInt64(v.Get("skip"))
	if err != nil {
		skip = 0
	}
	responds, err := s.store.Respond().FindMany(
		ctx,
		filters.New().Add(filters.RespondByWorkerID(workerID)).Filters(),
		limit,
		skip)
	switch {
	case err == mongoStore.ErrNoRespond:
		return nil, httpresponse.NewError(404, err.Error())
	case err != nil:
		s.log.Error("find responds", slog.String("err", err.Error()))
		return nil, err
	}

	return responds, nil
}
