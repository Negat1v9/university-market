package taskservice

import (
	"net/url"
	"strings"

	filters "github.com/Negat1v9/work-marketplace/internal/storage/mongo/filter"
	taskmodel "github.com/Negat1v9/work-marketplace/model/taskModel"
	httpresponse "github.com/Negat1v9/work-marketplace/pkg/httpResponse"
)

func FindFilterTasks(v url.Values) *filters.CmplxFilters {
	f := filters.New()
	tags := []string{}
	if filt := v.Get("form_education"); filt != "" {
		tags = append(tags, filt)
	}
	if filt := v.Get("university"); filt != "" {
		tags = append(tags, filt)
	}
	if filt := v.Get("subject"); filt != "" {
		tags = append(tags, filt)
	}
	if filt := v.Get("task_type"); filt != "" {
		tags = append(tags, filt)
	}
	if len(tags) != 0 {
		f.Add(filters.TaskByTags(tags))
	}
	// first selection by tags then status for indexes mongodb
	if filt := v.Get("status"); filt != "" {
		f.Add(filters.TaskByStatus(taskmodel.TaskStatus(filt)))
	}
	return f
}

// check worker repsponded on task
func CheckWorkerRespond(workerID string, responds []string) error {
	for i := 0; i < len(responds); i++ {
		if responds[i] == workerID {
			return nil
		}
	}
	return httpresponse.NewError(404, "no worker in responds")
}
func beforeCreateUpdate(t *taskmodel.TaskMeta) error {
	switch {
	case t.MinPrice == 0:
		return httpresponse.NewError(422, "field \"min_price\" is required")
	case t.MaxPrice == 0:
		return httpresponse.NewError(422, "field \"max_price\" is required")
	case t.FormEducation == "":
		return httpresponse.NewError(422, "field \"form_education\" is required")
	case t.University == "":
		return httpresponse.NewError(422, "field \"university\" is required")
	case t.Subject == "":
		return httpresponse.NewError(422, "field \"subject\" is required")
	case t.TaskType == "":
		return httpresponse.NewError(422, "field \"task_type\" is required")
	case t.Description == "":
		return httpresponse.NewError(422, "field \"description\" is required")
	}

	return nil
}

// consists of an array of tags of the key fields of the task, having previously removed spaces along
// the edges and converted the string to lower case
func createTags(t *taskmodel.TaskMeta) []string {
	tags := []string{}
	if t.FormEducation != "" {
		tags = append(tags,
			strings.ToLower(strings.TrimSpace(t.FormEducation)))
	}
	if t.University != "" {
		tags = append(tags,
			strings.ToLower(strings.TrimSpace(t.University)))
	}
	if t.Subject != "" {
		tags = append(tags,
			strings.ToLower(strings.TrimSpace(t.Subject)))
	}
	if t.TaskType != "" {
		tags = append(tags,
			strings.ToLower(strings.TrimSpace(t.TaskType)))
	}
	return tags
}
