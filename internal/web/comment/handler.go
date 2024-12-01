package commentHttp

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Negat1v9/work-marketplace/internal/config"
	commentservice "github.com/Negat1v9/work-marketplace/internal/services/comment"
	"github.com/Negat1v9/work-marketplace/internal/web/middleware"
	commentmodel "github.com/Negat1v9/work-marketplace/model/comment"
	httpresponse "github.com/Negat1v9/work-marketplace/pkg/httpResponse"
)

type CommentHandler struct {
	cfg     *config.WebCfg
	service commentservice.CommentService
}

func New(cfg *config.WebCfg, s commentservice.CommentService) *CommentHandler {
	return &CommentHandler{
		cfg:     cfg,
		service: s,
	}
}

func (h *CommentHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), h.cfg.CtxTimeOut)
	defer cancel()

	var comment commentmodel.Comment
	err := json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {

		httpresponse.ResponseError(w, 422, httpresponse.NewError(422, "invalid data"))
		return
	}

	defer r.Body.Close()

	userID := r.Context().Value(middleware.CtxUserIDKey).(string)

	commentID, err := h.service.Create(ctx, userID, &comment)
	if err != nil {
		httpresponse.ResponseError(w, 500, err)
		return
	}

	httpresponse.Response(w, 201, commentmodel.CommentCreateRes{CommentID: commentID})
}

func (h *CommentHandler) UserComments(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), h.cfg.CtxTimeOut)
	defer cancel()

	userID := r.Context().Value(middleware.CtxUserIDKey).(string)

	comments, err := h.service.UserComments(ctx, userID, r.URL.Query())
	if err != nil {
		httpresponse.ResponseError(w, 500, err)
		return
	}

	httpresponse.Response(w, 200, comments)
}

func (h *CommentHandler) UserWorkerComments(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), h.cfg.CtxTimeOut)
	defer cancel()

	workerID := r.PathValue("id")

	comments, err := h.service.WorkerComments(ctx, workerID, r.URL.Query())
	if err != nil {
		httpresponse.ResponseError(w, 500, err)
		return
	}

	httpresponse.Response(w, 200, comments)
}

func (h *CommentHandler) WorkerComments(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), h.cfg.CtxTimeOut)
	defer cancel()

	workerID := r.Context().Value(middleware.CtxUserIDKey).(string)

	comments, err := h.service.WorkerComments(ctx, workerID, r.URL.Query())
	if err != nil {
		httpresponse.ResponseError(w, 500, err)
		return
	}

	httpresponse.Response(w, 200, comments)
}
