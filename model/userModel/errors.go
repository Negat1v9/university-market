package usermodel

import (
	"fmt"
)

func ErrWorkerInfoIsNil(userID string) error {
	return fmt.Errorf("user %s has no field worker_info", userID)
}
