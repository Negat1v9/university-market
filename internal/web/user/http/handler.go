package userHttp

import (
	"context"
	"net/http"

	"github.com/Negat1v9/work-marketplace/internal/config"
	userservice "github.com/Negat1v9/work-marketplace/internal/services/user"
	"github.com/Negat1v9/work-marketplace/internal/web/middleware"
	httpresponse "github.com/Negat1v9/work-marketplace/pkg/httpResponse"
)

type UserHandler struct {
	cfg     *config.WebCfg
	service userservice.UserService
}

func New(cfg *config.WebCfg, s userservice.UserService) *UserHandler {
	return &UserHandler{
		cfg:     cfg,
		service: s,
	}
}

// func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
// 	ctx, cancel := context.WithTimeout(context.Background(), h.cfg.CtxTimeOut)
// 	defer cancel()

// 	var data usermodel.UserCreate
// 	err := json.NewDecoder(r.Body).Decode(&data)
// 	if err != nil {
// 		httpresponse.ResponseError(w, 422, httpresponse.NewError(422, "invalid data"))
// 		return
// 	}
// 	defer r.Body.Close()

// 	err = tgvalidation.ValidateInitData(data.InitData, h.cfg.TgBotToken)
// 	if err != nil {
// 		httpresponse.ResponseError(w, 401, httpresponse.NewError(401, "forbiden"))
// 		return
// 	}

// 	userID, err := h.service.Create(ctx, &data)
// 	if err != nil {
// 		httpresponse.ResponseError(w, 500, err)
// 		return
// 	}

// 	claims := &utils.Claims{
// 		UserID: userID,
// 	}

// 	token, err := utils.GenerateJwtToken(claims, h.cfg.JwtSecret)
// 	if err != nil {
// 		httpresponse.ResponseError(w, 500, httpresponse.ServerError())
// 		return
// 	}

// 	httpresponse.Response(w, 201, usermodel.LoginRes{TokenType: "Bearer", Token: token})
// }

func (h *UserHandler) User(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), h.cfg.CtxTimeOut)
	defer cancel()

	userID := r.Context().Value(middleware.CtxUserIDKey).(string)

	user, err := h.service.User(ctx, userID)
	if err != nil {
		httpresponse.ResponseError(w, 500, err)
		return
	}

	httpresponse.Response(w, 200, user)
}
