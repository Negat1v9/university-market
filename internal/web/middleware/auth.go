package middleware

import (
	"context"
	"net/http"

	httpresponse "github.com/Negat1v9/work-marketplace/pkg/httpResponse"
	"github.com/Negat1v9/work-marketplace/pkg/utils"
)

// Info: user authorization via an authentication token. Transmits userID in request context (type *string)
func (m *MiddleWareManager) AuthUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenStr := r.Header.Get("Authorization")

		tokenClaims, err := utils.JwtClaimsFromToken(tokenStr, m.cfg.JwtSecret)
		if err != nil {
			httpresponse.ResponseError(w, 401, httpresponse.NewError(401, "not authorized"))
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), m.cfg.CtxTimeOut)
		defer cancel()

		err = m.userserice.Auth(ctx, tokenClaims.UserID)
		if err != nil {
			httpresponse.ResponseError(w, 500, err)
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), CtxUserIDKey, tokenClaims.UserID)))
	})
}

func (m *MiddleWareManager) AuthWorker(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenStr := r.Header.Get("Authorization")

		tokenClaims, err := utils.JwtClaimsFromToken(tokenStr, m.cfg.JwtSecret)
		if err != nil {
			httpresponse.ResponseError(w, 403, httpresponse.NewError(403, "not authorized"))
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), m.cfg.CtxTimeOut)
		defer cancel()

		err = m.userserice.AuthWorker(ctx, tokenClaims.UserID)
		if err != nil {
			httpresponse.ResponseError(w, 500, err)
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), CtxUserIDKey, tokenClaims.UserID)))
	})
}
