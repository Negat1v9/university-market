package respondmodel

import "time"

type Respond struct {
	ID           string    `bson:"_id,omitempty" json:"id"`
	TaskID       string    `bson:"task_id,omitempty" json:"task_id"`
	WorkerID     string    `bson:"worker_id,omitempty" json:"worker_id"`
	PriceRespond int       `bson:"price_respond,omitempty" json:"price_respond"`   // respond price amount in telegram start "XTR"
	TaskType     string    `bson:"task_type,omitempty" json:"task_type,omitempty"` // task type from task info
	CreatedAt    time.Time `bson:"created_at,omitempty" json:"created_at"`         // Creation date UTC
}

func New(taskID, workerID, taskType string, priceRespond int) *Respond {
	return &Respond{
		TaskID:       taskID,
		WorkerID:     workerID,
		PriceRespond: priceRespond,
		TaskType:     taskType,
		CreatedAt:    time.Now().UTC(),
	}
}
