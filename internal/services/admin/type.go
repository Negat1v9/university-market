package adminservice

import (
	"context"

	eventmodel "github.com/Negat1v9/work-marketplace/model/event"
)

type AdminService interface {
	IsAdmin(userID int64) bool
	CreateEvent(ctx context.Context, event *eventmodel.Event) (string, error)
	StartSendingEvent(ctx context.Context, eventID string) error
}
