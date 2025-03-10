package services

import (
	"log/slog"

	"github.com/Negat1v9/work-marketplace/internal/config"
	adminservice "github.com/Negat1v9/work-marketplace/internal/services/admin"
	authservice "github.com/Negat1v9/work-marketplace/internal/services/auth"
	commentservice "github.com/Negat1v9/work-marketplace/internal/services/comment"
	paymentservice "github.com/Negat1v9/work-marketplace/internal/services/payment"
	reportservice "github.com/Negat1v9/work-marketplace/internal/services/report"
	taskservice "github.com/Negat1v9/work-marketplace/internal/services/task"
	userservice "github.com/Negat1v9/work-marketplace/internal/services/user"
	workerservice "github.com/Negat1v9/work-marketplace/internal/services/worker"
	"github.com/Negat1v9/work-marketplace/internal/storage"
	tgbot "github.com/Negat1v9/work-marketplace/internal/tgBot"
)

type Services struct {
	AuthService    authservice.AuthService
	PaymentService paymentservice.PaymentService
	TaskService    taskservice.TaskService
	UserService    userservice.UserService
	WorkerService  workerservice.WorkerService
	CommentService commentservice.CommentService
	ReportService  reportservice.ReportService
	AdmindService  adminservice.AdminService
}

func NewServiceBuilder(cfg *config.Config, log *slog.Logger, tgClient tgbot.WebTgClient, store storage.Store) *Services {
	return &Services{
		AuthService:    authservice.NewServiceAuth(log, store),
		PaymentService: paymentservice.NewServicePayment(log, cfg.WebConfig.TgBotToken, store),
		TaskService:    taskservice.NewServiceTask(log, tgClient, store),
		UserService:    userservice.NewServiceUser(log, store),
		WorkerService:  workerservice.NewServiceWorker(log, tgClient, store, cfg.WebConfig.TgBotToken),
		CommentService: commentservice.NewServiceComment(log, store),
		ReportService:  reportservice.NewServiceReport(log, store),
		AdmindService:  adminservice.NewServiceAdmin(log, store, tgClient, cfg.WebConfig.AdminsIDs),
	}
}
