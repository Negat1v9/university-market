package sercicesMock

import (
	"log/slog"

	"github.com/Negat1v9/work-marketplace/internal/config"
	authservice "github.com/Negat1v9/work-marketplace/internal/services/auth"
	auth_service_mock "github.com/Negat1v9/work-marketplace/internal/services/auth/mock"
	commentservice "github.com/Negat1v9/work-marketplace/internal/services/comment"
	paymentservice "github.com/Negat1v9/work-marketplace/internal/services/payment"
	taskservice "github.com/Negat1v9/work-marketplace/internal/services/task"
	task_service_mock "github.com/Negat1v9/work-marketplace/internal/services/task/mock"
	userservice "github.com/Negat1v9/work-marketplace/internal/services/user"
	user_service_mock "github.com/Negat1v9/work-marketplace/internal/services/user/mock"
	workerservice "github.com/Negat1v9/work-marketplace/internal/services/worker"
	worker_service_mock "github.com/Negat1v9/work-marketplace/internal/services/worker/mock"
	"github.com/Negat1v9/work-marketplace/internal/storage"
	tgbot "github.com/Negat1v9/work-marketplace/internal/tgBot"
	"github.com/golang/mock/gomock"
)

type ServicesMock struct {
	AuthService    authservice.AuthService
	PaymentService paymentservice.PaymentService
	TaskService    taskservice.TaskService
	UserService    userservice.UserService
	WorkerService  workerservice.WorkerService
	CommentService commentservice.CommentService
}

func NewServiceMockBuilder(cfg *config.Config, log *slog.Logger, tgClient tgbot.WebTgClient, store storage.Store, cntl *gomock.Controller) *ServicesMock {
	return &ServicesMock{
		AuthService:    auth_service_mock.NewMockAuthService(cntl),
		PaymentService: paymentservice.NewServicePayment(log, cfg.BotConfig.BotToken, store),
		TaskService:    task_service_mock.NewMockTaskService(cntl),
		UserService:    user_service_mock.NewMockUserService(cntl),
		WorkerService:  worker_service_mock.NewMockWorkerService(cntl),
		CommentService: commentservice.NewServiceComment(log, store),
	}
}
