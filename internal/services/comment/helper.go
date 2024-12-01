package commentservice

import (
	commentmodel "github.com/Negat1v9/work-marketplace/model/comment"
	httpresponse "github.com/Negat1v9/work-marketplace/pkg/httpResponse"
)

func beforeCreate(creatorID string, c *commentmodel.Comment) error {
	switch {
	case c.TaskID == "":
		return httpresponse.NewError(422, "field \"task_id\" is required")
	case c.TaskType == "":
		return httpresponse.NewError(422, "field \"task_type\" is required")
	case c.CreatorID == "":
		return httpresponse.NewError(422, "field \"creator_id\" is required")
	case c.CreatorID != creatorID:
		return httpresponse.NewError(422, "\"creator_id\" doesn't match")
	case c.WorkerID == "":
		return httpresponse.NewError(422, "field \"worker_id\" is required")
	case c.Description == "":
		return httpresponse.NewError(422, "field \"description\" is required")
	}

	return nil
}
