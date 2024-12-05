package sercicesMock

import (
	authservice "github.com/Negat1v9/work-marketplace/internal/services/auth"
	auth_service_mock "github.com/Negat1v9/work-marketplace/internal/services/auth/mock"
	commentservice "github.com/Negat1v9/work-marketplace/internal/services/comment"
	comment_service_mock "github.com/Negat1v9/work-marketplace/internal/services/comment/mock"
	paymentservice "github.com/Negat1v9/work-marketplace/internal/services/payment"
	payment_service_mock "github.com/Negat1v9/work-marketplace/internal/services/payment/mock"
	taskservice "github.com/Negat1v9/work-marketplace/internal/services/task"
	task_service_mock "github.com/Negat1v9/work-marketplace/internal/services/task/mock"
	userservice "github.com/Negat1v9/work-marketplace/internal/services/user"
	user_service_mock "github.com/Negat1v9/work-marketplace/internal/services/user/mock"
	workerservice "github.com/Negat1v9/work-marketplace/internal/services/worker"
	worker_service_mock "github.com/Negat1v9/work-marketplace/internal/services/worker/mock"
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

func NewServiceMockBuilder(cntl *gomock.Controller) *ServicesMock {
	return &ServicesMock{
		AuthService:    auth_service_mock.NewMockAuthService(cntl),
		PaymentService: payment_service_mock.NewMockPaymentService(cntl),
		TaskService:    task_service_mock.NewMockTaskService(cntl),
		UserService:    user_service_mock.NewMockUserService(cntl),
		WorkerService:  worker_service_mock.NewMockWorkerService(cntl),
		CommentService: comment_service_mock.NewMockCommentService(cntl),
	}
}
