package reportservice

import (
	reportmodel "github.com/Negat1v9/work-marketplace/model/report"
	httpresponse "github.com/Negat1v9/work-marketplace/pkg/httpResponse"
)

func beforeCreate(newReport *reportmodel.NewReportReq) error {
	switch {
	case newReport.ReportUser == "":
		return httpresponse.NewError(422, "field \"reported_user\" is required")
	case newReport.Reason == "":
		return httpresponse.NewError(422, "field \"reason\" is required")
	case newReport.Description == "":
		return httpresponse.NewError(422, "field \"description\" is required")
	}
	if err := newReport.Validate(); err != nil {
		return httpresponse.NewError(406, err.Error())
	}
	return nil
}
