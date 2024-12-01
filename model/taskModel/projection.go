package taskmodel

import "go.mongodb.org/mongo-driver/bson"

var (
	OnlyStatus = bson.M{
		"_id":    1,
		"status": 1,
	}
	OnlyMeta = bson.M{
		"_id":  1,
		"meta": 1,
	}
	ManyTasks = bson.M{
		"_id":                 1,
		"status":              1,
		"meta.max_days_work":  1,
		"meta.min_price":      1,
		"meta.max_price":      1,
		"meta.form_education": 1,
		"meta.university":     1,
		"meta.subject":        1,
		"meta.task_type":      1,
		"created_at":          1,
	}
	ProjOnRespond = bson.M{
		"_id":            1,
		"created_by":     1,
		"status":         1,
		"meta.min_price": 1,
		"meta.task_type": 1,
		"responds":       1,
	}
	ProjOnAttachFiles = bson.M{
		"_id":      1,
		"status":   1,
		"files_id": 1,
	}
	ProjOnCreateComment = bson.M{
		"_id":         1,
		"is_comment":  1,
		"assigned_to": 1,
		"status":      1,
	}
)
