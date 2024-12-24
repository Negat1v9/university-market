package taskHttp

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Negat1v9/work-marketplace/internal/config"
	taskservice "github.com/Negat1v9/work-marketplace/internal/services/task"
	"github.com/Negat1v9/work-marketplace/internal/web/middleware"
	taskmodel "github.com/Negat1v9/work-marketplace/model/taskModel"
	httpresponse "github.com/Negat1v9/work-marketplace/pkg/httpResponse"
)

type TaskHandler struct {
	cfg     *config.WebCfg
	service taskservice.TaskService
}

func New(cfg *config.WebCfg, s taskservice.TaskService) *TaskHandler {
	return &TaskHandler{
		cfg:     cfg,
		service: s,
	}
}

func (h *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), h.cfg.CtxTimeOut)
	defer cancel()

	var newTask taskmodel.TaskMeta
	err := json.NewDecoder(r.Body).Decode(&newTask)
	if err != nil {
		httpresponse.ResponseError(w, 422, httpresponse.NewError(422, "invalid data"))
		return
	}

	defer r.Body.Close()

	userID := r.Context().Value(middleware.CtxUserIDKey).(string)

	taskID, err := h.service.Create(ctx, userID, &newTask)

	if err != nil {
		httpresponse.ResponseError(w, 500, err)
		return
	}

	httpresponse.Response(w, 201, taskmodel.NewTaskCreatedRes{TaskID: taskID})
}

func (h *TaskHandler) FindOne(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), h.cfg.CtxTimeOut)
	defer cancel()

	userID := r.Context().Value(middleware.CtxUserIDKey).(string)
	taskID := r.PathValue("id")

	task, err := h.service.FindOne(ctx, userID, taskID)
	if err != nil {
		httpresponse.ResponseError(w, 500, err)
		return
	}

	httpresponse.Response(w, 200, task)
}

func (h *TaskHandler) UpdateTaskMeta(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), h.cfg.CtxTimeOut)
	defer cancel()

	var upd taskmodel.UpdateTaskMeta
	err := json.NewDecoder(r.Body).Decode(&upd)
	if err != nil {
		httpresponse.ResponseError(w, 422, httpresponse.NewError(422, "invalid data"))
		return
	}

	defer r.Body.Close()

	taskID := r.PathValue("id")
	userID := r.Context().Value(middleware.CtxUserIDKey).(string)

	updTask, err := h.service.UpdateTaskMeta(ctx, taskID, userID, &upd)
	if err != nil {
		httpresponse.ResponseError(w, 500, err)
		return
	}

	httpresponse.Response(w, 200, updTask)
}

func (h *TaskHandler) RaiseTask(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), h.cfg.CtxTimeOut)
	defer cancel()

	userID := r.Context().Value(middleware.CtxUserIDKey).(string)
	taskID := r.PathValue("id")

	taskInfo, err := h.service.RaiseTask(ctx, taskID, userID)
	if err != nil {
		httpresponse.ResponseError(w, 500, err)
		return
	}

	httpresponse.Response(w, 200, taskInfo)
}

func (h *TaskHandler) FindUserTasks(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), h.cfg.CtxTimeOut)
	defer cancel()

	userID := r.Context().Value(middleware.CtxUserIDKey).(string)

	tasks, err := h.service.FindUserTasks(ctx, userID, r.URL.Query())
	if err != nil {
		httpresponse.ResponseError(w, 500, err)
		return
	}

	httpresponse.Response(w, 200, tasks)
}

func (h *TaskHandler) SelectWorker(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), h.cfg.CtxTimeOut)
	defer cancel()

	userID := r.Context().Value(middleware.CtxUserIDKey).(string)
	taskID := r.PathValue("id")
	workerID := r.PathValue("worker_id")

	taskInfo, err := h.service.SelectWorker(ctx, taskID, userID, workerID)
	if err != nil {
		httpresponse.ResponseError(w, 500, err)
		return
	}

	httpresponse.Response(w, 200, taskInfo)
}

func (h *TaskHandler) CompleteTask(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), h.cfg.CtxTimeOut)
	defer cancel()

	userID := r.Context().Value(middleware.CtxUserIDKey).(string)
	taskID := r.PathValue("id")

	taskInfo, err := h.service.CompleteTask(ctx, taskID, userID)
	if err != nil {
		httpresponse.ResponseError(w, 500, err)
		return
	}

	httpresponse.Response(w, 200, taskInfo)
}

func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), h.cfg.CtxTimeOut)
	defer cancel()

	userID := r.Context().Value(middleware.CtxUserIDKey).(string)
	taskID := r.PathValue("id")

	err := h.service.DeleteTask(ctx, taskID, userID)
	if err != nil {
		httpresponse.ResponseError(w, 500, err)
		return
	}

	httpresponse.Response(w, 200, nil)
}
