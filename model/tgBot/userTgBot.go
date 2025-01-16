package tgbotmodel

type TgExpectedAction int16

var (
	// Waiting for files from the user for a task
	WaitingForFiles       TgExpectedAction = 0
	WorkerShareContact    TgExpectedAction = 1
	WaitingEventCaption   TgExpectedAction = 2
	WaitingStartSendEvent TgExpectedAction = 3
)

type UserCommand struct {
	ID int64 `bson:"_id"`
	// ExpectedaAction what action is expected from the user, for example, the bot is
	// waiting for a file to be sent to it
	ExpectedAction TgExpectedAction `bson:"expected_action"`
	// optional field if need to perform an action on a task,
	// for example, the user must add files
	TaskID string `bson:"task_id,omitempty"`
	// optional field for placing the employee ID, for example, so that the user shares his
	// contact with him
	WorkerID string `bson:"worker_id,omitempty"`
	// optional field for placing the user ID, for example, so that the worker shares his
	// contact with him
	UserID string `bson:"user_id,omitempty"`
	// optional field the field stores the ID of the event to store and add its photo
	EventID string `bson:"event_id,omitempty"`
}
