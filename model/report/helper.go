package reportmodel

import "fmt"

func (r *NewReportReq) Validate() error {
	switch {
	case len(r.ReportUser) > 30:
		return fmt.Errorf("max field \"reported_user\" length is %d", 30)
	case len(r.TaskID) > 30:
		return fmt.Errorf("max field \"task_id\" length is %d", 30)
	case len(r.Reason) > 50:
		return fmt.Errorf("max field \"reason\" length is %d", 50)
	case len(r.Description) > 400:
		return fmt.Errorf("max field \"description\" length is %d", 400)
	}
	return nil
}
