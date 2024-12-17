package reportservice

import (
	"context"
	"log/slog"

	"github.com/Negat1v9/work-marketplace/internal/storage"
	mongoStore "github.com/Negat1v9/work-marketplace/internal/storage/mongo"
	filters "github.com/Negat1v9/work-marketplace/internal/storage/mongo/filter"
	reportmodel "github.com/Negat1v9/work-marketplace/model/report"
	usermodel "github.com/Negat1v9/work-marketplace/model/userModel"
	httpresponse "github.com/Negat1v9/work-marketplace/pkg/httpResponse"
)

type ReportServiceImpl struct {
	log   *slog.Logger
	store storage.Store
}

func NewServiceReport(log *slog.Logger, store storage.Store) ReportService {
	return &ReportServiceImpl{
		log:   log,
		store: store,
	}
}

func (s *ReportServiceImpl) CreateReportOnWorker(ctx context.Context, userID string, newReport *reportmodel.NewReportReq) (string, error) {
	if err := beforeCreate(newReport); err != nil {
		return "", err
	}
	if userID == newReport.ReportUser {
		return "", httpresponse.NewError(409, "report self")
	}
	worker, err := s.store.User().FindProj(
		ctx,
		filters.New().Add(filters.UserByID(newReport.ReportUser)).Filters(),
		usermodel.AuthWorker,
	)
	switch {
	case err == mongoStore.ErrNoUser:
		return "", httpresponse.NewError(404, "no worker")
	case err != nil:
		s.log.Error("create report on worker", slog.String("err", err.Error()))
		return "", httpresponse.ServerError()
	case worker.Role != usermodel.Worker:
		return "", httpresponse.NewError(404, "no worker")
	}

	report, err := s.store.Report().FindOne(
		ctx,
		filters.New().
			Add(filters.ReportByReportedBy(userID)).
			Add(filters.ReportByReportedUser(newReport.ReportUser)).
			Add(filters.ReportByReporterByRole(usermodel.RegularUser)).Filters(),
	)
	switch {
	// there has already been a report from this user
	case report != nil:
		return report.ID, nil
	case err != nil && err != mongoStore.ErrNoReport:
		s.log.Error("create report on worker", slog.String("err", err.Error()))
		return "", httpresponse.ServerError()
	}

	report = reportmodel.NewReport(userID, usermodel.RegularUser, newReport)
	newReportID, err := s.store.Report().Create(ctx, report)
	if err != nil {
		s.log.Error("create report on worker", slog.String("err", err.Error()))
		return "", httpresponse.ServerError()
	}

	return newReportID, nil
}
func (s *ReportServiceImpl) CreateReportOnUser(ctx context.Context, workerID string, newReport *reportmodel.NewReportReq) (string, error) {
	if err := beforeCreate(newReport); err != nil {
		return "", err
	}
	if workerID == newReport.ReportUser {
		return "", httpresponse.NewError(409, "report self")
	}
	_, err := s.store.User().FindProj(
		ctx,
		filters.New().Add(filters.UserByID(newReport.ReportUser)).Filters(),
		usermodel.OnlyID,
	)
	switch {
	case err == mongoStore.ErrNoUser:
		return "", httpresponse.NewError(404, "no user")
	case err != nil:
		s.log.Error("create report on user", slog.String("err", err.Error()))
		return "", httpresponse.ServerError()
	}

	report, err := s.store.Report().FindOne(
		ctx,
		filters.New().
			Add(filters.ReportByReportedBy(workerID)).
			Add(filters.ReportByReportedUser(newReport.ReportUser)).
			Add(filters.ReportByReporterByRole(usermodel.Worker)).Filters(),
	)
	switch {
	// there has already been a report from this worker
	case report != nil:
		return report.ID, nil
	case err != nil && err != mongoStore.ErrNoReport:
		s.log.Error("create report on worker", slog.String("err", err.Error()))
		return "", httpresponse.ServerError()
	}

	report = reportmodel.NewReport(workerID, usermodel.Worker, newReport)
	newReportID, err := s.store.Report().Create(ctx, report)
	if err != nil {
		s.log.Error("create report on worker", slog.String("err", err.Error()))
		return "", httpresponse.ServerError()
	}

	return newReportID, nil
}
