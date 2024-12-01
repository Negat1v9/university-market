package userservice

import (
	"context"

	usermodel "github.com/Negat1v9/work-marketplace/model/userModel"
)

type UserService interface {
	Create(ctx context.Context, newUser *usermodel.User) (string, error)
	Auth(ctx context.Context, userID string) error
	AuthWorker(ctx context.Context, userID string) error
	User(ctx context.Context, userID string) (*usermodel.User, error)
}
