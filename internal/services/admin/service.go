package adminservice

import (
	"context"
	"log/slog"
	"time"

	"github.com/Negat1v9/work-marketplace/internal/storage"
	mongoStore "github.com/Negat1v9/work-marketplace/internal/storage/mongo"
	filters "github.com/Negat1v9/work-marketplace/internal/storage/mongo/filter"
	tgbot "github.com/Negat1v9/work-marketplace/internal/tgBot"
	eventmodel "github.com/Negat1v9/work-marketplace/model/event"
	usermodel "github.com/Negat1v9/work-marketplace/model/userModel"
	httpresponse "github.com/Negat1v9/work-marketplace/pkg/httpResponse"
)

type AdminServiceImpl struct {
	log         *slog.Logger
	store       storage.Store
	tgClient    tgbot.WebTgClient
	adminsTgIDs []int64
}

func NewServiceAdmin(log *slog.Logger, store storage.Store, c tgbot.WebTgClient, adminsTgIDs []int64) AdminService {
	return &AdminServiceImpl{
		log:         log,
		store:       store,
		tgClient:    c,
		adminsTgIDs: adminsTgIDs,
	}
}

func (s *AdminServiceImpl) IsAdmin(userID int64) bool {
	return checkIsAdmin(userID, s.adminsTgIDs)
}

func (s *AdminServiceImpl) CreateEvent(ctx context.Context, event *eventmodel.Event) (string, error) {
	if err := validateEvent(event); err != nil {
		return "", err
	}

	if isAdmin := checkIsAdmin(event.CreatorID, s.adminsTgIDs); !isAdmin {
		return "", httpresponse.NewError(403, "denied no rights")
	}

	eventID, err := s.store.Event().Create(ctx, event)
	if err != nil {
		s.log.Error("AdminService.CreateEvent.Event.Create", slog.String("err", err.Error()))
		return "", httpresponse.ServerError()
	}

	return eventID, nil
}

func (s *AdminServiceImpl) StartSendingEvent(ctx context.Context, eventID string) error {
	event, err := s.store.Event().FindOne(ctx, filters.New().Add(filters.EventByID(eventID)).Filters())
	if err != nil {
		s.log.Error("AdminService.StartSendingEvent.FindOne", slog.String("err", err.Error()))
		return err
	}

	// Start sending msg in background
	go s.startSendingEvent(event)

	return nil
}

func (s *AdminServiceImpl) startSendingEvent(event *eventmodel.Event) {
	var limit int64 = 1
	var skip int64 = 0
	totalSendMessages := 0
	totalErrorCount := 0
	for {
		users, err := s.getUsersForSednigEvent(limit, skip)
		switch {
		case totalErrorCount >= 5:
			s.log.Error("AdminService.startSendingEvent", slog.Int("total error count", totalErrorCount))
			return
		// exit loop if there are no more users
		case err == mongoStore.ErrNoUser:
			s.log.Debug("exit from sending event messages")
			return
		case err != nil:
			totalErrorCount++
			s.log.Error("AdminService.startSendingEvent.getUsersForSednigEvent", slog.String("err", err.Error()))
			time.Sleep(time.Second * 60)
		default:
			skip += int64(len(users))
			for i := 0; i < len(users); i++ {
				err = s.tgClient.SendEventMsg(users[i].TelegramID, event)
				if err != nil {
					s.log.Error("AdminService.startSendingEvent send msg to "+users[i].ID, slog.String("err", err.Error()))
				} else {
					totalSendMessages += 1
				}
				time.Sleep(time.Millisecond * 100)
			}
		}
	}
}

func (s *AdminServiceImpl) getUsersForSednigEvent(limit, skip int64) ([]usermodel.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	return s.store.User().FindManyProj(ctx, filters.New().Filters(), usermodel.OnlyTgID, limit, skip)

}
