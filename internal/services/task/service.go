package taskservice

import (
	"context"
	"log/slog"
	"net/url"
	"time"

	"github.com/Negat1v9/work-marketplace/internal/storage"
	mongoStore "github.com/Negat1v9/work-marketplace/internal/storage/mongo"
	filters "github.com/Negat1v9/work-marketplace/internal/storage/mongo/filter"
	tgbot "github.com/Negat1v9/work-marketplace/internal/tgBot"
	taskmodel "github.com/Negat1v9/work-marketplace/model/taskModel"
	tgbotmodel "github.com/Negat1v9/work-marketplace/model/tgBot"
	usermodel "github.com/Negat1v9/work-marketplace/model/userModel"
	httpresponse "github.com/Negat1v9/work-marketplace/pkg/httpResponse"
	"github.com/Negat1v9/work-marketplace/pkg/utils"
)

var (
	// minimum time that must pass after the last update to update the task
	deltaOnUpdateTask time.Duration = time.Minute * 60
)

type TaskServiceImpl struct {
	log      *slog.Logger
	tgClient tgbot.WebTgClient
	store    storage.Store
}

func NewServiceTask(log *slog.Logger, tgClient tgbot.WebTgClient, store storage.Store) TaskService {
	return &TaskServiceImpl{
		log:      log,
		tgClient: tgClient,
		store:    store,
	}
}

// Info: Create - create new task in database return taskID if created and return httpresponse.HttpError on error
func (s *TaskServiceImpl) Create(ctx context.Context, userID string, meta *taskmodel.TaskMeta) (string, error) {
	err := beforeCreateUpdate(meta)
	if err != nil {
		return "", err
	}
	if err := meta.ValidateMetaFields(); err != nil {
		return "", httpresponse.NewError(406, err.Error())
	}
	numberWeTasks, err := s.store.Task().Count(
		ctx,
		filters.New().Add(filters.TaskByCreator(userID)).Add(filters.TaskByStatus(taskmodel.WaitingExecution)).Filters(),
	)
	switch {
	case err != nil:
		s.log.Error("count we tasks on create", slog.String("err", err.Error()))
		return "", httpresponse.ServerError()
	case numberWeTasks >= 10: // limit
		return "", httpresponse.NewError(406, "too many unfinished tasks")
	}
	newTask := taskmodel.NewTask(userID, meta, createTags(meta))
	// send message in chat with a description of how to attach files and create task
	var taskID string
	switch meta.WithFiles {
	case true:
		taskID, err = s.createTaskWithFiles(ctx, userID, newTask)
		if err != nil {
			return "", err
		}

	case false:
		taskID, err = s.createTaskNoFiles(ctx, newTask)
		if err != nil {
			return "", err
		}
	}

	return taskID, nil
}

// Info: FindOne - finding a task by its ID, returns all information about the task
func (s *TaskServiceImpl) FindOne(ctx context.Context, userID, taskID string) (*taskmodel.InfoTaskRes, error) {
	task, err := s.store.Task().Find(
		ctx,
		filters.New().Add(filters.TaskByID(taskID)).Add(filters.TaskByCreator(userID)).Filters())
	switch {
	case err == mongoStore.ErrNoTask:
		return nil, httpresponse.NewError(404, err.Error())
	case err != nil:
		s.log.Debug("find task", slog.String("err", err.Error()))
		return nil, httpresponse.NewError(404, err.Error())
	case task.Status == taskmodel.Deleted:
		return nil, httpresponse.NewError(404, mongoStore.ErrNoTask.Error())
	}

	return &taskmodel.InfoTaskRes{Task: task, QuantityFiles: len(task.FilesID)}, nil
}

