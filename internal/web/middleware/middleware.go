package middleware

import (
	"net/http"

	"github.com/Negat1v9/work-marketplace/internal/config"
	userservice "github.com/Negat1v9/work-marketplace/internal/services/user"
)

const (
	CtxUserIDKey int = 0
)

type Middleware func(http.Handler) http.Handler

// Info: used to build all the necessary middleware as a constructor
type MiddleWareManager struct {
	cfg        *config.WebCfg
	userserice userservice.UserService
}

func New(cfg *config.WebCfg, userserice userservice.UserService) *MiddleWareManager {
	return &MiddleWareManager{
		cfg:        cfg,
		userserice: userserice,
	}
}

// Info: create middleware for all requests for CORS, logging and same
func (mw *MiddleWareManager) BasicMW() Middleware {
	return createStack(cors)
}

func createStack(xs ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := 0; i < len(xs); i++ {
			x := xs[i]
			next = x(next)
		}
		return next
	}
}
