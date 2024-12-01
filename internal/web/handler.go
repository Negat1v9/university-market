package web

import (
	"net/http"

	authHttp "github.com/Negat1v9/work-marketplace/internal/web/auth/http"
	commentHttp "github.com/Negat1v9/work-marketplace/internal/web/comment"
	"github.com/Negat1v9/work-marketplace/internal/web/middleware"
	paymentHttp "github.com/Negat1v9/work-marketplace/internal/web/payment/http"
	taskHttp "github.com/Negat1v9/work-marketplace/internal/web/task/http"
	userHttp "github.com/Negat1v9/work-marketplace/internal/web/user/http"
	workerHttp "github.com/Negat1v9/work-marketplace/internal/web/worker/http"
)

func (s *Server) InitRoutes() {
	router := http.NewServeMux()

	mw := middleware.New(s.cfg.WebConfig, s.services.UserService)

	authHandler := authHttp.New(s.cfg.WebConfig, s.services.AuthService)
	taskHandler := taskHttp.New(s.cfg.WebConfig, s.services.TaskService)
	userHandler := userHttp.New(s.cfg.WebConfig, s.services.UserService)
	workerHandler := workerHttp.New(s.cfg.WebConfig, s.services.WorkerService)
	paymentHandler := paymentHttp.New(s.cfg.WebConfig, s.services.PaymentService)
	commentHandler := commentHttp.New(s.cfg.WebConfig, s.services.CommentService)

	authRouter := authHttp.RestAuthRouter(authHandler, mw)
	router.Handle("/auth/", http.StripPrefix("/auth", authRouter))

	taskRouter := taskHttp.RestTaskRouter(taskHandler, mw)
	router.Handle("/task/", http.StripPrefix("/task", taskRouter))

	userRouter := userHttp.RestUserRouter(userHandler, mw)
	router.Handle("/user/", http.StripPrefix("/user", userRouter))

	workerRouter := workerHttp.RestWorkerRouter(workerHandler, mw)
	router.Handle("/worker/", http.StripPrefix("/worker", workerRouter))

	paymentRouter := paymentHttp.RestPaymentRouter(paymentHandler, mw)
	router.Handle("/payment/", http.StripPrefix("/payment", paymentRouter))

	commentRouter := commentHttp.RestPaymentRouter(commentHandler, mw)
	router.Handle("/comment/", http.StripPrefix("/comment", commentRouter))

	apiV1 := http.NewServeMux()

	apiV1.Handle("/api/v1/", http.StripPrefix("/api/v1", router))

	basicMW := mw.BasicMW()

	s.Server.Handler = basicMW(apiV1)
}
