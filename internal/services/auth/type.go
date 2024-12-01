package authservice

import "context"

type AuthService interface {
	Login(context.Context, string) (string, error)
}
