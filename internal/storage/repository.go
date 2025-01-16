package storage

import (
	"context"

	commentmodel "github.com/Negat1v9/work-marketplace/model/comment"
	eventmodel "github.com/Negat1v9/work-marketplace/model/event"
	paymentmodel "github.com/Negat1v9/work-marketplace/model/payment"
	reportmodel "github.com/Negat1v9/work-marketplace/model/report"
	respondmodel "github.com/Negat1v9/work-marketplace/model/respond"
	taskmodel "github.com/Negat1v9/work-marketplace/model/taskModel"
	tgbotmodel "github.com/Negat1v9/work-marketplace/model/tgBot"
	usermodel "github.com/Negat1v9/work-marketplace/model/userModel"
	"go.mongodb.org/mongo-driver/bson"
)

type UserRepository interface {
	Create(context.Context, *usermodel.User) (string, error)
	Find(context.Context, bson.D) (*usermodel.User, error)
	FindProj(ctx context.Context, filter bson.D, proj bson.M) (*usermodel.User, error)
	FindManyProj(ctx context.Context, filter bson.D, proj bson.M, limit, skip int64) ([]usermodel.User, error)
	Edit(ctx context.Context, filter bson.D, upd *usermodel.User) (*usermodel.User, error)
}

type PaymentRepository interface {
	Create(context.Context, *paymentmodel.Payment) (string, error)
	Find(context.Context, bson.D) (*paymentmodel.Payment, error)
	Edit(context.Context, bson.D, *paymentmodel.Payment) error
	Delete(context.Context, bson.D) error
}

type TaskRepository interface {
	Create(ctx context.Context, task *taskmodel.Task) (string, error)
	Find(ctx context.Context, filter bson.D) (*taskmodel.Task, error)
	FindProj(ctx context.Context, filter bson.D, proj bson.M) (*taskmodel.Task, error)
	FindMany(ctx context.Context, filter bson.D, proj bson.M, limit, skip int64) ([]taskmodel.Task, error)
	Count(ctx context.Context, filter bson.D) (int64, error)
	Update(ctx context.Context, filter bson.D, task *taskmodel.Task) (*taskmodel.Task, error)
	Delete(ctx context.Context, filter bson.D) error
}

type TgCommandRepository interface {
	Create(context.Context, *tgbotmodel.UserCommand) error
	// CreateDelete - create new command if command alredy exist delete it and try create again
	// CreateDelete(context.Context, *tgbotmodel.UserCommand) error
	Find(context.Context, int64) (*tgbotmodel.UserCommand, error)
	Delete(context.Context, int64) error
}

type RespondRepository interface {
	Create(context.Context, *respondmodel.Respond) (string, error)
	Find(context.Context, bson.D) (*respondmodel.Respond, error)
	FindMany(ctx context.Context, filter bson.D, limit, skip int64) ([]respondmodel.Respond, error)
	Delete(context.Context, bson.D) error
}

type CommentRepository interface {
	Create(context.Context, *commentmodel.Comment) (string, error)
	FindMany(ctx context.Context, filter bson.D, limit, skip int64) ([]commentmodel.Comment, error)
	CountWorkerLikesDislikes(ctx context.Context, workerID string) (*commentmodel.CountLikeDislikeWorker, error)
	Update(context.Context, bson.D, *commentmodel.Comment) error
	Delete(context.Context, bson.D) error
}

type ReportRepository interface {
	Create(context.Context, *reportmodel.Report) (string, error)
	FindOne(context.Context, bson.D) (*reportmodel.Report, error)
}

type EventRepository interface {
	Create(ctx context.Context, event *eventmodel.Event) (string, error)
	Update(ctx context.Context, filter bson.D, event *eventmodel.Event) (*eventmodel.Event, error)
	FindOne(ctx context.Context, filter bson.D) (*eventmodel.Event, error)
}
