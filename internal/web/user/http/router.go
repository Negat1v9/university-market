package userHttp

import (
	"net/http"

	"github.com/Negat1v9/work-marketplace/internal/web/middleware"
)

// RestUserRouter: all routers that return json
func RestUserRouter(h *UserHandler, mw *middleware.MiddleWareManager) http.Handler {
	userMux := http.NewServeMux()

	// userMux.HandleFunc("POST /create", h.CreateUser)

	priviteMux := http.NewServeMux()
	priviteMux.HandleFunc("GET /info/{id}", h.User)
	// userMux.HandleFunc("PUT /update/{id}", h.Update)

	userMux.Handle("/", http.StripPrefix("/", mw.AuthUser(priviteMux)))

	return userMux
}
