package authHttp

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Negat1v9/work-marketplace/internal/config"
	authservice "github.com/Negat1v9/work-marketplace/internal/services/auth"
	usermodel "github.com/Negat1v9/work-marketplace/model/userModel"
	httpresponse "github.com/Negat1v9/work-marketplace/pkg/httpResponse"
	tgvalidation "github.com/Negat1v9/work-marketplace/pkg/tgValidation"
	"github.com/Negat1v9/work-marketplace/pkg/utils"
)

type AuthHandler struct {
	cfg     *config.WebCfg
	service authservice.AuthService
}

func New(cfg *config.WebCfg, s authservice.AuthService) *AuthHandler {
	return &AuthHandler{
		cfg:     cfg,
		service: s,
	}
}

// Info: user login to the application, checking whether it exists in the database and
// checking authentication information
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), h.cfg.CtxTimeOut)
	defer cancel()

	var data usermodel.UserLoginReq

	err := json.NewDecoder(r.Body).Decode(&data)
	defer r.Body.Close()
	if err != nil {
		httpresponse.ResponseError(w, 422, httpresponse.NewError(422, "not valid"))
		return
	}

	err = tgvalidation.ValidateInitData(data.InitData, h.cfg.TgBotToken)
	if err != nil {
		httpresponse.ResponseError(w, 401, httpresponse.NewError(401, "forbiden"))
		return
	}

	userID, err := h.service.Login(ctx, data.InitData)
	if err != nil {
		httpresponse.ResponseError(w, 500, err)
		return
	}

	claims := &utils.Claims{
		UserID: userID,
	}

	token, err := utils.GenerateJwtToken(claims, h.cfg.JwtSecret)
	if err != nil {
		httpresponse.ResponseError(w, 500, httpresponse.ServerError())
		return
	}

	res := usermodel.LoginRes{
		TokenType: "Bearer",
		Token:     token,
	}

	httpresponse.Response(w, 200, res)
}
