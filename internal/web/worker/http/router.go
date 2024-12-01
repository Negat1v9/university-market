package workerHttp

import (
	"net/http"

	"github.com/Negat1v9/work-marketplace/internal/web/middleware"
)

// RestUserRouter: all routers that return json
func RestWorkerRouter(h *WorkerHandler, mw *middleware.MiddleWareManager) http.Handler {
	workerMux := http.NewServeMux()

	authUserMux := http.NewServeMux()
	authUserMux.HandleFunc("POST /new", h.Create)
	authUserMux.HandleFunc("POST /isworker", h.IsWorker)

	authUserMux.HandleFunc("GET /worker/{id}", h.WorkerPublicInfo)

	priviteMux := http.NewServeMux()
	priviteMux.HandleFunc("GET /profile", h.Worker)
	priviteMux.HandleFunc("PUT /edit/info", h.Update)
	priviteMux.HandleFunc("GET /task/all", h.AvailableTasks)
	priviteMux.HandleFunc("GET /task/info/{id}", h.TaskInfo)
	priviteMux.HandleFunc("POST /task/respond/{id}", h.RespondOnTask)

	priviteMux.HandleFunc("GET /task/responds", h.TasksByResponds)
	priviteMux.HandleFunc("GET /responds", h.Responds)

	workerMux.Handle("/", mw.AuthWorker(priviteMux))
	workerMux.Handle("/user/", http.StripPrefix("/user", mw.AuthUser(authUserMux)))
	return workerMux
}
