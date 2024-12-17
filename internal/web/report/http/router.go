package reportHttp

import (
	"net/http"

	"github.com/Negat1v9/work-marketplace/internal/web/middleware"
)

func RestReportRouter(h *ReportHandler, mw *middleware.MiddleWareManager) http.Handler {
	reportMux := http.NewServeMux()

	userAuthMux := http.NewServeMux()
	userAuthMux.HandleFunc("POST /create", h.CreateReportOnWorker)
	workerAuthMux := http.NewServeMux()
	workerAuthMux.HandleFunc("POST /create", h.CreateReportOnUser)

	reportMux.Handle("/user/", http.StripPrefix("/user", mw.AuthUser(userAuthMux)))
	reportMux.Handle("/worker/", http.StripPrefix("/worker", mw.AuthWorker(workerAuthMux)))

	return reportMux
}
