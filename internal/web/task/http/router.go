package taskHttp

import (
	"net/http"

	"github.com/Negat1v9/work-marketplace/internal/web/middleware"
)

// RestUserRouter: all routers that return json
func RestTaskRouter(h *TaskHandler, mw *middleware.MiddleWareManager) http.Handler {
	taskMux := http.NewServeMux()

	priviteMux := http.NewServeMux()

	priviteMux.HandleFunc("POST /create", h.Create)
	priviteMux.HandleFunc("GET /my", h.FindUserTasks)
	priviteMux.HandleFunc("GET /find/{id}", h.FindOne)
	priviteMux.HandleFunc("PUT /edit/meta/{id}", h.UpdateTaskMeta)
	// priviteMux.HandleFunc("GET /info/{id}/worker/{worker_id}", h.TaskWithWorkerInfo)
	priviteMux.HandleFunc("PUT /raise/{id}", h.RaiseTask)
	priviteMux.HandleFunc("PUT /{id}/select/worker/{worker_id}", h.SelectWorker)

	priviteMux.HandleFunc("PUT /complete/{id}", h.CompleteTask)

	priviteMux.HandleFunc("DELETE /delete/{id}", h.DeleteTask)

	taskMux.Handle("/", mw.AuthUser(priviteMux))
	return taskMux
}