// Info: UpdateTaskMeta - updating the meta field of a task, returns the updated fields
// only the user who created the task can change it
func (s *TaskServiceImpl) UpdateTaskMeta(ctx context.Context, taskID, userID string, data *taskmodel.UpdateTaskMeta) (*taskmodel.InfoTaskRes, error) {
	if err := beforeCreateUpdate(&data.Meta); err != nil {
		return nil, err
	}

	if err := data.Meta.ValidateMetaFields(); err != nil {
		return nil, httpresponse.NewError(406, err.Error())
	}

	task, err := s.store.Task().FindProj(
		ctx,
		filters.New().Add(filters.TaskByID(taskID)).Add(filters.TaskByCreator(userID)).Filters(),
		taskmodel.ProjOnUpdateTask)
	switch {
	case err == mongoStore.ErrNoTask:
		return nil, httpresponse.NewError(404, err.Error())
	case err != nil:
		s.log.Debug("update task", slog.String("err", err.Error()))
		return nil, httpresponse.ServerError()
	case task.Status == taskmodel.Deleted:
		return nil, httpresponse.NewError(404, mongoStore.ErrNoTask.Error())
	case !canUpdateTaskTime(task.UpdatedAt, deltaOnUpdateTask):
		return nil, httpresponse.NewError(406, "last update recently")
	}
	// meta is a pointer check before working with the field
	if task.Meta == nil {
		task.Meta = &taskmodel.TaskMeta{}
	}

	upd := taskmodel.Task{
		Meta:      &data.Meta,
		Tags:      createTags(&data.Meta),
		UpdatedAt: time.Now().UTC(),
	}

	afterTask, err := s.store.Task().Update(
		ctx,
		filters.New().Add(filters.TaskByID(taskID)).Filters(),
		&upd)
	if err != nil {
		s.log.Debug("update task", slog.String("err", err.Error()))
		return nil, httpresponse.ServerError()
	}

	return &taskmodel.InfoTaskRes{Task: afterTask, QuantityFiles: len(afterTask.FilesID)}, nil
}

// Info: FindUserTasks - find all user tasks support quary filters
// return only tasks created by user with userID
func (s *TaskServiceImpl) FindUserTasks(ctx context.Context, userID string, v url.Values) ([]taskmodel.Task, error) {
	filter := FindFilterTasks(v)
	filter.Add(filters.TaskByCreator(userID))
	filter.Add(filters.TaskByNoDeleted())

	limit, err := utils.ConvertStringToInt64(v.Get("limit"))
	if err != nil || limit > 20 {
		limit = 20 // max limit
	}
	skip, err := utils.ConvertStringToInt64(v.Get("skip"))
	if err != nil {
		skip = 0
	}

	tasks, err := s.store.Task().FindMany(ctx, filter.Filters(), taskmodel.ProjManyTaskForUser, limit, skip)
	switch {
	case err == mongoStore.ErrNoTask:
		return nil, httpresponse.NewError(404, err.Error())
	case err != nil:
		s.log.Debug("update task status", slog.String("err", err.Error()))
		return nil, err
	}

	return tasks, nil
}

func (s *TaskServiceImpl) SelectWorker(ctx context.Context, taskID, userID, workerID string) (*taskmodel.InfoTaskRes, error) {
	task, err := s.store.Task().FindProj(
		ctx,
		filters.New().Add(filters.TaskByID(taskID)).Add(filters.TaskByCreator(userID)).Filters(),
		taskmodel.ProjOnRespond)
	switch {
	case err == mongoStore.ErrNoTask:
		return nil, httpresponse.NewError(404, err.Error())
	case err != nil:
		s.log.Error("select worker", slog.String("err", err.Error()))
		return nil, err
	case task.Status == taskmodel.Deleted:
		return nil, httpresponse.NewError(404, mongoStore.ErrNoTask.Error())
	case task.Status != taskmodel.WaitingExecution:
		return nil, httpresponse.NewError(406, "task.status is not "+string(taskmodel.WaitingExecution))
	}
	if err = CheckWorkerRespond(workerID, task.Responds); err != nil {
		return nil, err
	}

	workerInfo, err := s.store.User().FindProj(
		ctx,
		filters.New().Add(filters.UserByID(workerID)).Filters(),
		usermodel.OnlyTgID)
	switch {
	case err == mongoStore.ErrNoUser:
		return nil, httpresponse.NewError(404, err.Error())
	case err != nil:
		s.log.Error("select worker", slog.String("err", err.Error()))
		return nil, err
	}

	err = s.store.TgCmd().Delete(ctx, workerInfo.TelegramID)
	if err != nil && err != mongoStore.ErrNoTgCmd {
		s.log.Error("create task", slog.String("err", err.Error()))
		return nil, httpresponse.ServerError()
	}
	tgCmd := &tgbotmodel.UserCommand{
		ID:             workerInfo.TelegramID,
		ExpectedAction: tgbotmodel.WorkerShareContact,
		TaskID:         task.ID,
		UserID:         task.CreatedBy,
	}

	upd := &taskmodel.Task{
		ID:         taskID,
		AssignedTo: workerID,
		Status:     taskmodel.InProgress,
	}

	afterUpdate, err := s.selectWorkerOnTaskTrx(ctx, workerInfo.TelegramID, upd, tgCmd)

	if err != nil {
		s.log.Error("select worker", slog.String("err", err.Error()))
		return nil, err
	}

	return &taskmodel.InfoTaskRes{Task: afterUpdate, QuantityFiles: len(afterUpdate.FilesID)}, nil
}

