package paymentHttp

import (
	"net/http"

	"github.com/Negat1v9/work-marketplace/internal/web/middleware"
)

func RestPaymentRouter(h *PaymenetHandler, mw *middleware.MiddleWareManager) http.Handler {
	payMux := http.NewServeMux()

	priviteMux := http.NewServeMux()

	priviteMux.HandleFunc("POST /create", h.CreateInvoiceLink)

	payMux.Handle("/", mw.AuthWorker(priviteMux))

	return payMux
}
