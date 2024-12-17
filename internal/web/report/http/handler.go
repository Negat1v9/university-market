package reportHttp

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Negat1v9/work-marketplace/internal/config"
	reportservice "github.com/Negat1v9/work-marketplace/internal/services/report"
	"github.com/Negat1v9/work-marketplace/internal/web/middleware"
	reportmodel "github.com/Negat1v9/work-marketplace/model/report"
	httpresponse "github.com/Negat1v9/work-marketplace/pkg/httpResponse"
)

type ReportHandler struct {
	cfg     *config.WebCfg
	service reportservice.ReportService
}

func New(cfg *config.WebCfg, s reportservice.ReportService) *ReportHandler {
	return &ReportHandler{
		cfg:     cfg,
		service: s,
	}
}

func (h *ReportHandler) CreateReportOnWorker(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), h.cfg.CtxTimeOut)
	defer cancel()

	var newReport reportmodel.NewReportReq
	err := json.NewDecoder(r.Body).Decode(&newReport)
	if err != nil {
		httpresponse.ResponseError(w, 422, httpresponse.NewError(422, "invalid data"))
		return
	}
	defer r.Body.Close()

	userID := r.Context().Value(middleware.CtxUserIDKey).(string)

	id, err := h.service.CreateReportOnWorker(ctx, userID, &newReport)
	if err != nil {
		httpresponse.ResponseError(w, 500, err)
		return
	}

	httpresponse.Response(w, 201, reportmodel.NewReportRes{ID: id})
}

func (h *ReportHandler) CreateReportOnUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), h.cfg.CtxTimeOut)
	defer cancel()

	var newReport reportmodel.NewReportReq
	err := json.NewDecoder(r.Body).Decode(&newReport)
	if err != nil {
		httpresponse.ResponseError(w, 422, httpresponse.NewError(422, "invalid data"))
		return
	}
	defer r.Body.Close()

	workerID := r.Context().Value(middleware.CtxUserIDKey).(string)

	id, err := h.service.CreateReportOnUser(ctx, workerID, &newReport)
	if err != nil {
		httpresponse.ResponseError(w, 500, err)
		return
	}

	httpresponse.Response(w, 201, reportmodel.NewReportRes{ID: id})
}