func (s *TaskServiceImpl) CompleteTask(ctx context.Context, taskID, userID string) (*taskmodel.InfoTaskRes, error) {
	task, err := s.store.Task().FindProj(
		ctx,
		filters.New().Add(filters.TaskByID(taskID)).Add(filters.TaskByCreator(userID)).Filters(),
		taskmodel.OnlyStatus)
	switch {
	case err == mongoStore.ErrNoTask:
		return nil, httpresponse.NewError(404, err.Error())
	case err != nil:
		s.log.Error("select worker", slog.String("err", err.Error()))
		return nil, err
	case task.Status == taskmodel.Deleted:
		return nil, httpresponse.NewError(404, mongoStore.ErrNoTask.Error())
	case task.Status != taskmodel.InProgress:
		return nil, httpresponse.NewError(406, "task.status is not "+string(taskmodel.InProgress))
	}
	upd := &taskmodel.Task{
		Status:      taskmodel.Completed,
		UpdatedAt:   time.Now().UTC(),
		ComplitedAT: time.Now().UTC(),
	}

	afterTask, err := s.store.Task().Update(ctx, filters.New().Add(filters.TaskByID(taskID)).Filters(), upd)
	if err != nil {
		s.log.Error("complete task", slog.String("err", err.Error()))
		return nil, err
	}
	return &taskmodel.InfoTaskRes{Task: afterTask, QuantityFiles: len(afterTask.FilesID)}, nil
}

func (s *TaskServiceImpl) DeleteTask(ctx context.Context, taskID, userID string) error {
	numberDeletedTasks, err := s.store.Task().Count(
		ctx,
		filters.New().Add(filters.TaskByCreator(userID)).Add(filters.TaskByIsDeleted()).Filters())
	switch {
	case err != nil:
		s.log.Error("count for delete task", slog.String("ere", err.Error()))
		return httpresponse.ServerError()
	case numberDeletedTasks >= 3:
		return httpresponse.NewError(429, "too many deletion attempts")
	}

	task, err := s.store.Task().FindProj(
		ctx,
		filters.New().Add(filters.TaskByID(taskID)).Add(filters.TaskByCreator(userID)).Filters(),
		taskmodel.OnlyStatus,
	)
	switch {
	case err == mongoStore.ErrNoTask:
		return httpresponse.NewError(404, err.Error())
	case err != nil:
		s.log.Error("delete task", slog.String("err", err.Error()))
		return httpresponse.ServerError()
	case task.Status == taskmodel.Deleted:
		return httpresponse.NewError(404, mongoStore.ErrNoTask.Error())
	}

	delTime := time.Now().UTC()
	delTask := taskmodel.Task{
		ID:       taskID,
		Status:   taskmodel.Deleted,
		DeleteAt: &delTime,
	}
	_, err = s.store.Task().Update(
		ctx,
		filters.New().Add(filters.TaskByID(taskID)).Filters(),
		&delTask)

	if err != nil {
		s.log.Error("delete task", slog.String("err", err.Error()))
		return httpresponse.ServerError()
	}
	return nil
}

