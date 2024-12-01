package usermodel

import "go.mongodb.org/mongo-driver/bson"

var (
	OnlyID = bson.M{
		"_id": 1,
	}
	AuthWorker = bson.M{
		"_id":             1,
		"role":            1,
		"worker_info.ban": 1,
	}
	OnlyWorkerInfo = bson.M{
		"_id":         1,
		"worker_info": 1,
	}
	OnlyTgID = bson.M{
		"_id":         1,
		"telegram_id": 1,
	}
	WorkerPublic = bson.M{
		"_id":                     1,
		"username":                1,
		"worker_info.karma":       1,
		"worker_info.description": 1,
		"role":                    1,
	}
)
