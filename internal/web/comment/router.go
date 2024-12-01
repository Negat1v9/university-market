package commentHttp

import (
	"net/http"

	"github.com/Negat1v9/work-marketplace/internal/web/middleware"
)

func RestPaymentRouter(h *CommentHandler, mw *middleware.MiddleWareManager) http.Handler {
	commentMux := http.NewServeMux()

	userAuth := http.NewServeMux()

	userAuth.HandleFunc("POST /create", h.Create)
	userAuth.HandleFunc("GET /my", h.UserComments)
	userAuth.HandleFunc("GET /worker/{id}", h.UserWorkerComments)

	workerAuth := http.NewServeMux()
	workerAuth.HandleFunc("GET /my", h.WorkerComments)

	commentMux.Handle("/user/", http.StripPrefix("/user", mw.AuthUser(userAuth)))
	commentMux.Handle("/worker/", http.StripPrefix("/worker", mw.AuthWorker(workerAuth)))

	return commentMux
}