// AttachFiles - calls from bot attach files to task
func (s *TaskServiceImpl) AttachFiles(ctx context.Context, taskID string, fileID string) error {
	task, err := s.store.Task().FindProj(
		ctx, filters.New().Add(filters.TaskByID(taskID)).Filters(),
		taskmodel.ProjOnAttachFiles)

	switch {
	case err == mongoStore.ErrNoTask:
		return httpresponse.NewError(404, err.Error())
	case err != nil:
		s.log.Error("attach files", slog.String("err", err.Error()))
		return err
	case task.Status == taskmodel.Deleted:
		return httpresponse.NewError(404, mongoStore.ErrNoTask.Error())
	case task.Status != taskmodel.Pending:
		return httpresponse.NewError(409, "cannot be changed")
	case len(task.FilesID) == 5:
		return httpresponse.NewError(406, "the largest number files is 5")
	}

	updTask := &taskmodel.Task{
		FilesID:   append(task.FilesID, fileID),
		UpdatedAt: time.Now().UTC(),
	}
	_, err = s.store.Task().Update(ctx, filters.New().Add(filters.TaskByID(taskID)).Filters(), updTask)
	if err != nil {
		s.log.Error("attach files", slog.String("err", err.Error()))
		return err
	}
	return nil
}

func (s *TaskServiceImpl) PublishTask(ctx context.Context, taskID string) error {
	task, err := s.store.Task().FindProj(
		ctx,
		filters.New().Add(filters.TaskByID(taskID)).Filters(),
		taskmodel.OnlyStatus)

	switch {
	case err == mongoStore.ErrNoTask:
		return httpresponse.NewError(404, err.Error())
	case err != nil:
		s.log.Error("publich task", slog.String("err", err.Error()))
		return err
	case task.Status == taskmodel.Deleted:
		return httpresponse.NewError(404, mongoStore.ErrNoTask.Error())
	case task.Status != taskmodel.Pending:
		return httpresponse.NewError(409, "alredy publiched")
	}

	upd := taskmodel.Task{
		Status: taskmodel.WaitingExecution,
	}

	_, err = s.store.Task().Update(
		ctx,
		filters.New().Add(filters.TaskByID(taskID)).Filters(),
		&upd)
	switch {
	case err == mongoStore.ErrNoTask:
		return httpresponse.NewError(404, err.Error())
	case err != nil:
		s.log.Error("publich task", slog.String("err", err.Error()))
		return err
	}

	return nil
}

func (s *TaskServiceImpl) createTaskWithFiles(ctx context.Context, userID string, task *taskmodel.Task) (string, error) {
	numberPendingTasks, err := s.store.Task().Count(
		ctx,
		filters.New().Add(filters.TaskByCreator(userID)).Add(filters.TaskByStatus(taskmodel.Pending)).Filters(),
	)
	switch {
	case err != nil:
		s.log.Error("count pending task on create", slog.String("err", err.Error()))
		return "", httpresponse.ServerError()
	case numberPendingTasks >= 1:
		return "", httpresponse.NewError(409, "task with status pending exists")
	}

	user, err := s.store.User().FindProj(ctx, filters.New().Add(filters.UserByID(userID)).Filters(), usermodel.OnlyTgID)
	if err != nil {
		s.log.Error("create task", slog.String("err", err.Error()))
		return "", httpresponse.ServerError()
	}
	// deleting a command if it exists
	err = s.store.TgCmd().Delete(ctx, user.TelegramID)
	if err != nil && err != mongoStore.ErrNoTgCmd {
		s.log.Error("create task", slog.String("err", err.Error()))
		return "", httpresponse.ServerError()
	}

	taskID, err := s.CreateTaskWithFilesTrx(ctx, user.TelegramID, task)
	if err != nil {
		s.log.Error("create task", slog.String("err", err.Error()))
		return "", httpresponse.ServerError()
	}
	return taskID, nil
}

func (s *TaskServiceImpl) createTaskNoFiles(ctx context.Context, task *taskmodel.Task) (string, error) {
	task.Status = taskmodel.WaitingExecution
	taskID, err := s.store.Task().Create(ctx, task)
	if err != nil {
		s.log.Error("create task", slog.String("err", err.Error()))
		return "", httpresponse.ServerError()
	}
	return taskID, nil
}
