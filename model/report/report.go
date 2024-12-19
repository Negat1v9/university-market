package reportmodel

import (
	"time"

	usermodel "github.com/Negat1v9/work-marketplace/model/userModel"
)

type ReportUser struct {
	ID   string             `bson:"id,omitempty" json:"id"`  // unique identifier generated by mongo
	Role usermodel.UserType `bson:"role,omitempty" json:"-"` // user role not public
}

type Report struct {
	ID          string     `bson:"_id,omitempty" json:"id"`                          // unique identifier generated by mongo
	ReportBy    ReportUser `bson:"reported_by,omitempty" json:"reported_by"`         // user who create report
	ReportUser  ReportUser `bson:"reported_user,omitempty" json:"reported_user"`     // reported user
	TaskID      string     `bson:"task_id,omitempty" json:"task_id,omitempty"`       // (not required) // task id if report is by task
	Reason      string     `bson:"reason,omitempty" json:"reason"`                   // reason of report
	Description string     `bson:"description,omitempty" json:"description"`         // description report
	CreatedAt   time.Time  `bson:"created_at,omitempty" json:"created_at"`           // time then report was created
	UpdatedAt   time.Time  `bson:"updated_at,omitempty" json:"updated_at,omitempty"` // (not required) time then report was updated
}

// creating a new report, all required fields must not be empty also check that the user against whom the report exists
func NewReport(reportBy string, reportRole usermodel.UserType, newReport *NewReportReq) *Report {
	report := &Report{
		ReportBy: ReportUser{
			ID: reportBy,
			// Role: usermodel.RegularUser,
		},
		ReportUser: ReportUser{
			ID: newReport.ReportUser,
			// Role: usermodel.Worker,
		},
		TaskID:      newReport.TaskID,
		Reason:      newReport.Reason,
		Description: newReport.Description,
		CreatedAt:   time.Now().UTC(),
	}
	if reportRole == usermodel.RegularUser {
		report.ReportBy.Role = usermodel.RegularUser
		report.ReportUser.Role = usermodel.Worker
	} else {
		report.ReportBy.Role = usermodel.Worker
		report.ReportUser.Role = usermodel.RegularUser
	}
	return report
}

type NewReportReq struct {
	ReportUser  string `json:"reported_user"`
	TaskID      string `json:"task_id,omitempty"`
	Reason      string `json:"reason"`
	Description string `json:"description"`
}
type NewReportRes struct {
	ID string `json:"id"`
}
