package paymentHttp

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Negat1v9/work-marketplace/internal/config"
	paymentservice "github.com/Negat1v9/work-marketplace/internal/services/payment"
	"github.com/Negat1v9/work-marketplace/internal/web/middleware"
	paymentmodel "github.com/Negat1v9/work-marketplace/model/payment"
	httpresponse "github.com/Negat1v9/work-marketplace/pkg/httpResponse"
)

type PaymenetHandler struct {
	cfg     *config.WebCfg
	service paymentservice.PaymentService
}

func New(cfg *config.WebCfg, s paymentservice.PaymentService) *PaymenetHandler {
	return &PaymenetHandler{
		cfg:     cfg,
		service: s,
	}
}

func (h *PaymenetHandler) CreateInvoiceLink(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), h.cfg.CtxTimeOut)
	defer cancel()

	var data paymentmodel.PaymentLinkReq
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		httpresponse.ResponseError(w, 422, httpresponse.NewError(422, "invalid data"))
		return
	}

	defer r.Body.Close()

	userID := r.Context().Value(middleware.CtxUserIDKey).(string)

	link, err := h.service.CreateInvoiceLink(ctx, userID, &data)
	if err != nil {
		httpresponse.ResponseError(w, 500, err)
		return
	}

	httpresponse.Response(w, 200, paymentmodel.PaymentLinkRes{Link: link})
}
