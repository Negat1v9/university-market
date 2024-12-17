package reportservice

import (
	"context"

	reportmodel "github.com/Negat1v9/work-marketplace/model/report"
)

type ReportService interface {
	CreateReportOnWorker(ctx context.Context, userID string, report *reportmodel.NewReportReq) (string, error)
	CreateReportOnUser(ctx context.Context, workerID string, report *reportmodel.NewReportReq) (string, error)
}
