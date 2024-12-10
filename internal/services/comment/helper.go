package commentservice

import (
	"net/url"

	filters "github.com/Negat1v9/work-marketplace/internal/storage/mongo/filter"
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

func filtersComment(f *filters.CmplxFilters, v url.Values) {
	if likes := v.Get("likes"); likes != "" {
		f.Add(filters.CommentByIsLike(true))
	} else if dislikes := v.Get("dislikes"); dislikes != "" {
		f.Add(filters.CommentByIsLike(false))
	}
}
