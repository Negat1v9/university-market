package userservice

import (
	"context"
	"log/slog"

	"github.com/Negat1v9/work-marketplace/internal/storage"
	mongoStore "github.com/Negat1v9/work-marketplace/internal/storage/mongo"
	filters "github.com/Negat1v9/work-marketplace/internal/storage/mongo/filter"
	usermodel "github.com/Negat1v9/work-marketplace/model/userModel"
	httpresponse "github.com/Negat1v9/work-marketplace/pkg/httpResponse"
)

type UserServiceImpl struct {
	log   *slog.Logger
	store storage.Store
}

func NewServiceUser(log *slog.Logger, store storage.Store) UserService {
	return &UserServiceImpl{
		log:   log,
		store: store,
	}
}

// Create - create new user in database if not exist return userId new user must contain telegram id
func (s *UserServiceImpl) Create(ctx context.Context, newUser *usermodel.User) (string, error) {
	filter := filters.NewCmplxFilter().Add(filters.UserByTgID(newUser.TelegramID))
	user, err := s.store.User().FindProj(ctx, filter.Filters(), usermodel.OnlyID)

	var userID string
	switch {
	case err == mongoStore.ErrNoUser: // no user in database
		userID, err = s.store.User().Create(ctx, newUser)
		if err != nil {
			s.log.Error("create user", slog.String("err", err.Error()))
			return "", err
		}

		return userID, nil

	case err != nil: // mongo error
		s.log.Error("create user", slog.String("err", err.Error()))
		return "", httpresponse.ServerError()
	}

	// return only if user alredy exist
	return user.ID, nil
}

func (s *UserServiceImpl) Auth(ctx context.Context, userID string) error {
	filter := filters.NewCmplxFilter().Add(filters.UserByID(userID))
	_, err := s.store.User().FindProj(ctx, filter.Filters(), usermodel.OnlyID)
	switch {
	case err == mongoStore.ErrNoUser:
		return httpresponse.NewError(401, "not authorized")
	case err != nil:
		return httpresponse.ServerError()
	}
	return nil
}

func (s *UserServiceImpl) AuthWorker(ctx context.Context, userID string) error {
	filter := filters.NewCmplxFilter().Add(filters.UserByID(userID))
	worker, err := s.store.User().FindProj(ctx, filter.Filters(), usermodel.AuthWorker)
	switch {
	case err == mongoStore.ErrNoUser:
		return httpresponse.NewError(401, "not authorized")
	case err != nil:
		return httpresponse.ServerError()
	// the user is not an employee but wants access
	case worker.Role != usermodel.Worker:
		return httpresponse.NewError(403, "forbidden")
	}
	// user is banned from the site
	if worker.WorkerInfo != nil {
		if worker.WorkerInfo.Ban {
			return httpresponse.NewError(403, "forbidden")
		}
	}
	return nil
}

func (s *UserServiceImpl) User(ctx context.Context, userID string) (*usermodel.User, error) {
	user, err := s.store.User().Find(ctx, filters.New().Add(filters.UserByID(userID)).Filters())
	switch {
	case err == mongoStore.ErrNoUser:
		return nil, httpresponse.NewError(404, err.Error())
	case err != nil:
		s.log.Error("user info", slog.String("err", err.Error()))
		return nil, err
	}

	return user, nil
}
