package authHttp

import (
	"net/http"

	"github.com/Negat1v9/work-marketplace/internal/web/middleware"
)

// RestAuthRouter: all routers that return json
func RestAuthRouter(h *AuthHandler, mw *middleware.MiddleWareManager) http.Handler {
	authMux := http.NewServeMux()
	authMux.HandleFunc("POST /login", h.Login)
	// authMux.HandleFunc("POST /refresh", )

	return authMux
}
