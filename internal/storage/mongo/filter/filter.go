package filters

import (
	paymentmodel "github.com/Negat1v9/work-marketplace/model/payment"
	taskmodel "github.com/Negat1v9/work-marketplace/model/taskModel"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CmplxFilters struct {
	filters bson.D
}

// refactored
func New() *CmplxFilters {
	return &CmplxFilters{
		filters: make(primitive.D, 0),
	}
}

func (f *CmplxFilters) Add(filter bson.E) *CmplxFilters {
	f.filters = append(f.filters, filter)
	return f
}
func (f *CmplxFilters) Filters() bson.D {
	return f.filters
}

// -------------------------------------------------
//	USER FILTERS
// -------------------------------------------------

func UserByID(id string) bson.E {
	userID, _ := primitive.ObjectIDFromHex(id)
	return bson.E{Key: "_id", Value: userID}
}
func UserByTgID(id int64) bson.E {
	return bson.E{Key: "telegram_id", Value: id}
}

// -------------------------------------------------
//	TASK FILTERS
// -------------------------------------------------

func TaskByID(id string) bson.E {
	taskID, _ := primitive.ObjectIDFromHex(id)
	return bson.E{Key: "_id", Value: taskID}
}

func TaskByCreator(creator string) bson.E {
	return bson.E{Key: "created_by", Value: creator}
}
func TaskByAssigned(workerID string) bson.E {
	return bson.E{Key: "assigned_to", Value: workerID}
}
func TaskByIsOpen(open bool) bson.E {
	return bson.E{Key: "is_open", Value: open}
}
func TaskByStatus(status taskmodel.TaskStatus) bson.E {
	return bson.E{Key: "status", Value: status}
}
func TaskByTags(tags []string) bson.E {
	return bson.E{Key: "tags", Value: bson.M{"$in": tags}}
}
func TaskByIsDeleted() bson.E {
	return bson.E{Key: "delete_at", Value: bson.M{"$exists": true}}
}
func TaskByNoDeleted() bson.E {
	return bson.E{Key: "delete_at", Value: bson.M{"$exists": false}}
}

// -------------------------------------------------
//	Payment FILTERS
// -------------------------------------------------

func PaymentByID(id string) bson.E {
	paymentID, _ := primitive.ObjectIDFromHex(id)
	return bson.E{Key: "_id", Value: paymentID}
}

func PaymentByUser(userID string) bson.E {
	return bson.E{Key: "user_id", Value: userID}
}
func PaymentByStatus(status paymentmodel.PaymentStatus) bson.E {
	return bson.E{Key: "status", Value: status}
}

// -------------------------------------------------
//	Respond FILTERS
// -------------------------------------------------

func RespondByWorkerID(workerID string) bson.E {
	return bson.E{Key: "worker_id", Value: workerID}
}

// -------------------------------------------------
//	Comments FILTERS
// -------------------------------------------------

func CommentByCreator(userID string) bson.E {
	return bson.E{Key: "creator_id", Value: userID}
}

func CommentByWorker(workerID string) bson.E {
	return bson.E{Key: "worker_id", Value: workerID}
}
