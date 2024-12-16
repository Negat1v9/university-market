package workerHttp

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Negat1v9/work-marketplace/internal/config"
	workerservice "github.com/Negat1v9/work-marketplace/internal/services/worker"
	"github.com/Negat1v9/work-marketplace/internal/web/middleware"
	usermodel "github.com/Negat1v9/work-marketplace/model/userModel"
	httpresponse "github.com/Negat1v9/work-marketplace/pkg/httpResponse"
)

type WorkerHandler struct {
	cfg     *config.WebCfg
	service workerservice.WorkerService
}

func New(cfg *config.WebCfg, s workerservice.WorkerService) *WorkerHandler {
	return &WorkerHandler{
		cfg:     cfg,
		service: s,
	}
}

func (h *WorkerHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), h.cfg.CtxTimeOut)
	defer cancel()

	userID := r.Context().Value(middleware.CtxUserIDKey).(string)

	var data usermodel.WorkerCreate
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		httpresponse.ResponseError(w, 422, httpresponse.NewError(422, "invalid data"))
		return
	}

	defer r.Body.Close()

	worker, err := h.service.Create(ctx, userID, &data)
	if err != nil {
		httpresponse.ResponseError(w, 500, err)
		return
	}

	httpresponse.Response(w, 201, worker)

}

func (h *WorkerHandler) IsWorker(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), h.cfg.CtxTimeOut)
	defer cancel()

	userID := r.Context().Value(middleware.CtxUserIDKey).(string)

	isWorker, err := h.service.IsWorker(ctx, userID)
	if err != nil {
		httpresponse.ResponseError(w, 500, err)
		return
	}

	httpresponse.Response(w, 200, usermodel.IsWorkerRes{
		UserID:   userID,
		IsWorker: isWorker,
	})

}

func (h *WorkerHandler) WorkerPublicInfo(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), h.cfg.CtxTimeOut)
	defer cancel()

	workerID := r.PathValue("id")

	workerInfo, err := h.service.WorkerPublicInfo(ctx, workerID)
	if err != nil {
		httpresponse.ResponseError(w, 500, err)
		return
	}

	httpresponse.Response(w, 200, workerInfo)
}

func (h *WorkerHandler) Worker(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), h.cfg.CtxTimeOut)
	defer cancel()

	workerID := r.Context().Value(middleware.CtxUserIDKey).(string)

	worker, err := h.service.Worker(ctx, workerID)
	if err != nil {
		httpresponse.ResponseError(w, 500, err)
		return
	}

	httpresponse.Response(w, 200, worker)
}

func (h *WorkerHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), h.cfg.CtxTimeOut)
	defer cancel()

	var data usermodel.WorkerInfo

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		httpresponse.ResponseError(w, 422, httpresponse.NewError(422, "invalid data"))
		return
	}

	defer r.Body.Close()

	workerID := r.Context().Value(middleware.CtxUserIDKey).(string)

	worker, err := h.service.Update(ctx, workerID, &data)
	if err != nil {
		httpresponse.ResponseError(w, 500, err)
		return
	}

	httpresponse.Response(w, 200, worker)
}

func (h *WorkerHandler) AvailableTasks(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), h.cfg.CtxTimeOut)
	defer cancel()

	workerID := r.Context().Value(middleware.CtxUserIDKey).(string)

	tasks, err := h.service.AvailableTasks(ctx, workerID, r.URL.Query())
	if err != nil {
		httpresponse.ResponseError(w, 500, err)
		return
	}

	httpresponse.Response(w, 200, tasks)
}

func (h *WorkerHandler) TaskInfo(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), h.cfg.CtxTimeOut)
	defer cancel()

	taskID := r.PathValue("id")
	workerID := r.Context().Value(middleware.CtxUserIDKey).(string)

	task, err := h.service.TaskInfo(ctx, workerID, taskID)
	if err != nil {
		httpresponse.ResponseError(w, 500, err)
		return
	}

	httpresponse.Response(w, 200, task)
}

func (h *WorkerHandler) TaskFiles(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), h.cfg.CtxTimeOut)
	defer cancel()

	taskID := r.PathValue("id")
	workerID := r.Context().Value(middleware.CtxUserIDKey).(string)

	err := h.service.SendTaskFiles(ctx, workerID, taskID)
	if err != nil {
		httpresponse.ResponseError(w, 500, err)
		return
	}

	httpresponse.Response(w, 200, nil)
}

func (h *WorkerHandler) RespondOnTask(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), h.cfg.CtxTimeOut)
	defer cancel()

	taskID := r.PathValue("id")
	workerID := r.Context().Value(middleware.CtxUserIDKey).(string)

	err := h.service.RespondOnTask(ctx, workerID, taskID)
	if err != nil {
		httpresponse.ResponseError(w, 500, err)
		return
	}

	httpresponse.Response(w, 200, nil)
}

func (h *WorkerHandler) TasksByResponds(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), h.cfg.CtxTimeOut)
	defer cancel()

	workerID := r.Context().Value(middleware.CtxUserIDKey).(string)

	tasks, err := h.service.TasksResponded(ctx, workerID, r.URL.Query())
	if err != nil {
		httpresponse.ResponseError(w, 500, err)
		return
	}

	httpresponse.Response(w, 200, tasks)
}

func (h *WorkerHandler) Responds(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), h.cfg.CtxTimeOut)
	defer cancel()

	workerID := r.Context().Value(middleware.CtxUserIDKey).(string)

	responds, err := h.service.Responds(ctx, workerID, r.URL.Query())
	if err != nil {
		httpresponse.ResponseError(w, 500, err)
		return
	}

	httpresponse.Response(w, 200, responds)
}
