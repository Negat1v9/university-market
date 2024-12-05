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
		"worker_info.full_name":   1,
		"worker_info.education":   1,
		"worker_info.experience":  1,
		"worker_info.description": 1,
		"role":                    1,
	}
	ProjSuccessPayment = bson.M{
		"_id":         1,
		"referral_id": 1,
		"balance":     1,
	}
	ProjOnlyBalance = bson.M{
		"_id":     1,
		"balance": 1,
	}
	ProjCreateWorker = bson.M{
		"_id":                   1,
		"telegram_id":           1,
		"worker_info.full_name": 1,
		"role":                  1,
	}
)
