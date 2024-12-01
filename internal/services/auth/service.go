package authservice

import (
	"context"
	"log/slog"

	"github.com/Negat1v9/work-marketplace/internal/storage"
	mongoStore "github.com/Negat1v9/work-marketplace/internal/storage/mongo"
	filters "github.com/Negat1v9/work-marketplace/internal/storage/mongo/filter"
	usermodel "github.com/Negat1v9/work-marketplace/model/userModel"
	httpresponse "github.com/Negat1v9/work-marketplace/pkg/httpResponse"
	tgvalidation "github.com/Negat1v9/work-marketplace/pkg/tgValidation"
)

type AuthServiceImpl struct {
	log   *slog.Logger
	store storage.Store
}

func NewServiceAuth(log *slog.Logger, store storage.Store) AuthService {
	return &AuthServiceImpl{
		log:   log,
		store: store,
	}
}

// Login: checking whether the user exists in the database; if not, it returns an error with code 401
// initData - from telegram.WebApp.InidData
func (s *AuthServiceImpl) Login(ctx context.Context, initData string) (string, error) {

	data, err := tgvalidation.ParseInitData(initData)
	if err != nil {
		return "", httpresponse.NewError(401, "forbiden")
	}

	filter := filters.NewCmplxFilter().Add(filters.UserByTgID(data.User.ID))
	user, err := s.store.User().FindProj(ctx, filter.Filters(), usermodel.OnlyID)
	switch {
	case err == mongoStore.ErrNoUser:
		return "", httpresponse.NewError(401, "forbiden")
	case err != nil:
		s.log.Debug("login user", slog.String("err", err.Error()))
		return "", httpresponse.ServerError()
	}

	return user.ID, nil
}
